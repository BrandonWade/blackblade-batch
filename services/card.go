package services

import (
	"github.com/BrandonWade/blackblade-batch/clients"
	"github.com/BrandonWade/blackblade-batch/models"
	"github.com/BrandonWade/blackblade-batch/repositories"
	"github.com/sirupsen/logrus"
)

// CardService interface for working with a cardService
type CardService interface {
	GetAllCards() (models.ScryfallBulkData, error)
	DownloadAllCardData(uri, filepath string) error
	GetRulings() (models.ScryfallBulkData, error)
	DownloadRulingsData(uri, filepath string) error
	UpsertCards(cards []models.ScryfallCard) error
	GenerateFacesJSON() error
	GenerateSetsJSON() error
	InsertRulings(rulings []models.ScryfallRuling) error
	GenerateRulingsJSON() error
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
func (c *cardService) GetAllCards() (models.ScryfallBulkData, error) {
	return c.scryfallClient.GetBulkData("all-cards")
}

// DownloadAllCardData downloads the all_cards bulk data file from the scryfall API.
func (c *cardService) DownloadAllCardData(uri, filepath string) error {
	return c.scryfallClient.DownloadBulkData("all-cards", uri, filepath)
}

// GetRulings returns the rulings bulk data from the Scryfall API.
func (c *cardService) GetRulings() (models.ScryfallBulkData, error) {
	return c.scryfallClient.GetBulkData("rulings")
}

// DownloadRulingsData downloads the rulings bulk data file from the scryfall API.
func (c *cardService) DownloadRulingsData(uri, filepath string) error {
	return c.scryfallClient.DownloadBulkData("rulings", uri, filepath)
}

// UpsertCards upserts the provided cards into the database.
func (c *cardService) UpsertCards(cards []models.ScryfallCard) error {
	return c.cardRepo.UpsertCards(cards)
}

// GenerateFacesJSON calculates the set name and images for each card in the database and saves the result.
func (c *cardService) GenerateFacesJSON() error {
	return c.cardRepo.GenerateFacesJSON()
}

// GenerateSetsJSON aggregates the faces JSON for distinct card in the database and saves the result.
func (c *cardService) GenerateSetsJSON() error {
	return c.cardRepo.GenerateSetsJSON()
}

// GenerateRulingsJSON aggregates the rulings for each distinct card in the database and saves the result.
func (c *cardService) GenerateRulingsJSON() error {
	return c.cardRepo.GenerateRulingsJSON()
}

// InsertRulings inserts the provided rulings into the database.
func (c *cardService) InsertRulings(rulings []models.ScryfallRuling) error {
	return c.cardRepo.InsertRulings(rulings)
}
