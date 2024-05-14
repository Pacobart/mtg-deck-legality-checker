package main

import (
	"fmt"

	"github.com/Pacobart/mtg-deck-legality-checker/internal/helpers"
	"github.com/Pacobart/mtg-deck-legality-checker/internal/scryfall"
)

func performLogic(deckUrl string, deckFormat string) string {
	deckList := scryfall.GetDeckList(deckUrl)
	helpers.Debug(fmt.Sprintf("Decklist name: %s", deckList.Name))
	cards := deckList.GetCards()
	areCardsLegalForMyFormat := scryfall.CheckCardsForFormatLegal(deckFormat, cards)

	var legalCards []string
	var illegalCards []string
	for _, card := range areCardsLegalForMyFormat {
		if card.Legal {
			legalCards = append(legalCards, card.Name)
		} else {
			illegalCards = append(illegalCards, card.Name)
		}
	}

	result := struct {
		DeckName         string
		Format           string
		LegalCardCount   int
		IllegalCardCount int
		LegalCards       []string
		IllegalCards     []string
	}{
		DeckName:         deckList.Name,
		Format:           deckFormat,
		LegalCardCount:   len(legalCards),
		IllegalCardCount: len(illegalCards),
		LegalCards:       legalCards,
		IllegalCards:     illegalCards,
	}

	resultJson := helpers.StructToJson(result)
	return resultJson
}

func main() {
	//helpers.DEBUG = true

	//deckUrl := "https://scryfall.com/@Pacobart/decks/2ca4c348-b07a-4930-8b4e-3496db97199e" //Commander
	//deckFormat := "Commander"
	deckUrl := "https://scryfall.com/@Pacobart/decks/a8305a46-ab83-4897-957e-2ef128a67e10"
	deckFormat := "Standard"
	result := performLogic(deckUrl, deckFormat)
	fmt.Println(result)
}
