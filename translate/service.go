package translate

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/zhucenichenghao/alfred-translator/provider"
)

const maxResponseSize = 2 << 20

type HTTPDoer interface {
	Do(*http.Request) (*http.Response, error)
}

type Service struct {
	client HTTPDoer
}

func NewService(client HTTPDoer) *Service {
	return &Service{client: client}
}

func (s *Service) Translate(ctx context.Context, p provider.Provider, query provider.Query) ([]provider.Result, error) {
	request, err := p.BuildRequest(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("构造翻译请求失败：%w", err)
	}
	response, err := s.client.Do(request)
	if err != nil {
		if ctx.Err() != nil {
			return nil, fmt.Errorf("翻译请求超时或已取消：%w", ctx.Err())
		}
		return nil, fmt.Errorf("翻译网络请求失败：%w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("翻译服务返回 HTTP %d", response.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(response.Body, maxResponseSize+1))
	if err != nil {
		return nil, fmt.Errorf("读取翻译响应失败：%w", err)
	}
	if len(body) > maxResponseSize {
		return nil, fmt.Errorf("翻译响应过大")
	}
	results, err := p.ParseResponse(query, body)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, provider.ErrInvalidResponse
	}
	return results, nil
}
