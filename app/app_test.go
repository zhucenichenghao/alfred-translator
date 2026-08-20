package app

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/zhucenichenghao/alfred-translator/provider"
)

type appProvider struct{}

func (appProvider) BuildRequest(ctx context.Context, _ provider.Query) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, http.MethodGet, "https://example.com", nil)
}

func (appProvider) ParseResponse(_ provider.Query, _ []byte) ([]provider.Result, error) {
	return []provider.Result{{Title: "你好", Subtitle: "hello", Arg: "你好", Pronounce: "hello"}}, nil
}

type appRoundTripFunc func(*http.Request) (*http.Response, error)

func (f appRoundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) { return f(request) }

func TestRunOutputsAlfredJSON(t *testing.T) {
	client := &http.Client{Transport: appRoundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{}`))}, nil
	})}
	var output bytes.Buffer
	code := RunWithFactory(context.Background(), []string{"binary", "hello"}, env(map[string]string{
		"key": "k", "secret": "s", "platform": "Youdao",
	}), &output, client, func(_, _, _ string) (provider.Provider, error) { return appProvider{}, nil })
	if code != 0 {
		t.Fatalf("exit code %d", code)
	}
	var workflow map[string][]map[string]any
	if err := json.Unmarshal(output.Bytes(), &workflow); err != nil {
		t.Fatalf("invalid JSON %q: %v", output.String(), err)
	}
	if workflow["items"][0]["title"] != "你好" {
		t.Fatalf("unexpected output: %s", output.String())
	}
}

func TestRunTurnsConfigErrorIntoItem(t *testing.T) {
	var output bytes.Buffer
	code := Run(context.Background(), []string{"binary"}, env(nil), &output, http.DefaultClient)
	if code != 0 {
		t.Fatalf("exit code %d", code)
	}
	if !strings.Contains(output.String(), `"valid":false`) || !strings.Contains(output.String(), "缺少环境变量 key") {
		t.Fatalf("unexpected output: %s", output.String())
	}
}
