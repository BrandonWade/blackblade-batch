package runner

import (
	"time"

	"github.com/BrandonWade/poseidon-batch/services"
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
	page := 1

	for {
		// Request cards from the API
		res, err := b.cardService.ListCards(page)
		if err != nil {
			b.logger.Errorf("error downloading cards from page %d: %s", page, err.Error())
			return
		}

		// Upsert the response into the database
		_, err = b.cardService.UpsertCards(res.Cards)
		if err != nil {
			b.logger.Errorf("error upserting cards for page %d: %s", page, err.Error())
			return
		}

		// Finish if there are no more cards to fetch
		if res.HasMore == false || res.NextPage == nil {
			break
		}

		page++
		time.Sleep(150 * time.Millisecond)
	}
}
