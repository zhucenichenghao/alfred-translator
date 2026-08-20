package provider

import "fmt"

func New(platform, key, secret string, saltSource SaltSource) (Provider, error) {
	if saltSource == nil {
		saltSource = randomSalt
	}
	switch platform {
	case "Youdao":
		return NewYoudao(key, secret, saltSource), nil
	case "Baidu":
		return NewBaidu(key, secret, saltSource), nil
	default:
		return nil, fmt.Errorf("platform 必须为 Youdao 或 Baidu")
	}
}
