package internal

type keyword string

var Keywords = []keyword{
	"artifact",
	"block",
	"dexterity",
	"draw",
	"gain",
	"lose",
	"poison",
	"strength",
	"vulnerable",
	"vulnerable",
	"frail",
	"weak",
	"attack card",
	"curse",
	"ethereal",
	"exhaust",
	"innate",
	"power card",
	"retain",
	"skill card",
	"status",
	"unplayable",
	"channel",
	"evoke",
	"orb",
	"stance",
	"turn",
	"combat",
	"act",
	"run",
	"gain",
	"lose",
	"orb",
	"lightning",
	"frost",
	"dark",
	"plasma",
	"focus",
	"heal",
}

type Card struct {
	Title					string
	EnergyCost 		string
	Description		string
	Color 				string
	Character			string
	Type					string
	Rarity				string
	Source				string
	Img						string
	URL						string
}

type Relic struct {
	Name 					string
	Description		string
	Character			string
	Rarity				string
	Img						string
}

type Potion struct {
	Name 					string
	Description		string
	Character			string
	Rarity				string
	Img						string
}
