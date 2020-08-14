package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	scryfall "github.com/BlueMonday/go-scryfall"
	"github.com/BrandonWade/blackblade-batch/models"
	"github.com/sirupsen/logrus"
)

// ScryfallClient interface for working with a scryfallClient.
type ScryfallClient interface {
	GetAllCards() (models.BulkData, error)
	DownloadAllCardData(url string) (io.ReadCloser, error)
}

type scryfallClient struct {
	baseURL string
	logger  *logrus.Logger
	client  *scryfall.Client
}

// NewScryfallClient create a new ScryfallClient instance.
func NewScryfallClient(logger *logrus.Logger, client *scryfall.Client) ScryfallClient {
	return &scryfallClient{
		"https://api.scryfall.com",
		logger,
		client,
	}
}

// GetAllCards returns the all_cards bulk data from the Scryfall API.
func (s *scryfallClient) GetAllCards() (models.BulkData, error) {
	url := fmt.Sprintf("%s/bulk-data/all-cards", s.baseURL)

	res, err := http.Get(url)
	if err != nil {
		return models.BulkData{}, err
	}
	defer res.Body.Close()

	allCards := models.BulkData{}
	err = json.NewDecoder(res.Body).Decode(&allCards)
	if err != nil {
		return models.BulkData{}, err
	}

	return allCards, nil
}

// DownloadAllCardData downloads the contents of the all_cards bulk data file from the Scryfall API.
func (s *scryfallClient) DownloadAllCardData(url string) (io.ReadCloser, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}
