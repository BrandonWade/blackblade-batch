package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	scryfall "github.com/BlueMonday/go-scryfall"
	"github.com/BrandonWade/blackblade-batch/models"
	"github.com/sirupsen/logrus"
)

// ScryfallClient interface for working with a scryfallClient.
type ScryfallClient interface {
	GetAllCards() (models.BulkData, error)
	DownloadAllCardData(uri, filepath string) error
}

type scryfallClient struct {
	baseURL string
	logger  *logrus.Logger
	client  *scryfall.Client
}

// NewScryfallClient create a new ScryfallClient instance.
func NewScryfallClient(baseURL string, logger *logrus.Logger, client *scryfall.Client) ScryfallClient {
	return &scryfallClient{
		baseURL,
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
func (s *scryfallClient) DownloadAllCardData(uri, filepath string) error {
	res, err := http.Get(uri)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, res.Body)
	if err != nil {
		return err
	}

	s.logger.Println("Successfully downloaded all-cards bulk data file")

	return nil
}
