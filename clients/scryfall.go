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
	GetBulkData(dataType string) (models.ScryfallBulkData, error)
	DownloadBulkData(dataType, uri, filepath string) error
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

// GetBulkData returns the bulk data of the specified type from the Scryfall API.
func (s *scryfallClient) GetBulkData(dataType string) (models.ScryfallBulkData, error) {
	url := fmt.Sprintf("%s/bulk-data/%s", s.baseURL, dataType)

	res, err := http.Get(url)
	if err != nil {
		return models.ScryfallBulkData{}, err
	}
	defer res.Body.Close()

	data := models.ScryfallBulkData{}
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return models.ScryfallBulkData{}, err
	}

	return data, nil
}

// DownloadBulkData downloads the contents of the specified bulk data file from the Scryfall API.
func (s *scryfallClient) DownloadBulkData(dataType, uri, filepath string) error {
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

	s.logger.Printf("Successfully downloaded %s bulk data file %s", dataType, filepath)

	return nil
}
