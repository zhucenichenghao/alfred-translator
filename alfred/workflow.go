package alfred

import "github.com/zhucenichenghao/alfred-translator/provider"

const (
	translateIcon = "assets/translate.png"
	pronounceIcon = "assets/translate-say.png"
)

func FromResults(results []provider.Result) Workflow {
	items := make([]Item, 0, len(results))
	for _, result := range results {
		icon := translateIcon
		if len(result.Arg) > 0 && result.Arg[0] == '~' {
			icon = pronounceIcon
		}

		items = append(items, Item{
			Title:    result.Title,
			Subtitle: result.Subtitle,
			Arg:      result.Arg,
			Valid:    true,
			Icon:     &Icon{Path: icon},
			Mods: map[string]Modifier{
				"cmd": {Subtitle: "🔊 " + result.Pronounce, Arg: result.Pronounce, Valid: true},
				"alt": {Subtitle: "📣 " + result.Pronounce, Arg: result.Pronounce, Valid: true},
			},
			Text:         &Text{Copy: result.Title},
			QuickLookURL: result.QuickLookURL,
		})
	}
	return Workflow{Items: items}
}

func Error(message string) Workflow {
	return Workflow{Items: []Item{{
		Title:    "👻 翻译出错啦",
		Subtitle: message,
		Arg:      "Ooops...",
		Valid:    false,
		Icon:     &Icon{Path: translateIcon},
	}}}
}
