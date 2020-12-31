package services

import (
	"github.com/BrandonWade/blackblade-batch/clients"
	"github.com/BrandonWade/blackblade-batch/models"
	"github.com/BrandonWade/blackblade-batch/repositories"
	"github.com/sirupsen/logrus"
)

// CardService interface for working with a cardService
type CardService interface {
	GetDefaultCards() (models.ScryfallBulkData, error)
	DownloadDefaultCardData(uri, filepath string) error
	GetRulings() (models.ScryfallBulkData, error)
	DownloadRulingsData(uri, filepath string) error
	UpsertCards(cards []models.ScryfallCard) error
	GenerateCardFacesJSON() error
	GenerateCardSetsJSON() error
	GenerateSets() error
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

// GetDefaultCards returns the default_cards bulk data from the Scryfall API.
func (c *cardService) GetDefaultCards() (models.ScryfallBulkData, error) {
	return c.scryfallClient.GetBulkData("default-cards")
}

// DownloadDefaultCardData downloads the default_cards bulk data file from the scryfall API.
func (c *cardService) DownloadDefaultCardData(uri, filepath string) error {
	return c.scryfallClient.DownloadBulkData("default-cards", uri, filepath)
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

// GenerateCardFacesJSON calculates the set name and images for each card in the database and saves the result.
func (c *cardService) GenerateCardFacesJSON() error {
	return c.cardRepo.GenerateCardFacesJSON()
}

// GenerateCardSetsJSON aggregates the faces JSON for distinct card in the database and saves the result.
func (c *cardService) GenerateCardSetsJSON() error {
	return c.cardRepo.GenerateCardSetsJSON()
}

// GenerateSets calculates a list of unique set names and codes in the database and saves the result.
func (c *cardService) GenerateSets() error {
	return c.cardRepo.GenerateSets()
}

// GenerateRulingsJSON aggregates the rulings for each distinct card in the database and saves the result.
func (c *cardService) GenerateRulingsJSON() error {
	return c.cardRepo.GenerateRulingsJSON()
}

// InsertRulings inserts the provided rulings into the database.
func (c *cardService) InsertRulings(rulings []models.ScryfallRuling) error {
	return c.cardRepo.InsertRulings(rulings)
}
