package scryfall

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"regexp"

	"github.com/Pacobart/mtg-deck-legality-checker/internal/helpers"
)

type CardEntry struct {
	Object            string     `json:"object"`
	ID                string     `json:"id"`
	DeckID            string     `json:"deck_id"`
	Section           string     `json:"section"`
	Cardinality       float64    `json:"cardinality"`
	Count             int        `json:"count"`
	RawText           string     `json:"raw_text"`
	Found             bool       `json:"found"`
	PrintingSpecified bool       `json:"printing_specified"`
	Finish            bool       `json:"finish"`
	CardDigest        CardDigest `json:"card_digest"`
}

type CardDigest struct {
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
}

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
		Maybeboard []CardEntry `json:"maybeboard"`
		Nonlands   []CardEntry `json:"nonlands"`
		Outside    []CardEntry `json:"outside"`
		Commanders []CardEntry `json:"commanders"`
		Lands      []CardEntry `json:"lands"`
	} `json:"entries"`
}

type Card struct {
	Object        string `json:"object"`
	ID            string `json:"id"`
	OracleID      string `json:"oracle_id"`
	MultiverseIds []int  `json:"multiverse_ids"`
	MtgoID        int    `json:"mtgo_id"`
	ArenaID       int    `json:"arena_id"`
	TcgplayerID   int    `json:"tcgplayer_id"`
	CardmarketID  int    `json:"cardmarket_id"`
	Name          string `json:"name"`
	Lang          string `json:"lang"`
	ReleasedAt    string `json:"released_at"`
	URI           string `json:"uri"`
	ScryfallURI   string `json:"scryfall_uri"`
	Layout        string `json:"layout"`
	HighresImage  bool   `json:"highres_image"`
	ImageStatus   string `json:"image_status"`
	ImageUris     struct {
		Small      string `json:"small"`
		Normal     string `json:"normal"`
		Large      string `json:"large"`
		Png        string `json:"png"`
		ArtCrop    string `json:"art_crop"`
		BorderCrop string `json:"border_crop"`
	} `json:"image_uris"`
	ManaCost      string   `json:"mana_cost"`
	Cmc           float64  `json:"cmc"`
	TypeLine      string   `json:"type_line"`
	OracleText    string   `json:"oracle_text"`
	Power         string   `json:"power"`
	Toughness     string   `json:"toughness"`
	Colors        []string `json:"colors"`
	ColorIdentity []string `json:"color_identity"`
	Keywords      []any    `json:"keywords"`
	AllParts      []struct {
		Object    string `json:"object"`
		ID        string `json:"id"`
		Component string `json:"component"`
		Name      string `json:"name"`
		TypeLine  string `json:"type_line"`
		URI       string `json:"uri"`
	} `json:"all_parts"`
	Legalities struct {
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
	Games           []string `json:"games"`
	Reserved        bool     `json:"reserved"`
	Foil            bool     `json:"foil"`
	Nonfoil         bool     `json:"nonfoil"`
	Finishes        []string `json:"finishes"`
	Oversized       bool     `json:"oversized"`
	Promo           bool     `json:"promo"`
	Reprint         bool     `json:"reprint"`
	Variation       bool     `json:"variation"`
	SetID           string   `json:"set_id"`
	Set             string   `json:"set"`
	SetName         string   `json:"set_name"`
	SetType         string   `json:"set_type"`
	SetURI          string   `json:"set_uri"`
	SetSearchURI    string   `json:"set_search_uri"`
	ScryfallSetURI  string   `json:"scryfall_set_uri"`
	RulingsURI      string   `json:"rulings_uri"`
	PrintsSearchURI string   `json:"prints_search_uri"`
	CollectorNumber string   `json:"collector_number"`
	Digital         bool     `json:"digital"`
	Rarity          string   `json:"rarity"`
	FlavorText      string   `json:"flavor_text"`
	CardBackID      string   `json:"card_back_id"`
	Artist          string   `json:"artist"`
	ArtistIds       []string `json:"artist_ids"`
	IllustrationID  string   `json:"illustration_id"`
	BorderColor     string   `json:"border_color"`
	Frame           string   `json:"frame"`
	FrameEffects    []string `json:"frame_effects"`
	SecurityStamp   string   `json:"security_stamp"`
	FullArt         bool     `json:"full_art"`
	Textless        bool     `json:"textless"`
	Booster         bool     `json:"booster"`
	StorySpotlight  bool     `json:"story_spotlight"`
	EdhrecRank      int      `json:"edhrec_rank"`
	PennyRank       int      `json:"penny_rank"`
	Prices          struct {
		Usd       string `json:"usd"`
		UsdFoil   string `json:"usd_foil"`
		UsdEtched any    `json:"usd_etched"`
		Eur       string `json:"eur"`
		EurFoil   string `json:"eur_foil"`
		Tix       string `json:"tix"`
	} `json:"prices"`
	RelatedUris struct {
		Gatherer                  string `json:"gatherer"`
		TcgplayerInfiniteArticles string `json:"tcgplayer_infinite_articles"`
		TcgplayerInfiniteDecks    string `json:"tcgplayer_infinite_decks"`
		Edhrec                    string `json:"edhrec"`
	} `json:"related_uris"`
	PurchaseUris struct {
		Tcgplayer   string `json:"tcgplayer"`
		Cardmarket  string `json:"cardmarket"`
		Cardhoarder string `json:"cardhoarder"`
	} `json:"purchase_uris"`
}

type CardLegalForFormat struct {
	ID    string
	Name  string
	Legal bool
}

func GetDeckList(url string) DeckList {
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
	categoryMap := map[string][]CardEntry{
		"Commanders": d.Entries.Commanders,
		"Lands":      d.Entries.Lands,
		"Maybeboard": d.Entries.Maybeboard,
		"Nonlands":   d.Entries.Nonlands,
		"Outside":    d.Entries.Outside,
	}

	var cards []CardEntry

	for _, entries := range categoryMap {
		cards = append(cards, entries...)
	}
	return cards
}

func CheckCardsForFormatLegal(format string, cards []CardEntry) []CardLegalForFormat {
	// Pass a list of Cards in to iterate over and see if card is valid for specified format
	var cardLegalities []CardLegalForFormat
	for _, cardentry := range cards {
		cardId := cardentry.CardDigest.ID
		card := GetCard(string(cardId))
		cardLegality := cardentry.CheckCardForFormatLegal(format, card)
		cardLegalities = append(cardLegalities, cardLegality)
	}
	return cardLegalities
}

func (c *CardEntry) CheckCardForFormatLegal(format string, card Card) CardLegalForFormat {
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
