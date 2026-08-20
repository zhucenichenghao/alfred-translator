package provider

import (
	"context"
	"strings"
	"testing"
)

func TestBaiduRequest(t *testing.T) {
	provider := NewBaidu("appid", "secret", fixedSalt)
	request, err := provider.BuildRequest(context.Background(), Query{Text: "hello", HasHan: false})
	if err != nil {
		t.Fatal(err)
	}
	values := request.URL.Query()
	if values.Get("from") != "auto" || values.Get("to") != "zh" || values.Get("dict") != "1" || values.Get("action") != "1" {
		t.Fatalf("unexpected query: %s", request.URL.RawQuery)
	}
}

func TestBaiduUsesTranslationAsArg(t *testing.T) {
	results, err := NewBaidu("appid", "secret", fixedSalt).ParseResponse(Query{Text: "hello"}, []byte(`{"trans_result":[{"src":"hello","dst":"你好"}]}`))
	if err != nil {
		t.Fatal(err)
	}
	if results[0].Arg != "你好" || results[0].Pronounce != "hello" {
		t.Fatalf("unexpected result: %#v", results[0])
	}
}

func TestBaiduAcceptsNumericErrorCode(t *testing.T) {
	_, err := NewBaidu("appid", "secret", fixedSalt).ParseResponse(Query{Text: "hello"}, []byte(`{"error_code":54001}`))
	if err == nil || !strings.Contains(err.Error(), "签名检验失败") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBaiduRejectsEmptyResponse(t *testing.T) {
	if _, err := NewBaidu("appid", "secret", fixedSalt).ParseResponse(Query{Text: "hello"}, []byte(`{}`)); err == nil {
		t.Fatal("expected invalid response error")
	}
}
