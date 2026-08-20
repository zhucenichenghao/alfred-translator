package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const baiduEndpoint = "https://fanyi-api.baidu.com/api/trans/vip/translate"

type Baidu struct {
	key        string
	secret     string
	saltSource SaltSource
}

type baiduResponse struct {
	ErrorCode   flexibleString     `json:"error_code"`
	TransResult []baiduTranslation `json:"trans_result"`
}

type baiduTranslation struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

type flexibleString string

func (s *flexibleString) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		*s = ""
		return nil
	}
	var text string
	if err := json.Unmarshal(data, &text); err == nil {
		*s = flexibleString(text)
		return nil
	}
	var number json.Number
	if err := json.Unmarshal(data, &number); err != nil {
		return err
	}
	*s = flexibleString(number.String())
	return nil
}

func NewBaidu(key, secret string, saltSource SaltSource) *Baidu {
	return &Baidu{key: key, secret: secret, saltSource: saltSource}
}

func (p *Baidu) BuildRequest(ctx context.Context, query Query) (*http.Request, error) {
	salt, err := p.saltSource()
	if err != nil {
		return nil, err
	}
	from, to := "auto", "zh"
	if query.HasHan {
		from, to = "zh", "en"
	}
	values := url.Values{
		"q":      {query.Text},
		"from":   {from},
		"to":     {to},
		"appid":  {p.key},
		"salt":   {salt},
		"sign":   {signature(p.key, query.Text, salt, p.secret)},
		"dict":   {"1"},
		"action": {"1"},
	}
	return http.NewRequestWithContext(ctx, http.MethodGet, baiduEndpoint+"?"+values.Encode(), nil)
}

func (p *Baidu) ParseResponse(query Query, body []byte) ([]Result, error) {
	var response baiduResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("Baidu 返回无效 JSON：%w", err)
	}
	if response.ErrorCode != "" {
		code := string(response.ErrorCode)
		return nil, providerError("Baidu", code, baiduErrorMessage(code))
	}
	if len(response.TransResult) == 0 {
		return nil, fmt.Errorf("Baidu：%w", ErrInvalidResponse)
	}

	quickLookURL := "https://fanyi.baidu.com/#auto/auto/" + encodedPath(query.Text)
	results := make([]Result, 0, len(response.TransResult))
	for _, translation := range response.TransResult {
		if !nonEmpty(translation.Dst) {
			continue
		}
		pronounce := query.Text
		if query.HasHan {
			pronounce = translation.Dst
		}
		addResult(&results, Result{
			Title:        translation.Dst,
			Subtitle:     translation.Src,
			Arg:          translation.Dst,
			Pronounce:    pronounce,
			QuickLookURL: quickLookURL,
		})
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("Baidu：%w", ErrInvalidResponse)
	}
	return results, nil
}

func baiduErrorMessage(code string) string {
	messages := map[string]string{
		"54000": "缺少必填的参数",
		"58001": "不支持的语言类型",
		"54005": "翻译文本过长",
		"52003": "应用ID无效",
		"58002": "无相关服务的有效实例",
		"90107": "开发者账号无效",
		"54001": "签名检验失败,检查 KEY 和 SECRET",
		"54004": "账户已经欠费",
		"54003": "访问频率受限",
	}
	if message, ok := messages[code]; ok {
		return message
	}
	return "请参考错误码：" + code
}
