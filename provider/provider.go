package provider

import (
	"context"
	"net/http"
)

type Query struct {
	Text   string
	HasHan bool
}

type Result struct {
	Title        string
	Subtitle     string
	Arg          string
	Pronounce    string
	QuickLookURL string
}

type Provider interface {
	BuildRequest(context.Context, Query) (*http.Request, error)
	ParseResponse(Query, []byte) ([]Result, error)
}
