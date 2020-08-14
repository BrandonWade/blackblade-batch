package runner

import (
	"encoding/json"

	scryfall "github.com/BlueMonday/go-scryfall"
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

	resBody, err := b.cardService.DownloadAllCardData(allCards.URI)
	if err != nil {
		b.logger.Errorf("error downloading all cards data from api: %s", err.Error())
		return
	}

	dec := json.NewDecoder(resBody)

	// read open bracket
	_, err = dec.Token()
	if err != nil {
		b.logger.Errorf("error parsing card data: %s", err.Error())
		return
	}

	cards := []scryfall.Card{}
	for dec.More() {
		var card scryfall.Card
		err = dec.Decode(&card)
		if err != nil {
			b.logger.Errorf("error decoding card: %s", err.Error())
		}

		cards = append(cards, card)
		if len(cards) == 100 {
			_, err = b.cardService.UpsertCards(cards)
			if err != nil {
				b.logger.Errorf("error upserting cards: %s", err.Error())
				return
			}

			cards = []scryfall.Card{}
		}
	}

	if len(cards) > 0 {
		_, err = b.cardService.UpsertCards(cards)
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
}
