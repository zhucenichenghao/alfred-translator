package translate

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/zhucenichenghao/alfred-translator/provider"
)

type fakeProvider struct{}

func (fakeProvider) BuildRequest(ctx context.Context, _ provider.Query) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, http.MethodGet, "https://example.com", nil)
}

func (fakeProvider) ParseResponse(_ provider.Query, body []byte) ([]provider.Result, error) {
	if string(body) != "ok" {
		return nil, errors.New("bad body")
	}
	return []provider.Result{{Title: "result"}}, nil
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) { return f(request) }

func TestServiceSuccessAndHTTPError(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("ok"))}, nil
	})}
	results, err := NewService(client).Translate(context.Background(), fakeProvider{}, provider.Query{Text: "hello"})
	if err != nil || len(results) != 1 {
		t.Fatalf("results=%#v err=%v", results, err)
	}

	client.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 500, Body: io.NopCloser(strings.NewReader("fail"))}, nil
	})
	if _, err := NewService(client).Translate(context.Background(), fakeProvider{}, provider.Query{}); err == nil || !strings.Contains(err.Error(), "HTTP 500") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestServiceRejectsOversizedResponse(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(strings.Repeat("x", maxResponseSize+1)))}, nil
	})}
	if _, err := NewService(client).Translate(context.Background(), fakeProvider{}, provider.Query{}); err == nil || !strings.Contains(err.Error(), "过大") {
		t.Fatalf("unexpected error: %v", err)
	}
}
