package scryfall

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDeckList(t *testing.T) {
	deckUrl := "https://scryfall.com/@Pacobart/decks/2ca4c348-b07a-4930-8b4e-3496db97199e"
	want := "Judith, Carnage Connoisseur Spellslinger"
	actual := GetDeckList(deckUrl).Name
	if want != actual {
		t.Errorf("want [%s], got [%s]", want, actual)
	}
}

func TestGetCards(t *testing.T) {
	jsonData := `{
		"object": "deck",
		"id": "a8305a46-ab83-4897-957e-2ef128a67e10",
		"name": "Judith rakdos standard",
		"format": "standard",
		"layout": "constructed",
		"uri": "https://api.scryfall.com/decks/a8305a46-ab83-4897-957e-2ef128a67e10",
		"scryfall_uri": "https://scryfall.com/@Pacobart/decks/a8305a46-ab83-4897-957e-2ef128a67e10",
		"description": "possible: swap 2 chimil for 2 judith",
		"trashed": false,
		"in_compliance": true,
		"sections": {
			"primary": [
				"mainboard"
			],
			"secondary": [
				"sideboard",
				"maybeboard"
			]
		},
		"entries": {
			"mainboard": [
				{
					"object": "deck_entry",
					"id": "004c5972-19d4-4daa-a2b5-ee425da944b6",
					"deck_id": "a8305a46-ab83-4897-957e-2ef128a67e10",
					"section": "mainboard",
					"count": 3,
					"raw_text": "3 Judith, Carnage Connoisseur",
					"card_digest": {
						"object": "card_digest",
						"id": "3eaa19ce-cace-499e-8b23-ef9e56b23700",
						"oracle_id": "e86c965b-ac10-4fff-b682-dddd6d9747c6",
						"name": "Judith, Carnage Connoisseur",
						"scryfall_uri": "https://scryfall.com/card/mkm/210/judith-carnage-connoisseur"
					}
				}
			]
		}
	}`

	var deck DeckList
	err := json.Unmarshal([]byte(jsonData), &deck)
	assert.NoError(t, err, "Unmarshal should not return an error")
	cards := deck.GetCards()
	want := "Judith, Carnage Connoisseur"
	actual := cards[0].CardDigest.Name
	assert.Equal(t, want, actual, "Card digest name should match")
}

func TestCheckCardsForFormatLegal(t *testing.T) {
	jsonData := `{
		"object": "deck",
		"id": "a8305a46-ab83-4897-957e-2ef128a67e10",
		"name": "Judith rakdos standard",
		"format": "standard",
		"layout": "constructed",
		"uri": "https://api.scryfall.com/decks/a8305a46-ab83-4897-957e-2ef128a67e10",
		"scryfall_uri": "https://scryfall.com/@Pacobart/decks/a8305a46-ab83-4897-957e-2ef128a67e10",
		"description": "possible: swap 2 chimil for 2 judith",
		"trashed": false,
		"in_compliance": true,
		"sections": {
			"primary": [
				"mainboard"
			],
			"secondary": [
				"sideboard",
				"maybeboard"
			]
		},
		"entries": {
			"mainboard": [
				{
					"object": "deck_entry",
					"id": "004c5972-19d4-4daa-a2b5-ee425da944b6",
					"deck_id": "a8305a46-ab83-4897-957e-2ef128a67e10",
					"section": "mainboard",
					"count": 3,
					"raw_text": "3 Judith, Carnage Connoisseur",
					"card_digest": {
						"object": "card_digest",
						"id": "3eaa19ce-cace-499e-8b23-ef9e56b23700",
						"oracle_id": "e86c965b-ac10-4fff-b682-dddd6d9747c6",
						"name": "Judith, Carnage Connoisseur",
						"scryfall_uri": "https://scryfall.com/card/mkm/210/judith-carnage-connoisseur"
					}
				}
			]
		}
	}`

	var deck DeckList
	err := json.Unmarshal([]byte(jsonData), &deck)
	assert.NoError(t, err, "Unmarshal should not return an error")
	cards := deck.GetCards()
	areCardsLegalForMyFormat := CheckCardsForFormatLegal("Commander", cards)

	assert.Equal(t, "Judith, Carnage Connoisseur", areCardsLegalForMyFormat[0].Name, "Card digest name should match")
	assert.Equal(t, true, areCardsLegalForMyFormat[0].Legal, "Card legal should be true")
}

