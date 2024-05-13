package scryfall

import (
	"testing"
)

func TestGetDeckList(t *testing.T) {
	deckUrl := "https://scryfall.com/@Pacobart/decks/2ca4c348-b07a-4930-8b4e-3496db97199e"
	want := "Judith, Carnage Connoisseur Spellslinger"
	actual := GetDeckList(deckUrl).Name
	if want != actual {
		t.Errorf("want [%s], got [%s]", want, actual)
	}
}
