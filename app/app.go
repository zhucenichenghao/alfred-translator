package app

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/zhucenichenghao/alfred-translator/alfred"
	"github.com/zhucenichenghao/alfred-translator/provider"
	"github.com/zhucenichenghao/alfred-translator/translate"
)

type ProviderFactory func(platform, key, secret string) (provider.Provider, error)

func Run(ctx context.Context, args []string, getenv func(string) string, output io.Writer, client *http.Client) int {
	return RunWithFactory(ctx, args, getenv, output, client, func(platform, key, secret string) (provider.Provider, error) {
		return provider.New(platform, key, secret, nil)
	})
}

func RunWithFactory(ctx context.Context, args []string, getenv func(string) string, output io.Writer, client *http.Client, factory ProviderFactory) int {
	workflow := execute(ctx, args, getenv, client, factory)
	if err := json.NewEncoder(output).Encode(workflow); err != nil {
		return 1
	}
	return 0
}

func execute(ctx context.Context, args []string, getenv func(string) string, client *http.Client, factory ProviderFactory) alfred.Workflow {
	config, err := LoadConfig(getenv)
	if err != nil {
		return alfred.Error(err.Error())
	}
	rawQuery, err := LastQuery(args)
	if err != nil {
		return alfred.Error(err.Error())
	}
	query, err := translate.ParseQuery(rawQuery)
	if err != nil {
		return alfred.Error(err.Error())
	}
	selectedProvider, err := factory(config.Platform, config.Key, config.Secret)
	if err != nil {
		return alfred.Error(err.Error())
	}
	results, err := translate.NewService(client).Translate(ctx, selectedProvider, query)
	if err != nil {
		return alfred.Error(err.Error())
	}
	return alfred.FromResults(results)
}
