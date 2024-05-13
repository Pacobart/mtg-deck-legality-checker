package main

import (
	"fmt"

	"github.com/Pacobart/mtg-deck-legality-checker/internal/helpers"
	"github.com/Pacobart/mtg-deck-legality-checker/internal/scryfall"
)

func main() {
	//helpers.DEBUG = true

	deckUrl := "https://scryfall.com/@Pacobart/decks/2ca4c348-b07a-4930-8b4e-3496db97199e"
	deckList := scryfall.GetDeckList(deckUrl)
	helpers.Debug(fmt.Sprintf("Decklist name: %s", deckList.Name))
	cards := deckList.GetCards()
	areCardsLegalForCommander := scryfall.CheckCardsForFormatLegal("Commander", cards)
	for _, card := range areCardsLegalForCommander {
		name := card.Name
		legal := card.Legal
		fmt.Println(fmt.Sprintf("Card Legality for %s is %t", name, legal))
	}
	// var cardNames []string
	// for _, card := range cards {
	// 	name := card.CardDigest.Name
	// 	cardNames = append(cardNames, name)
	// }
	// fmt.Println(len(cardNames))
	// fmt.Print(cardNames)
}
