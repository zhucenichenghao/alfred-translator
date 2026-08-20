package translate

import "testing"

func TestNormalizeCamelCase(t *testing.T) {
	tests := map[string]string{
		"helloWorld":   "hello world",
		"HelloWorld":   "hello world",
		"HTTPServer":   "h t t p server",
		" already ok ": "already ok",
	}
	for input, want := range tests {
		if got := NormalizeCamelCase(input); got != want {
			t.Fatalf("NormalizeCamelCase(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestParseQueryDetectsHan(t *testing.T) {
	for _, input := range []string{"你好", "你好！", "hello中文"} {
		query, err := ParseQuery(input)
		if err != nil {
			t.Fatal(err)
		}
		if !query.HasHan {
			t.Fatalf("expected %q to contain Han characters", input)
		}
	}
	query, err := ParseQuery("hello!")
	if err != nil {
		t.Fatal(err)
	}
	if query.HasHan {
		t.Fatal("did not expect ASCII query to contain Han characters")
	}
}

func TestParseQueryRejectsWhitespace(t *testing.T) {
	if _, err := ParseQuery("  \t "); err == nil {
		t.Fatal("expected an error")
	}
}
