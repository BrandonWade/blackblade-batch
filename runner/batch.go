package runner

import (
	"fmt"

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

	// TODO: Parse cards from resBody and upsert into DB
	fmt.Printf("%#v", resBody)

	// Upsert the response into the database
	// _, err = b.cardService.UpsertCards(res.Cards)
	// if err != nil {
	// 	b.logger.Errorf("error upserting cards for page %d: %s", page, err.Error())
	// 	return
	// }
}
