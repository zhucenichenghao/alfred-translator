package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const youdaoEndpoint = "https://openapi.youdao.com/api"

type Youdao struct {
	key        string
	secret     string
	saltSource SaltSource
}

type youdaoResponse struct {
	ErrorCode   *string         `json:"errorCode"`
	Translation []string        `json:"translation"`
	Basic       *youdaoBasic    `json:"basic"`
	Web         []youdaoWebItem `json:"web"`
}

type youdaoBasic struct {
	Explains   []string `json:"explains"`
	Phonetic   string   `json:"phonetic"`
	USPhonetic string   `json:"us-phonetic"`
	UKPhonetic string   `json:"uk-phonetic"`
}

type youdaoWebItem struct {
	Key   string   `json:"key"`
	Value []string `json:"value"`
}

func NewYoudao(key, secret string, saltSource SaltSource) *Youdao {
	return &Youdao{key: key, secret: secret, saltSource: saltSource}
}

func (p *Youdao) BuildRequest(ctx context.Context, query Query) (*http.Request, error) {
	salt, err := p.saltSource()
	if err != nil {
		return nil, err
	}
	from, to := "auto", "zh-CHS"
	if query.HasHan {
		from, to = "zh-CHS", "en"
	}
	values := url.Values{
		"q":      {query.Text},
		"from":   {from},
		"to":     {to},
		"appKey": {p.key},
		"salt":   {salt},
		"sign":   {signature(p.key, query.Text, salt, p.secret)},
	}
	return http.NewRequestWithContext(ctx, http.MethodGet, youdaoEndpoint+"?"+values.Encode(), nil)
}

func (p *Youdao) ParseResponse(query Query, body []byte) ([]Result, error) {
	var response youdaoResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("Youdao 返回无效 JSON：%w", err)
	}
	if response.ErrorCode == nil {
		return nil, fmt.Errorf("Youdao：%w", ErrInvalidResponse)
	}
	if *response.ErrorCode != "0" {
		return nil, providerError("Youdao", *response.ErrorCode, youdaoErrorMessage(*response.ErrorCode))
	}

	quickLookURL := "https://www.youdao.com/w/" + encodedPath(query.Text)
	results := make([]Result, 0)
	if len(response.Translation) > 0 && nonEmpty(response.Translation[0]) {
		translation := response.Translation[0]
		pronounce := query.Text
		if query.HasHan {
			pronounce = translation
		}
		addResult(&results, Result{Title: translation, Subtitle: query.Text, Arg: translation, Pronounce: pronounce, QuickLookURL: quickLookURL})
	}

	if response.Basic != nil {
		pronounce := query.Text
		for _, explain := range response.Basic.Explains {
			if !nonEmpty(explain) {
				continue
			}
			if query.HasHan {
				pronounce = explain
			}
			addResult(&results, Result{Title: explain, Subtitle: query.Text, Arg: explain, Pronounce: pronounce, QuickLookURL: quickLookURL})
		}
		phonetic := formatPhonetic(query.HasHan, response.Basic)
		if phonetic != "" && pronounce != "" {
			addResult(&results, Result{Title: phonetic, Subtitle: "回车可听发音", Arg: "~" + pronounce, Pronounce: pronounce, QuickLookURL: quickLookURL})
		}
	}

	for _, item := range response.Web {
		values := nonEmptyStrings(item.Value)
		if !nonEmpty(item.Key) || len(values) == 0 {
			continue
		}
		pronounce := item.Key
		if query.HasHan {
			pronounce = values[0]
		}
		addResult(&results, Result{Title: strings.Join(values, ", "), Subtitle: item.Key, Arg: values[0], Pronounce: pronounce, QuickLookURL: quickLookURL})
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("Youdao：%w", ErrInvalidResponse)
	}
	return results, nil
}

func formatPhonetic(hasHan bool, basic *youdaoBasic) string {
	var builder strings.Builder
	if hasHan && basic.Phonetic != "" {
		builder.WriteString("[")
		builder.WriteString(basic.Phonetic)
		builder.WriteString("] ")
	}
	if basic.USPhonetic != "" {
		builder.WriteString(" [美: ")
		builder.WriteString(basic.USPhonetic)
		builder.WriteString("] ")
	}
	if basic.UKPhonetic != "" {
		builder.WriteString(" [英: ")
		builder.WriteString(basic.UKPhonetic)
		builder.WriteString("]")
	}
	return builder.String()
}

func nonEmptyStrings(values []string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		if nonEmpty(value) {
			result = append(result, value)
		}
	}
	return result
}

func youdaoErrorMessage(code string) string {
	messages := map[string]string{
		"101": "缺少必填的参数",
		"102": "不支持的语言类型",
		"103": "翻译文本过长",
		"108": "应用ID无效",
		"110": "无相关服务的有效实例",
		"111": "开发者账号无效",
		"112": "请求服务无效",
		"113": "查询为空",
		"202": "签名检验失败,检查 KEY 和 SECRET",
		"401": "账户已经欠费",
		"411": "访问频率受限",
	}
	if message, ok := messages[code]; ok {
		return message
	}
	return "请参考错误码：" + code
}
