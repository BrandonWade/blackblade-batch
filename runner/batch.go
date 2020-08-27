package runner

import (
	"encoding/json"
	"fmt"
	"os"
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

// Run download cards from the Scryfall API and upsert them into the database
func (b *batchRunner) Run() {
	b.logger.Println("Batch starting...")
	start := time.Now()

	allCards, err := b.cardService.GetAllCards()
	if err != nil {
		b.logger.Errorf("error fetching all cards from api: %s", err.Error())
		return
	}

	if (allCards == models.BulkData{}) {
		b.logger.Errorf("all cards not found")
		return
	}

	// TODO: Compare allCards.UpdatedAt against last run

	filepath := fmt.Sprintf("allcards-%v.json", int32(time.Now().Unix()))
	err = b.cardService.DownloadAllCardData(allCards.DownloadURI, filepath)
	if err != nil {
		b.logger.Fatalf("error downloading all cards data from api: %s", err.Error())
		return
	}

	b.logger.Println("Processing all-cards bulk data file...")

	file, err := os.Open(filepath)
	if err != nil {
		b.logger.Fatalf("error opening all cards data file: %s", err.Error())
		return
	}

	dec := json.NewDecoder(file)
	dec.DisallowUnknownFields()

	// read open bracket
	_, err = dec.Token()
	if err != nil {
		b.logger.Errorf("error parsing card data: %s", err.Error())
		return
	}

	cards := []models.ScryfallCard{}
	for dec.More() {
		var card models.ScryfallCard
		err = dec.Decode(&card)
		if err != nil {
			b.logger.Errorf("error decoding card: %s", err.Error())
		}

		if card.Lang == "en" && card.TypeLine != "Vanguard" {
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
			return
		}
	}

	// read closing bracket
	_, err = dec.Token()
	if err != nil {
		b.logger.Errorf("error parsing card data: %s", err.Error())
		return
	}

	// TODO: Populate card_sets_images table

	elapsed := time.Since(start)
	b.logger.Printf("Batch completed in %s.", elapsed)
}
