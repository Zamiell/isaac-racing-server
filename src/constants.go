package main

type AllItems struct {
	Name string
}

// JSONItem is a struct for reading in all the json items
type JSONItem struct {
	Name        string  `json:"name"`
	Damage      float32 `json:"dmg"`
	DamageX     float32 `json:"dmg_x"`
	Health      int     `json:"health"`
	SoulHearts  int     `json:"soul_hearts"`
	SinHearts   int     `json:"sin_hearts"`
	Tears       float32 `json:"tears"`
	Delay       float32 `json:"delay"`
	DelayX      float32 `json:"delay_x`
	Speed       float32 `json:"speed"`
	ShotSpeed   float32 `json:"shot_speed"`
	Height      float32 `json:"height"`
	Range       float32 `json:"range"`
	Luck        float32 `json:"luck"`
	Beelzebub   bool    `json:"beelzebub`
	Bob         bool    `json:"bob"`
	Bookworm    bool    `json:"bookworm"`
	Conjoined   bool    `json:"conjoined"`
	Funguy      bool    `json:"funguy"`
	Guppy       bool    `json:"guppy"`
	Leviathan   bool    `json:"leviathan"`
	OhCrap      bool    `json:"ohcrap"`
	Seraphim    bool    `json:"seraphim"`
	SpiderBaby  bool    `json:"spiderbaby"`
	Spun        bool    `json:"spun"`
	Superbum    bool    `json:"superbum"`
	YesMother   bool    `json:"yesmother"`
	SpaceBar    bool    `json:"space"`
	HealthOnly  bool    `json:"health_only"`
	Intro       string  `json:"introduced_in"`
	Shown       bool    `json:"shown"`
	Summary     bool    `json:"in_summary"`
	SummaryName string  `json:"summary_name"`
	SummaryCond string  `json:"condition_name"`
	Text        string  `json:"text"`
}
