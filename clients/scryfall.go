package clients

import (
	"context"
	scryfall "github.com/BlueMonday/go-scryfall"
	"github.com/sirupsen/logrus"
)

// ScryfallClient interface for working with a scryfallClient
type ScryfallClient interface {
	ListCards(int) (scryfall.CardListResponse, error)
}

type scryfallClient struct {
	logger *logrus.Logger
	client *scryfall.Client
}

// NewScryfallClient create a new ScryfallClient instance
func NewScryfallClient(logger *logrus.Logger, client *scryfall.Client) ScryfallClient {
	return &scryfallClient{
		logger,
		client,
	}
}

func (s *scryfallClient) ListCards(page int) (scryfall.CardListResponse, error) {
	return s.client.ListCards(context.Background(), scryfall.ListCardsOptions{
		Page: page,
	})
}
