package runner

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/BrandonWade/blackblade-batch/models"
	"github.com/BrandonWade/blackblade-batch/services"
	"github.com/sirupsen/logrus"
)

// BatchRunner interface for working with a batchRunner
type BatchRunner interface {
	Run()
}

type batchRunner struct {
	logger      *logrus.Logger
	cardService services.CardService
}

// NewBatchRunner create a new BatchRunner instance
func NewBatchRunner(logger *logrus.Logger, cardService services.CardService) BatchRunner {
	return &batchRunner{
		logger,
		cardService,
	}
}

// Run download data from the Scryfall API and process it
func (b *batchRunner) Run() {
	b.logger.Println("Batch starting...")
	start := time.Now()

	err := b.processCards()
	if err != nil {
		return
	}

	err = b.processRulings()
	if err != nil {
		return
	}

	elapsed := time.Since(start)
	b.logger.Printf("Batch completed in %s.", elapsed)
}

func (b *batchRunner) processCards() error {
	b.logger.Println("Downloading default-cards bulk data file...")
	defaultCards, err := b.cardService.GetDefaultCards()
	if err != nil {
		b.logger.Errorf("error fetching default cards from api: %s", err.Error())
		return err
	}

	if (defaultCards == models.ScryfallBulkData{}) {
		err = errors.New("default cards not found")
		b.logger.Errorf(err.Error())
		return err
	}

	filepath := fmt.Sprintf("defaultcards-%v.json", int32(time.Now().Unix()))
	err = b.cardService.DownloadDefaultCardData(defaultCards.DownloadURI, filepath)
	if err != nil {
		b.logger.Fatalf("error downloading default cards data from api: %s", err.Error())
		return err
	}

	b.logger.Println("Processing default-cards bulk data file...")

	file, err := os.Open(filepath)
	if err != nil {
		b.logger.Fatalf("error opening default cards data file: %s", err.Error())
		return err
	}

	dec := json.NewDecoder(file)
	// dec.DisallowUnknownFields()

	// read opening bracket
	_, err = dec.Token()
	if err != nil {
		b.logger.Errorf("error parsing card data: %s", err.Error())
		return err
	}

	cards := []models.ScryfallCard{}
	for dec.More() {
		var card models.ScryfallCard
		err = dec.Decode(&card)
		if err != nil {
			b.logger.Errorf("error decoding card: %s", err.Error())

			if err.Error() == "not at beginning of value" {
				b.logger.Fatalf("bulk data file contents in unexpected format - is the scryfall bulk data api broken?")
			}
		}

		validPrint := card.Lang == "en" && !card.Digital
		validCardType := card.TypeLine != "Vanguard" && card.Layout != "art_series" && card.Layout != "planar" && card.Layout != "scheme"
		validSetType := card.SetType != "memorabilia"
		validFunny := card.SetType != "funny" || (card.SetType == "funny" && (strings.Contains(card.TypeLine, "Plains") ||
			strings.Contains(card.TypeLine, "Island") ||
			strings.Contains(card.TypeLine, "Swamp") ||
			strings.Contains(card.TypeLine, "Mountain") ||
			strings.Contains(card.TypeLine, "Forest")))

		if validPrint && validCardType && validSetType && validFunny {
			cards = append(cards, card)
		}

		if len(cards) == 100 {
			err = b.cardService.UpsertCards(cards)
			if err != nil {
				b.logger.Errorf("error upserting cards: %s", err.Error())
			}

			cards = []models.ScryfallCard{}
		}
	}

	if len(cards) > 0 {
		err = b.cardService.UpsertCards(cards)
		if err != nil {
			b.logger.Errorf("error upserting cards: %s", err.Error())
			return err
		}
	}

	// read closing bracket
	_, err = dec.Token()
	if err != nil {
		b.logger.Errorf("error parsing card data: %s", err.Error())
		return err
	}

	b.logger.Println("Calculating cards.faces_json column values...")
	err = b.cardService.GenerateCardFacesJSON()
	if err != nil {
		b.logger.Errorf("error generating cards.faces_json values: %s", err.Error())
		return err
	}

	b.logger.Println("Calculating card_sets_list table...")
	err = b.cardService.GenerateCardSetsJSON()
	if err != nil {
		b.logger.Errorf("error generating card_sets_list table: %s", err.Error())
		return err
	}

	b.logger.Println("Calculating sets table...")
	err = b.cardService.GenerateSets()
	if err != nil {
		b.logger.Errorf("error generating sets table: %s", err.Error())
		return err
	}

	return nil
}

// When go gets generics, it might be possible to de-dupe a lot of this code...
func (b *batchRunner) processRulings() error {
	b.logger.Println("Downloading rulings bulk data file...")
	rulingsData, err := b.cardService.GetRulings()
	if err != nil {
		b.logger.Errorf("error fetching card rulings from api: %s", err.Error())
		return err
	}

	if (rulingsData == models.ScryfallBulkData{}) {
		err = errors.New("rulings not found")
		b.logger.Errorf(err.Error())
		return err
	}

	filepath := fmt.Sprintf("rulings-%v.json", int32(time.Now().Unix()))
	err = b.cardService.DownloadRulingsData(rulingsData.DownloadURI, filepath)
	if err != nil {
		b.logger.Fatalf("error downloading card rulings data from api: %s", err.Error())
		return err
	}

	b.logger.Println("Processing rulings bulk data file...")

	file, err := os.Open(filepath)
	if err != nil {
		b.logger.Fatalf("error opening rulings data file: %s", err.Error())
		return err
	}

	dec := json.NewDecoder(file)
	// dec.DisallowUnknownFields()

	// read opening bracket
	_, err = dec.Token()
	if err != nil {
		b.logger.Errorf("error parsing ruling data: %s", err.Error())
		return err
	}

	rulings := []models.ScryfallRuling{}
	for dec.More() {
		var ruling models.ScryfallRuling
		err = dec.Decode(&ruling)
		if err != nil {
			b.logger.Errorf("error decoding ruling: %s", err.Error())

			if err.Error() == "not at beginning of value" {
				b.logger.Fatalf("bulk data file contents in unexpected format - is the scryfall bulk data api broken?")
			}
		}

		rulings = append(rulings, ruling)

		if len(rulings) == 100 {
			err = b.cardService.InsertRulings(rulings)
			if err != nil {
				b.logger.Errorf("error inserting rulings: %s", err.Error())
			}

			rulings = []models.ScryfallRuling{}
		}
	}

	if len(rulings) > 0 {
		err = b.cardService.InsertRulings(rulings)
		if err != nil {
			b.logger.Errorf("error inserting rulings: %s", err.Error())
			return err
		}
	}

	// read closing bracket
	_, err = dec.Token()
	if err != nil {
		b.logger.Errorf("error parsing ruling data: %s", err.Error())
		return err
	}

	b.logger.Println("Calculating card_rulings_list table...")
	err = b.cardService.GenerateRulingsJSON()
	if err != nil {
		b.logger.Errorf("error generating card_rulings_list table: %s", err.Error())
		return err
	}

	return nil
}
