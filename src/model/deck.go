package model

type Deck struct {
	ID    string     `json:"id"`
	Name  string     `json:"name"`
	Cards []CardBase `json:"cards"`
}

type DeckAPI struct {
	Name  string `json:"name"`
	Details []string `json:"details"`
	Cards []CardBase `json:"cards"`
}
// El formato por linea es `${cantidad} ${name} - ${version}` 