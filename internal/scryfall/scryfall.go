package scryfall

import (
	"fmt"
	"regexp"
)

func GetDeckListFromScryfall(url string) {
	deckRegex, _ := regexp.Compile("[0-9A-Fa-f]{8}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{12}")
	deckId := deckRegex.FindString(url)
	fullUrl := fmt.Sprintf("https://api.scryfall.com/decks/%s/export/json", deckId)
	fmt.Print(fullUrl)
}
