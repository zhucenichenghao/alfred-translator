package app

import "testing"

func env(values map[string]string) func(string) string {
	return func(key string) string { return values[key] }
}

func TestLoadConfig(t *testing.T) {
	config, err := LoadConfig(env(map[string]string{"key": "k", "secret": "s", "platform": "Youdao"}))
	if err != nil {
		t.Fatal(err)
	}
	if config.Platform != "Youdao" {
		t.Fatalf("unexpected platform %q", config.Platform)
	}
}

func TestLoadConfigRejectsMissingAndUnknownValues(t *testing.T) {
	cases := []map[string]string{
		{"secret": "s", "platform": "Youdao"},
		{"key": "k", "platform": "Youdao"},
		{"key": "k", "secret": "s", "platform": "youdao"},
	}
	for _, values := range cases {
		if _, err := LoadConfig(env(values)); err == nil {
			t.Fatalf("expected error for %#v", values)
		}
	}
}

func TestLastQueryUsesLastArgument(t *testing.T) {
	query, err := LastQuery([]string{"binary", "ignored", "final query"})
	if err != nil {
		t.Fatal(err)
	}
	if query != "final query" {
		t.Fatalf("got %q", query)
	}
	if _, err := LastQuery([]string{"binary"}); err == nil {
		t.Fatal("expected missing query error")
	}
}
