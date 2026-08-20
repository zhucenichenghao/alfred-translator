package alfred

import (
	"encoding/json"
	"testing"

	"github.com/zhucenichenghao/alfred-translator/provider"
)

func TestFromResults(t *testing.T) {
	workflow := FromResults([]provider.Result{{
		Title: "你好", Subtitle: "hello", Arg: "你好", Pronounce: "hello", QuickLookURL: "https://example.com",
	}, {
		Title: "[həˈləʊ]", Subtitle: "回车可听发音", Arg: "~hello", Pronounce: "hello",
	}})
	if len(workflow.Items) != 2 {
		t.Fatalf("got %d items", len(workflow.Items))
	}
	if workflow.Items[0].Icon.Path != translateIcon || workflow.Items[1].Icon.Path != pronounceIcon {
		t.Fatal("unexpected icon mapping")
	}
	if workflow.Items[0].Mods["cmd"].Arg != "hello" || workflow.Items[0].Text.Copy != "你好" {
		t.Fatal("modifier or copy text was not preserved")
	}
}

func TestErrorSerializesValidFalse(t *testing.T) {
	data, err := json.Marshal(Error("bad config"))
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string][]map[string]any
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	item := decoded["items"][0]
	if valid, ok := item["valid"].(bool); !ok || valid {
		t.Fatalf("valid=false missing from %s", data)
	}
	if _, ok := item["mods"]; ok {
		t.Fatal("error item should not contain modifiers")
	}
}
