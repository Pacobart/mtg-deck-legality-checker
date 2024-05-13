package scryfall

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/Pacobart/mtg-deck-legality-checker/internal/scryfall/helpers"
)

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
	Sections     struct {
		Primary   []string `json:"primary"`
		Secondary []string `json:"secondary"`
	} `json:"sections"`
	Entries struct {
		Maybeboard []struct {
			Object            string  `json:"object"`
			ID                string  `json:"id"`
			DeckID            string  `json:"deck_id"`
			Section           string  `json:"section"`
			Cardinality       float64 `json:"cardinality"`
			Count             int     `json:"count"`
			RawText           string  `json:"raw_text"`
			Found             bool    `json:"found"`
			PrintingSpecified bool    `json:"printing_specified"`
			Finish            bool    `json:"finish"`
			CardDigest        any     `json:"card_digest"`
		} `json:"maybeboard"`
		Nonlands []struct {
			Object            string  `json:"object"`
			ID                string  `json:"id"`
			DeckID            string  `json:"deck_id"`
			Section           string  `json:"section"`
			Cardinality       float64 `json:"cardinality"`
			Count             int     `json:"count"`
			RawText           string  `json:"raw_text"`
			Found             bool    `json:"found"`
			PrintingSpecified bool    `json:"printing_specified"`
			Finish            any     `json:"finish"`
			CardDigest        struct {
				Object          string `json:"object"`
				ID              string `json:"id"`
				OracleID        string `json:"oracle_id"`
				Name            string `json:"name"`
				ScryfallURI     string `json:"scryfall_uri"`
				ManaCost        string `json:"mana_cost"`
				TypeLine        string `json:"type_line"`
				CollectorNumber string `json:"collector_number"`
				Set             string `json:"set"`
				ImageUris       struct {
					Front string `json:"front"`
				} `json:"image_uris"`
			} `json:"card_digest"`
		} `json:"nonlands"`
		Outside []struct {
			Object            string  `json:"object"`
			ID                string  `json:"id"`
			DeckID            string  `json:"deck_id"`
			Section           string  `json:"section"`
			Cardinality       float64 `json:"cardinality"`
			Count             int     `json:"count"`
			RawText           string  `json:"raw_text"`
			Found             bool    `json:"found"`
			PrintingSpecified bool    `json:"printing_specified"`
			Finish            bool    `json:"finish"`
			CardDigest        any     `json:"card_digest"`
		} `json:"outside"`
		Commanders []struct {
			Object            string  `json:"object"`
			ID                string  `json:"id"`
			DeckID            string  `json:"deck_id"`
			Section           string  `json:"section"`
			Cardinality       float64 `json:"cardinality"`
			Count             int     `json:"count"`
			RawText           string  `json:"raw_text"`
			Found             bool    `json:"found"`
			PrintingSpecified bool    `json:"printing_specified"`
			Finish            bool    `json:"finish"`
			CardDigest        any     `json:"card_digest"`
		} `json:"commanders"`
		Lands []struct {
			Object            string  `json:"object"`
			ID                string  `json:"id"`
			DeckID            string  `json:"deck_id"`
			Section           string  `json:"section"`
			Cardinality       float64 `json:"cardinality"`
			Count             int     `json:"count"`
			RawText           string  `json:"raw_text"`
			Found             bool    `json:"found"`
			PrintingSpecified bool    `json:"printing_specified"`
			Finish            any     `json:"finish"`
			CardDigest        struct {
				Object          string `json:"object"`
				ID              string `json:"id"`
				OracleID        string `json:"oracle_id"`
				Name            string `json:"name"`
				ScryfallURI     string `json:"scryfall_uri"`
				ManaCost        string `json:"mana_cost"`
				TypeLine        string `json:"type_line"`
				CollectorNumber string `json:"collector_number"`
				Set             string `json:"set"`
				ImageUris       struct {
					Front string `json:"front"`
				} `json:"image_uris"`
			} `json:"card_digest"`
		} `json:"lands"`
	} `json:"entries"`
}

func GetDeckList(url string) []Decklist {
	deckRegex, _ := regexp.Compile("[0-9A-Fa-f]{8}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{12}")
	deckId := deckRegex.FindString(url)
	fullUrl := fmt.Sprintf("https://api.scryfall.com/decks/%s/export/json", deckId)

	resp, err := http.Get(fullUrl)
	helpers.Check(err)
	if err != nil {
		fmt.Println("No response from request")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	helpers.Check(err)
	helpers.Debug(string(body))

	var result DeckList
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}
	helpers.Check(err)
	return result
}
