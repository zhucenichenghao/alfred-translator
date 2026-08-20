package provider

import (
	"context"
	"net/url"
	"strings"
	"testing"
)

func fixedSalt() (string, error) { return "1234", nil }

func TestYoudaoRequest(t *testing.T) {
	provider := NewYoudao("key", "secret", fixedSalt)
	request, err := provider.BuildRequest(context.Background(), Query{Text: "你好！", HasHan: true})
	if err != nil {
		t.Fatal(err)
	}
	values := request.URL.Query()
	if values.Get("from") != "zh-CHS" || values.Get("to") != "en" || values.Get("q") != "你好！" {
		t.Fatalf("unexpected query: %s", request.URL.RawQuery)
	}
	if values.Get("sign") != signature("key", "你好！", "1234", "secret") {
		t.Fatal("unexpected signature")
	}
	if strings.Contains(request.URL.String(), "secret") {
		t.Fatal("secret leaked into URL")
	}
}

func TestYoudaoResponseOrderAndPhonetic(t *testing.T) {
	body := []byte(`{"errorCode":"0","translation":["你好"],"basic":{"explains":["int. 你好"],"us-phonetic":"həˈloʊ","uk-phonetic":"həˈləʊ"},"web":[{"key":"hello","value":["你好","您好"]}]}`)
	results, err := NewYoudao("key", "secret", fixedSalt).ParseResponse(Query{Text: "hello"}, body)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 4 {
		t.Fatalf("got %d results", len(results))
	}
	if results[0].Title != "你好" || results[1].Title != "int. 你好" || !strings.HasPrefix(results[2].Arg, "~") || results[3].Title != "你好, 您好" {
		t.Fatalf("unexpected result order: %#v", results)
	}
}

func TestYoudaoErrorsAndQuickLookEncoding(t *testing.T) {
	provider := NewYoudao("key", "secret", fixedSalt)
	if _, err := provider.ParseResponse(Query{Text: "hello"}, []byte(`{"errorCode":"202"}`)); err == nil || !strings.Contains(err.Error(), "签名检验失败") {
		t.Fatalf("unexpected error: %v", err)
	}
	results, err := provider.ParseResponse(Query{Text: "a/b?#c"}, []byte(`{"errorCode":"0","translation":["ok"]}`))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := url.Parse(results[0].QuickLookURL); err != nil || strings.Contains(results[0].QuickLookURL, "a/b?#c") {
		t.Fatalf("query was not encoded: %s", results[0].QuickLookURL)
	}
}

func TestMD5UsesUTF8(t *testing.T) {
	if got := md5Hex("😀"); got != "2a02eac39d716a70ecf37579185927b6" {
		t.Fatalf("unexpected UTF-8 MD5: %s", got)
	}
}
