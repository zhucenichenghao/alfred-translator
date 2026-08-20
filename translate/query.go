package translate

import (
	"errors"
	"strings"
	"unicode"

	"github.com/zhucenichenghao/alfred-translator/provider"
)

func ParseQuery(raw string) (provider.Query, error) {
	text := NormalizeCamelCase(raw)
	if text == "" {
		return provider.Query{}, errors.New("请输入需要翻译的内容")
	}
	return provider.Query{Text: text, HasHan: ContainsHan(text)}, nil
}

func NormalizeCamelCase(raw string) string {
	raw = strings.TrimSpace(raw)
	var builder strings.Builder
	for _, r := range raw {
		if r >= 'A' && r <= 'Z' {
			builder.WriteByte(' ')
		}
		builder.WriteRune(unicode.ToLower(r))
	}
	return strings.TrimSpace(builder.String())
}

func ContainsHan(text string) bool {
	for _, r := range text {
		if unicode.Is(unicode.Han, r) {
			return true
		}
	}
	return false
}
