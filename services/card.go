package services

import (
	scryfall "github.com/BlueMonday/go-scryfall"
	"github.com/BrandonWade/poseidon-batch/repositories"
	"github.com/sirupsen/logrus"
)

// CardService interface for working with a cardService
type CardService interface {
	GetCards(int) (scryfall.CardListResponse, error)
	UpsertCards([]scryfall.Card) (int64, error)
}

type cardService struct {
	logger   *logrus.Logger
	cardRepo repositories.CardRepository
}

// NewCardService create a new CardService instance
func NewCardService(logger *logrus.Logger, cardRepo repositories.CardRepository) CardService {
	return &cardService{
		logger,
		cardRepo,
	}
}

// GetCards returns a paginated set of cards from the Scryfall API
func (c *cardService) GetCards(page int) (scryfall.CardListResponse, error) {
	// TODO: Make request to Scryfall API
	return scryfall.CardListResponse{}, nil
}

// UpsertCards upserts cards into the database
func (c *cardService) UpsertCards(cards []scryfall.Card) (int64, error) {
	return c.cardRepo.UpsertCards(cards)
}
