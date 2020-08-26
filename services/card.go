package services

import (
	"github.com/BrandonWade/blackblade-batch/clients"
	"github.com/BrandonWade/blackblade-batch/models"
	"github.com/BrandonWade/blackblade-batch/repositories"
	"github.com/sirupsen/logrus"
)

// CardService interface for working with a cardService
type CardService interface {
	GetAllCards() (models.BulkData, error)
	DownloadAllCardData(uri, filepath string) error
	UpsertCards(cards []models.ScryfallCard) error
}

type cardService struct {
	logger         *logrus.Logger
	scryfallClient clients.ScryfallClient
	cardRepo       repositories.CardRepository
}

// NewCardService create a new CardService instance
func NewCardService(logger *logrus.Logger, scryfallClient clients.ScryfallClient, cardRepo repositories.CardRepository) CardService {
	return &cardService{
		logger,
		scryfallClient,
		cardRepo,
	}
}

// GetAllCards returns the all_cards bulk data from the Scryfall API.
func (c *cardService) GetAllCards() (models.BulkData, error) {
	return c.scryfallClient.GetAllCards()
}

func (c *cardService) DownloadAllCardData(uri, filepath string) error {
	return c.scryfallClient.DownloadAllCardData(uri, filepath)
}

// UpsertCards upserts cards into the database.
func (c *cardService) UpsertCards(cards []models.ScryfallCard) error {
	return c.cardRepo.UpsertCards(cards)
}
