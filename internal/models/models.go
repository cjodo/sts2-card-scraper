package models

type Card struct {
	URL         string
	Title       string
	Description string
	Character   string
	Type        string
	EnergyCost  string
	Rarity      string
	Img         string
}

type Relic struct {
	Name        string
	Description string
	Character   string
	Rarity      string
	Img         string
}
