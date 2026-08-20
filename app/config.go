package app

import (
	"errors"
	"strings"
)

type Config struct {
	Key      string
	Secret   string
	Platform string
}

func LoadConfig(getenv func(string) string) (Config, error) {
	config := Config{
		Key:      getenv("key"),
		Secret:   getenv("secret"),
		Platform: strings.TrimSpace(getenv("platform")),
	}
	if strings.TrimSpace(config.Key) == "" {
		return Config{}, errors.New("缺少环境变量 key")
	}
	if strings.TrimSpace(config.Secret) == "" {
		return Config{}, errors.New("缺少环境变量 secret")
	}
	if config.Platform != "Youdao" && config.Platform != "Baidu" {
		return Config{}, errors.New("platform 必须为 Youdao 或 Baidu")
	}
	return config, nil
}

func LastQuery(args []string) (string, error) {
	if len(args) < 2 {
		return "", errors.New("请输入需要翻译的内容")
	}
	query := args[len(args)-1]
	if strings.TrimSpace(query) == "" {
		return "", errors.New("请输入需要翻译的内容")
	}
	return query, nil
}
