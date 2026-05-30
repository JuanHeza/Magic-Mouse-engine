package model

type CardBase struct {
	ID        string     `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	Version   string     `json:"version" db:"version"`
	Cost      int        `json:"cost" db:"cost"`
	Inkwell   *bool      `json:"inkwell" db:"inkwell"`
	Type      []string   `json:"type" db:"type"`
	Subtypes  []string   `json:"classifications" db:"classifications"`
	Keywords  []string   `json:"keywords" db:"keywords"`
	MoveCost  int        `json:"move_cost" db:"move_cost"`
	Strength  int        `json:"strength" db:"strength"`
	Willpower int        `json:"willpower" db:"willpower"`
	Lore      int        `json:"lore" db:"lore"`
	Images    CardImages `json:"images_uris"`
}

type CardAPI struct {
	CardBase

	Ink             string   `json:"ink" db:"ink"`
	Inks            string   `json:"inks" db:"inks"`
	Text            string   `json:"text" db:"text"`
	Lang            string   `json:"lang" db:"lang"`
	FlavorText      string   `json:"flavor_text" db:"flavor_text"`
	Rarity          string   `json:"rarity" db:"rarity"`
	Illustrator     []string `json:"illustrator" db:"illustrator"`
	CollectorNumber string   `json:"collector_number" db:"collector_number"`
	Legalities      Formats  `json:"legalities" db:"legalities"`
	Set             Set      `json:"set" db:"set"`
}

type CardSizes struct {
	Small  string `json:"small" db:"small"`
	Normal string `json:"normal" db:"normal"`
	Large  string `json:"large" db:"large"`
}

type CardState struct {
	CardBase

	Zone         string `json:"zone"`
	Owner        string `json:"owner"`
	Drying       bool   `json:"drying"`
	Exerted      bool   `json:"exerted"`
	Damage       int    `json:"damaged"`
	WillpowerMod int    `json:"willpower_mod"`
	StrengthMod  int    `json:"strength_mod"`
	LoreMod      int    `json:"lore_mod"`
	CostMod      int    `json:"cost_mod"`
	Challenger   int    `json:"challenger"`
	Resist       int    `json:"resist"`
	InLocation   bool   `json:"in_location"`
}

type CardImages struct {
	CardSizes CardSizes `json:"digital"`
}

type Formats struct {
	Standard string `json:"core"`
	Infinity string `json:"infinity"`
}

type Set struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}