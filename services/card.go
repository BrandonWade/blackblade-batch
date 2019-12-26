package services

import (
	scryfall "github.com/BlueMonday/go-scryfall"
	"github.com/BrandonWade/poseidon-batch/clients"
	"github.com/BrandonWade/poseidon-batch/repositories"
	"github.com/sirupsen/logrus"
)

// CardService interface for working with a cardService
type CardService interface {
	ListCards(int) (scryfall.CardListResponse, error)
	UpsertCards([]scryfall.Card) (int64, error)
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

// ListCards returns a paginated set of cards from the Scryfall API
func (c *cardService) ListCards(page int) (scryfall.CardListResponse, error) {
	return c.scryfallClient.ListCards(page)
}

// UpsertCards upserts cards into the database
func (c *cardService) UpsertCards(cards []scryfall.Card) (int64, error) {
	return c.cardRepo.UpsertCards(cards)
}
