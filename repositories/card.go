package repositories

import (
	"database/sql"
	"strings"

	"github.com/BrandonWade/blackblade-batch/models"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// CardRepository interface for working with a cardRepository
type CardRepository interface {
	UpsertCards(cards []models.ScryfallCard) error
	GenerateFacesJSON() error
	GenerateSetsJSON() error
	InsertRulings(rulings []models.ScryfallRuling) error
}

type cardRepository struct {
	logger *logrus.Logger
	db     *sqlx.DB
}

// NewCardRepository create a new CardRepository instance
func NewCardRepository(logger *logrus.Logger, db *sqlx.DB) CardRepository {
	return &cardRepository{
		logger,
		db,
	}
}

// UpsertCards upserts cards into the database
func (c *cardRepository) UpsertCards(cards []models.ScryfallCard) error {
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}

	for _, card := range cards {
		c.setLayout(&card)
		cardID, err := c.upsertCard(tx, card)
		if err != nil {
			return err
		}

		err = insertCardMultiverseIDs(tx, cardID, card.MultiverseIDs)
		if err != nil {
			return err
		}

		err = insertCardFrameEffects(tx, cardID, card.FrameEffects)
		if err != nil {
			return err
		}

		err = upsertCardPrices(tx, cardID, card.Prices)
		if err != nil {
			return err
		}

		cardFaces := c.getCardFaces(card)
		for i, cardFace := range cardFaces {
			_, err = c.upsertCardFace(tx, cardID, i, cardFace)
			if err != nil {
				return err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (c *cardRepository) setLayout(card *models.ScryfallCard) {
	for _, keyword := range card.Keywords {
		if strings.ToLower(keyword) == "aftermath" {
			(*card).Layout = "aftermath"
		}
	}
}

func (c *cardRepository) upsertCard(tx *sql.Tx, card models.ScryfallCard) (int64, error) {
	isWhite := contains(card.Colors, "W")
	isBlue := contains(card.Colors, "U")
	isBlack := contains(card.Colors, "B")
	isRed := contains(card.Colors, "R")
	isGreen := contains(card.Colors, "G")

	result, err := tx.Exec(`INSERT INTO cards (
		scryfall_id,
		oracle_id,
		tcgplayer_id,
		card_back_id,
		cmc,
		is_white,
		is_blue,
		is_black,
		is_red,
		is_green,
		set_code,
		set_name,
		rarity,
		layout,
		border_color,
		frame,
		released_at,
		has_foil,
		has_nonfoil,
		is_oversized,
		is_reserved,
		is_booster,
		is_digital_only,
		is_full_art,
		is_textless,
		is_reprint,
		has_highres_image,
		rulings_uri,
		scryfall_uri
	) VALUES (
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?
	) ON DUPLICATE KEY UPDATE
		scryfall_id = ?,
		oracle_id = ?,
		tcgplayer_id = ?,
		card_back_id = ?,
		cmc = ?,
		is_white = ?,
		is_blue = ?,
		is_black = ?,
		is_red = ?,
		is_green = ?,
		set_code = ?,
		set_name = ?,
		rarity = ?,
		layout = ?,
		border_color = ?,
		frame = ?,
		released_at = ?,
		has_foil = ?,
		has_nonfoil = ?,
		is_oversized = ?,
		is_reserved = ?,
		is_booster = ?,
		is_digital_only = ?,
		is_full_art = ?,
		is_textless = ?,
		is_reprint = ?,
		has_highres_image = ?,
		rulings_uri = ?,
		scryfall_uri = ?
	`,
		card.ID,
		card.OracleID,
		card.TCGPlayerID,
		card.CardBackID,
		card.CMC,
		isWhite,
		isBlue,
		isBlack,
		isRed,
		isGreen,
		card.Set,
		card.SetName,
		card.Rarity,
		card.Layout,
		card.BorderColor,
		card.Frame,
		card.ReleasedAt,
		card.Foil,
		card.Nonfoil,
		card.Oversized,
		card.Reserved,
		card.Booster,
		card.Digital,
		card.FullArt,
		card.Textless,
		card.Reprint,
		card.HighresImage,
		card.RulingsURI,
		card.ScryfallURI,
		card.ID,
		card.OracleID,
		card.TCGPlayerID,
		card.CardBackID,
		card.CMC,
		isWhite,
		isBlue,
		isBlack,
		isRed,
		isGreen,
		card.Set,
		card.SetName,
		card.Rarity,
		card.Layout,
		card.BorderColor,
		card.Frame,
		card.ReleasedAt,
		card.Foil,
		card.Nonfoil,
		card.Oversized,
		card.Reserved,
		card.Booster,
		card.Digital,
		card.FullArt,
		card.Textless,
		card.Reprint,
		card.HighresImage,
		card.RulingsURI,
		card.ScryfallURI,
	)
	if err != nil {
		return 0, err
	}

	cardID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// When using ON DUPLICATE KEY UPDATE, MySQL does allow you to return the ID for the existing row
	// using LAST_INSERT_ID(id), however this causes the AUTO_INCREMENT value to be increased by 2 per
	// row every time the batch runs as every row is always updated.
	if cardID == 0 {
		err := c.db.Get(&cardID, `SELECT
			id
			FROM cards c
			WHERE c.scryfall_id = ?
		`,
			card.ID,
		)
		if err != nil {
			return 0, err
		}
	}

	return cardID, nil
}

func insertCardMultiverseIDs(tx *sql.Tx, cardID int64, multiverseIDs []int) error {
	for _, multiverseID := range multiverseIDs {
		_, err := tx.Exec(`INSERT IGNORE INTO card_multiverse_ids (
			card_id,
			multiverse_id
		) VALUES (
			?,
			?
		)
		`,
			cardID,
			multiverseID,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func insertCardFrameEffects(tx *sql.Tx, cardID int64, frameEffects []string) error {
	for _, frameEffect := range frameEffects {
		_, err := tx.Exec(`INSERT IGNORE INTO card_frame_effects (
			card_id,
			frame_effect
		) VALUES (
			?,
			?
		)
		`,
			cardID,
			frameEffect,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func upsertCardPrices(tx *sql.Tx, cardID int64, prices models.ScryfallPrices) error {
	_, err := tx.Exec(`INSERT INTO card_prices (
		card_id,
		usd,
		usd_foil,
		eur,
		tix
	) VALUES (
		?,
		?,
		?,
		?,
		?
	) ON DUPLICATE KEY UPDATE
		usd = ?,
		usd_foil = ?,
		eur = ?,
		tix = ?
	`,
		cardID,
		prices.USD,
		prices.USDFoil,
		prices.EUR,
		prices.Tix,
		prices.USD,
		prices.USDFoil,
		prices.EUR,
		prices.Tix,
	)

	return err
}

func (c *cardRepository) getCardFaces(card models.ScryfallCard) []models.ScryfallCardFace {
	if len(card.CardFaces) > 0 {
		// Some card layouts have 2 faces but only a single set of image URIs
		if card.Layout == "flip" || card.Layout == "split" || card.Layout == "adventure" || card.Layout == "aftermath" {
			for i := range card.CardFaces {
				card.CardFaces[i].ImageURIs.Small = card.ImageURIs.Small
				card.CardFaces[i].ImageURIs.Normal = card.ImageURIs.Normal
				card.CardFaces[i].ImageURIs.Large = card.ImageURIs.Large
				card.CardFaces[i].ImageURIs.PNG = card.ImageURIs.PNG
				card.CardFaces[i].ImageURIs.ArtCrop = card.ImageURIs.ArtCrop
				card.CardFaces[i].ImageURIs.BorderCrop = card.ImageURIs.BorderCrop
			}
		}

		// Determine face derived types
		for i := range card.CardFaces {
			card.CardFaces[i].DerivedType = c.getDerivedType(card.CardFaces[i].TypeLine)
		}

		return card.CardFaces
	}

	cardFace := models.ScryfallCardFace{
		Artist:          card.Artist,
		ColorIndicator:  card.ColorIndicator,
		Colors:          card.Colors,
		FlavorText:      card.FlavorText,
		IllustrationID:  card.IllustrationID,
		ImageURIs:       card.ImageURIs,
		Loyalty:         card.Loyalty,
		ManaCost:        card.ManaCost,
		Name:            card.Name,
		OracleText:      card.OracleText,
		Power:           card.Power,
		PrintedName:     card.PrintedName,
		PrintedText:     card.PrintedText,
		PrintedTypeLine: card.PrintedTypeLine,
		Toughness:       card.Toughness,
		DerivedType:     c.getDerivedType(card.TypeLine),
		TypeLine:        card.TypeLine,
		Watermark:       card.Watermark,
	}

	return []models.ScryfallCardFace{
		cardFace,
	}
}

func (c *cardRepository) getDerivedType(cardType string) string {
	t := strings.ToLower(cardType)

	if strings.Contains(t, "creature") {
		return "creature"
	} else if strings.Contains(t, "artifact") {
		return "artifact"
	} else if strings.Contains(t, "enchantment") {
		return "enchantment"
	} else if strings.Contains(t, "instant") {
		return "instant"
	} else if strings.Contains(t, "sorcery") {
		return "sorcery"
	} else if strings.Contains(t, "planeswalker") {
		return "planeswalker"
	} else if strings.Contains(t, "land") {
		return "land"
	}

	return ""
}

func (c *cardRepository) upsertCardFace(tx *sql.Tx, cardID int64, index int, cardFace models.ScryfallCardFace) (int64, error) {
	result, err := tx.Exec(`INSERT INTO card_faces (
		card_id,
		face_index,
		artist,
		flavor_text,
		illustration_id,
		image_small,
		image_normal,
		image_large,
		image_png,
		image_art_crop,
		image_border_crop,
		mana_cost,
		name,
		oracle_text,
		power,
		toughness,
		loyalty,
		type_line,
		derived_type,
		watermark
	) VALUES (
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?
	) ON DUPLICATE KEY UPDATE
		artist = ?,
		flavor_text = ?,
		illustration_id = ?,
		image_small = ?,
		image_normal = ?,
		image_large = ?,
		image_png = ?,
		image_art_crop = ?,
		image_border_crop = ?,
		mana_cost = ?,
		name = ?,
		oracle_text = ?,
		power = ?,
		toughness = ?,
		loyalty = ?,
		type_line = ?,
		derived_type = ?,
		watermark = ?
	`,
		cardID,
		index,
		cardFace.Artist,
		cardFace.FlavorText,
		cardFace.IllustrationID,
		cardFace.ImageURIs.Small,
		cardFace.ImageURIs.Normal,
		cardFace.ImageURIs.Large,
		cardFace.ImageURIs.PNG,
		cardFace.ImageURIs.ArtCrop,
		cardFace.ImageURIs.BorderCrop,
		cardFace.ManaCost,
		cardFace.Name,
		cardFace.OracleText,
		cardFace.Power,
		cardFace.Toughness,
		cardFace.Loyalty,
		cardFace.TypeLine,
		cardFace.DerivedType,
		cardFace.Watermark,
		cardFace.Artist,
		cardFace.FlavorText,
		cardFace.IllustrationID,
		cardFace.ImageURIs.Small,
		cardFace.ImageURIs.Normal,
		cardFace.ImageURIs.Large,
		cardFace.ImageURIs.PNG,
		cardFace.ImageURIs.ArtCrop,
		cardFace.ImageURIs.BorderCrop,
		cardFace.ManaCost,
		cardFace.Name,
		cardFace.OracleText,
		cardFace.Power,
		cardFace.Toughness,
		cardFace.Loyalty,
		cardFace.TypeLine,
		cardFace.DerivedType,
		cardFace.Watermark,
	)
	if err != nil {
		return 0, err
	}

	cardFaceID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return cardFaceID, nil
}

// GenerateFacesJSON calculates card face info per card saves the result as JSON to the card row.
func (c *cardRepository) GenerateFacesJSON() error {
	_, err := c.db.Exec(`UPDATE cards c
		INNER JOIN (
			SELECT
			c.id card_id,
			JSON_ARRAYAGG(JSON_OBJECT(
				'face_id', f.id,
				'name', f.name,
				'mana_cost', f.mana_cost,
				'type_line', f.type_line,
				'derived_type', f.derived_type,
				'oracle_text', f.oracle_text,
				'flavor_text', f.flavor_text,
				'image', f.image_normal,
				'power', f.power,
				'toughness', f.toughness,
				'loyalty', f.loyalty,
				'artist', f.artist
			)) faces
			FROM cards c
			INNER JOIN card_faces f ON f.card_id = c.id
			GROUP BY c.id
		) a
		SET c.faces_json = a.faces
		WHERE c.id = a.card_id
	`)

	return err
}

// GenerateSetsJSON calculates card set info per card saves the result as JSON to the card row.
func (c *cardRepository) GenerateSetsJSON() error {
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`TRUNCATE TABLE card_sets_list`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT INTO card_sets_list (oracle_id, sets_json)
		SELECT
		a.oracle_id,
		JSON_ARRAYAGG(JSON_OBJECT(
			'set_name', a.set_name,
			'price', a.usd,
			'card_faces', a.faces_json
		)) sets
		FROM (
			SELECT
			c.oracle_id,
			c.set_name,
			p.usd,
			c.faces_json
			FROM cards c
			INNER JOIN card_faces f ON c.id = f.card_id
			INNER JOIN card_prices p ON p.card_id = c.id
			GROUP BY c.id
			ORDER BY c.released_at DESC
		) a
		GROUP BY a.oracle_id
	`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`UPDATE cards c
		INNER JOIN card_sets_list s ON s.oracle_id = c.oracle_id
		SET c.card_sets_list_id = s.id
		WHERE c.oracle_id = s.oracle_id
	`)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (c *cardRepository) InsertRulings(rulings []models.ScryfallRuling) error {
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}

	for _, ruling := range rulings {
		_, err := c.insertRuling(tx, ruling)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (c *cardRepository) insertRuling(tx *sql.Tx, ruling models.ScryfallRuling) (int64, error) {
	result, err := tx.Exec(`INSERT IGNORE INTO card_rulings (
		oracle_id,
		comment_hash,
		source,
		published_at,
		comment
	) VALUES (
		?,
		MD5(?),
		?,
		?,
		?
	)
	`,
		ruling.OracleID,
		ruling.Comment,
		ruling.Source,
		ruling.PublishedAt,
		ruling.Comment,
	)
	if err != nil {
		return 0, err
	}

	rulingID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return rulingID, nil
}

func contains(list []string, key string) bool {
	for _, item := range list {
		if item == key {
			return true
		}
	}

	return false
}
