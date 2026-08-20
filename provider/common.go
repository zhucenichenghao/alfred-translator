package provider

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"unicode"
)

var ErrInvalidResponse = errors.New("翻译服务返回无效响应")

func providerError(name, code, message string) error {
	if code == "" {
		return fmt.Errorf("%s：%s", name, message)
	}
	return fmt.Errorf("%s 错误 %s：%s", name, code, message)
}

func encodedPath(text string) string {
	return url.PathEscape(text)
}

func addResult(results *[]Result, result Result) {
	limit := 60
	if containsHan(result.Title) {
		limit = 27
	}
	runes := []rune(result.Title)
	if len(runes) > limit {
		result.Title = string(runes[:limit])
		result.Subtitle = string(runes[limit:])
	}
	*results = append(*results, result)
}

func containsHan(text string) bool {
	for _, r := range text {
		if unicode.Is(unicode.Han, r) {
			return true
		}
	}
	return false
}

func nonEmpty(value string) bool {
	return strings.TrimSpace(value) != ""
}