func TestGetCard(t *testing.T) {
	cardId := "3eaa19ce-cace-499e-8b23-ef9e56b23700"
	card := GetCard(string(cardId))

	want := "Judith, Carnage Connoisseur"
	actual := card.Name
	assert.Equal(t, want, actual, "Card names should match")

}
func TestCheckCardForFormatLegal(t *testing.T) {
	jsonData := `{
        "object": "card",
        "id": "3eaa19ce-cace-499e-8b23-ef9e56b23700",
        "name": "Judith, Carnage Connoisseur",
        "released_at": "2024-02-09",
        "uri": "https://api.scryfall.com/cards/3eaa19ce-cace-499e-8b23-ef9e56b23700",
        "scryfall_uri": "https://scryfall.com/card/mkm/210/judith-carnage-connoisseur?utm_source=api",
        "legalities": {
            "standard": "legal",
            "future": "legal",
            "historic": "legal",
            "timeless": "legal",
            "gladiator": "legal",
            "pioneer": "legal",
            "explorer": "legal",
            "modern": "legal",
            "legacy": "legal",
            "pauper": "not_legal",
            "vintage": "legal",
            "penny": "legal",
            "commander": "legal",
            "oathbreaker": "legal",
            "standardbrawl": "legal",
            "brawl": "legal",
            "alchemy": "legal",
            "paupercommander": "not_legal",
            "duel": "legal",
            "oldschool": "not_legal",
            "premodern": "not_legal",
            "predh": "not_legal"
        },
        "rarity": "rare"
    }`

	// Unmarshal JSON data into the Card struct
	var card Card
	err := json.Unmarshal([]byte(jsonData), &card)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var cardentry CardEntry
	cardLegality := cardentry.CheckCardForFormatLegal("Commander", card)
	var cardLegalities []CardLegalForFormat
	cardLegalities = append(cardLegalities, cardLegality)

	want := "Judith, Carnage Connoisseur"
	actual := cardLegalities[0].Name
	assert.Equal(t, want, actual, "Card names should match")
}

func TestIsCardLegalForFormatTrue(t *testing.T) {
	jsonData := `{
        "id": "3eaa19ce-cace-499e-8b23-ef9e56b23700",
        "name": "Judith, Carnage Connoisseur",
        "legalities": {
            "standard": "legal",
            "commander": "legal",
            "premodern": "not_legal",
            "predh": "not_legal"
        }
    }`

	// Unmarshal JSON data into the Card struct
	var card Card
	err := json.Unmarshal([]byte(jsonData), &card)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	formatLegal := card.IsCardLegalForFormat("Commander")

	want := true
	actual := formatLegal
	assert.Equal(t, want, actual, "Card legal should be true")
}

func TestIsCardLegalForFormatFalse(t *testing.T) {
	jsonData := `{
        "id": "3eaa19ce-cace-499e-8b23-ef9e56b23700",
        "name": "Judith, Carnage Connoisseur",
        "legalities": {
            "standard": "legal",
            "commander": "legal",
            "premodern": "not_legal",
            "predh": "not_legal"
        }
    }`

	// Unmarshal JSON data into the Card struct
	var card Card
	err := json.Unmarshal([]byte(jsonData), &card)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	formatLegal := card.IsCardLegalForFormat("Premodern")

	want := false
	actual := formatLegal
	assert.Equal(t, want, actual, "Card legal should be false")
}

func TestIsCardLegalForFormatInvalid(t *testing.T) {
	jsonData := `{
        "id": "3eaa19ce-cace-499e-8b23-ef9e56b23700",
        "name": "Judith, Carnage Connoisseur",
        "legalities": {
            "standard": "legal",
            "commander": "legal",
            "premodern": "not_legal",
            "predh": "not_legal"
        }
    }`

	// Unmarshal JSON data into the Card struct
	var card Card
	err := json.Unmarshal([]byte(jsonData), &card)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	formatLegal := card.IsCardLegalForFormat("Nonexistent") //invalid format

	want := false
	actual := formatLegal
	assert.Equal(t, want, actual, "Card legal should be false")
}
