package alfred

type Workflow struct {
	Items []Item `json:"items"`
}

type Item struct {
	Title        string              `json:"title"`
	Subtitle     string              `json:"subtitle"`
	Arg          string              `json:"arg,omitempty"`
	Valid        bool                `json:"valid"`
	Icon         *Icon               `json:"icon,omitempty"`
	Mods         map[string]Modifier `json:"mods,omitempty"`
	Text         *Text               `json:"text,omitempty"`
	QuickLookURL string              `json:"quicklookurl,omitempty"`
}

type Icon struct {
	Path string `json:"path"`
}

type Modifier struct {
	Subtitle string `json:"subtitle"`
	Arg      string `json:"arg"`
	Valid    bool   `json:"valid"`
}

type Text struct {
	Copy string `json:"copy"`
}
