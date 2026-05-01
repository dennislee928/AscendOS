package httpapi

import "strings"

// ModuleCatalogEntry describes one AscendOS module exposed by the gateway.
type ModuleCatalogEntry struct {
	Name        string `json:"name"`
	Key         string `json:"key"`
	Language    string `json:"language"`
	Domain      string `json:"domain"`
	Summary     string `json:"summary"`
	ServicePath string `json:"servicePath"`
}

type moduleCatalog struct {
	items []ModuleCatalogEntry
	byKey map[string]ModuleCatalogEntry
}

func newModuleCatalog() moduleCatalog {
	items := []ModuleCatalogEntry{
		{
			Name:        "chronos",
			Key:         "CHRONOS",
			Language:    "Go",
			Domain:      "Sleep science",
			Summary:     "Circadian tracking, sleep-stage logging, and light-exposure planning.",
			ServicePath: "services/chronos",
		},
		{
			Name:        "aegis",
			Key:         "AEGIS",
			Language:    "Rust",
			Domain:      "Psychological defense",
			Summary:     "Manipulation-pattern detection for chat logs and transcripts.",
			ServicePath: "services/aegis",
		},
		{
			Name:        "neuro",
			Key:         "NEURO",
			Language:    "Python",
			Domain:      "Neuroscience",
			Summary:     "Mood and neurotransmitter journaling with retrieval-backed explanations.",
			ServicePath: "services/neuro",
		},
		{
			Name:        "orator",
			Key:         "ORATOR",
			Language:    "Rust",
			Domain:      "Communication",
			Summary:     "Speech recording, prosody analysis, and narrative-frame scoring.",
			ServicePath: "services/orator",
		},
		{
			Name:        "metis",
			Key:         "METIS",
			Language:    "Go",
			Domain:      "Cognition",
			Summary:     "Spaced-repetition workflows and bias-detection prompts.",
			ServicePath: "services/metis",
		},
		{
			Name:        "argentum",
			Key:         "ARGENTUM",
			Language:    "Python",
			Domain:      "Finance",
			Summary:     "Cashflow tracking and behavioural-finance forecasting.",
			ServicePath: "services/argentum",
		},
		{
			Name:        "kairos",
			Key:         "KAIROS",
			Language:    "Python",
			Domain:      "Learning",
			Summary:     "VARK profiling and personalised study-path generation.",
			ServicePath: "services/kairos",
		},
		{
			Name:        "praxis",
			Key:         "PRAXIS",
			Language:    "Go",
			Domain:      "Habits",
			Summary:     "Habit graphing, streak tracking, and relapse forecasting.",
			ServicePath: "services/praxis",
		},
	}

	byKey := make(map[string]ModuleCatalogEntry, len(items))
	for _, item := range items {
		byKey[strings.ToLower(item.Name)] = item
	}

	return moduleCatalog{
		items: items,
		byKey: byKey,
	}
}

func (c moduleCatalog) all() []ModuleCatalogEntry {
	out := make([]ModuleCatalogEntry, len(c.items))
	copy(out, c.items)
	return out
}

func (c moduleCatalog) get(name string) (ModuleCatalogEntry, bool) {
	item, ok := c.byKey[strings.ToLower(strings.TrimSpace(name))]
	return item, ok
}
