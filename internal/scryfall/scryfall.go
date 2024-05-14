package scryfall

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"regexp"
	"sync"

	"github.com/Pacobart/mtg-deck-legality-checker/internal/helpers"
)

// CardEntry from DeckList
type CardEntry struct {
	Object     string     `json:"object"`
	ID         string     `json:"id"`
	DeckID     string     `json:"deck_id"`
	Count      int        `json:"count"`
	CardDigest CardDigest `json:"card_digest"`
}

// CardDigest from DeckList/CardEntry
type CardDigest struct {
	Object      string `json:"object"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	ScryfallURI string `json:"scryfall_uri"`
}

// Full DeckList
type DeckList struct {
	Object       string `json:"object"`
	ID           string `json:"id"`
	Name         string `json:"name"`
	Format       string `json:"format"`
	Layout       string `json:"layout"`
	URI          string `json:"uri"`
	ScryfallURI  string `json:"scryfall_uri"`
	Description  string `json:"description"`
	Trashed      bool   `json:"trashed"`
	InCompliance bool   `json:"in_compliance"`
	Entries      struct {
		Maybeboard []CardEntry `json:"maybeboard"`
		Nonlands   []CardEntry `json:"nonlands"`
		Outside    []CardEntry `json:"outside"`
		Commanders []CardEntry `json:"commanders"`
		Lands      []CardEntry `json:"lands"`
		Mainboard  []CardEntry `json:"mainboard"`
	} `json:"entries"`
}

// Individual Card api
type Card struct {
	Object      string `json:"object"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	ReleasedAt  string `json:"released_at"`
	URI         string `json:"uri"`
	ScryfallURI string `json:"scryfall_uri"`
	Legalities  struct {
		Standard        string `json:"standard"`
		Future          string `json:"future"`
		Historic        string `json:"historic"`
		Timeless        string `json:"timeless"`
		Gladiator       string `json:"gladiator"`
		Pioneer         string `json:"pioneer"`
		Explorer        string `json:"explorer"`
		Modern          string `json:"modern"`
		Legacy          string `json:"legacy"`
		Pauper          string `json:"pauper"`
		Vintage         string `json:"vintage"`
		Penny           string `json:"penny"`
		Commander       string `json:"commander"`
		Oathbreaker     string `json:"oathbreaker"`
		Standardbrawl   string `json:"standardbrawl"`
		Brawl           string `json:"brawl"`
		Alchemy         string `json:"alchemy"`
		Paupercommander string `json:"paupercommander"`
		Duel            string `json:"duel"`
		Oldschool       string `json:"oldschool"`
		Premodern       string `json:"premodern"`
		Predh           string `json:"predh"`
	} `json:"legalities"`
	Rarity string `json:"rarity"`
}

// Simplified Struct with Card data
type CardLegalForFormat struct {
	ID    string
	Name  string
	Legal bool
}

func GetDeckList(url string) DeckList {
	// Get Decklist from Scryfall and load into DeckList struct
	deckRegex, _ := regexp.Compile("[0-9A-Fa-f]{8}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{12}")
	deckId := deckRegex.FindString(url)
	fullUrl := fmt.Sprintf("https://api.scryfall.com/decks/%s/export/json", deckId)

	resp, err := http.Get(fullUrl)
	helpers.Check(err)
	if err != nil {
		helpers.Debug("No response from request")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	helpers.Check(err)
	helpers.Debug(string(body))

	var result DeckList
	if err := json.Unmarshal(body, &result); err != nil {
		helpers.Debug("Can not unmarshal JSON")
	}
	helpers.Check(err)
	return result
}

func (d *DeckList) GetCards() []CardEntry {
	// Get Cards from Decklist Struct and merge them into a single Array.
	categoryMap := map[string][]CardEntry{
		"Commanders": d.Entries.Commanders,
		"Lands":      d.Entries.Lands,
		"Maybeboard": d.Entries.Maybeboard,
		"Nonlands":   d.Entries.Nonlands,
		"Outside":    d.Entries.Outside,
		"Mainboard":  d.Entries.Mainboard,
	}

	var cards []CardEntry

	for _, entries := range categoryMap {
		cards = append(cards, entries...)
	}

	helpers.Debug(fmt.Sprintf("Cards in Deck: %v\nCount: %d\n", cards, len(cards)))

	nonEmptyCards := cards[:0]
	for _, card := range cards {
		if card.CardDigest.Name != "" {
			nonEmptyCards = append(nonEmptyCards, card)
		}
	}

	fmt.Println("--------------")
	fmt.Println(nonEmptyCards)
	fmt.Println("--------------")

	return nonEmptyCards
}

func CheckCardsForFormatLegal(format string, cards []CardEntry) []CardLegalForFormat {
	// Pass a list of Cards in to iterate over and see if card is valid for specified format
	// Uses go routine to speed up results
	resultChan := make(chan CardLegalForFormat, len(cards))
	var wg sync.WaitGroup

	for _, cardentry := range cards {
		wg.Add(1)
		go func(cardentry CardEntry) {
			defer wg.Done()
			cardId := cardentry.CardDigest.ID
			card := GetCard(string(cardId))
			cardLegality := cardentry.CheckCardForFormatLegal(format, card)
			resultChan <- cardLegality
		}(cardentry)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var cardLegalities []CardLegalForFormat
	for cardLegality := range resultChan {
		cardLegalities = append(cardLegalities, cardLegality)
	}

	return cardLegalities
}

func (c *CardEntry) CheckCardForFormatLegal(format string, card Card) CardLegalForFormat {
	// Pass a card in and validate its legality in the specificed game format. Returns struct with card info and legality
	cardId := card.ID
	helpers.Debug(fmt.Sprintf("looking up Card info from: %s", string(cardId)))
	formatLegal := card.IsCardLegalForFormat(format)
	helpers.Debug(fmt.Sprintf("Card Legality for %s is %t", card.Name, formatLegal))
	var cardLegalForFormat CardLegalForFormat
	cardLegalForFormat.ID = string(cardId)
	if card.Name != "" {
		cardLegalForFormat.Name = card.Name
	} else {
		helpers.Debug(fmt.Sprintf("Card Name empty for %s", cardId))
		cardLegalForFormat.Name = fmt.Sprintf("UNKNOWN NAME: %s", cardId)
	}
	cardLegalForFormat.Legal = formatLegal
	return cardLegalForFormat
}

func GetCard(cardId string) Card {
	// Get card data from scryfall api
	cardUri := fmt.Sprintf("https://api.scryfall.com/cards/%s?format=json&pretty=true", cardId)
	resp, err := http.Get(cardUri)
	helpers.Check(err)
	if err != nil {
		helpers.Debug("No response from request")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	helpers.Check(err)
	helpers.Debug(string(body))

	var result Card
	if err := json.Unmarshal(body, &result); err != nil {
		helpers.Debug("Can not unmarshal JSON")
	}
	helpers.Check(err)
	return result
}

func (c *Card) IsCardLegalForFormat(format string) bool {
	// Check card fields to see if one matches the Struct fields.
	// If so, get value and compare with format.
	// Return True if Legal in game format. Otherwise return false
	field := reflect.ValueOf(c.Legalities).FieldByName(format)
	if field.IsValid() {
		value := field.Interface().(string)
		helpers.Debug(fmt.Sprintf("Card Format found: %s", value))
		if value == "legal" {
			return true
		} else {
			return false
		}
	} else {
		helpers.Debug(fmt.Sprintf("Invalid field: %s looking up format for %s", field, c.ID))
		return false
	}
}
