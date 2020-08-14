package models

// BulkData represents a BulkData object from Scryfall.
type BulkData struct {
	Object          string `json:"object"`
	ID              string `json:"id"`
	Type            string `json:"type"`
	UpdatedAt       string `json:"updated_at"`
	URI             string `json:"uri"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	CompressedSize  int64  `json:"compressed_size"`
	DownloadURI     string `json:"download_uri"`
	ContentType     string `json:"content_type"`
	ContentEncoding string `json:"content_encoding"`
}
