package main

import "github.com/Pacobart/mtg-deck-legality-checker/internal/scryfall"

func main() {
	deckUrl := "https://scryfall.com/@Pacobart/decks/2ca4c348-b07a-4930-8b4e-3496db97199e"
	scryfall.GetDeckList(deckUrl)
}
