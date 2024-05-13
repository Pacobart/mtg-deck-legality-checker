package main

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
func main() {
	deckUrl := "https://scryfall.com/@Pacobart/decks/2ca4c348-b07a-4930-8b4e-3496db97199e"
	GetDeckListFromScryfall(deckUrl)
}
