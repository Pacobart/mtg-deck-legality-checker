package scryfall

import (
	"bytes"
	"testing"
)

func TestGetDeckList(t *testing.T) {
	deckUrl := "https://scryfall.com/@Pacobart/decks/2ca4c348-b07a-4930-8b4e-3496db97199e"
	wantString := `json = {
  doublespacedhere = things
  }
  `
	want := []byte(wantString)
	actual := GetDeckList(deckUrl)
	if bytes.Equal(want, actual) == false {
		t.Errorf("want [%s], got [%s]", want, actual)
	}
}
