package repositories

import (
	"database/sql"

	"github.com/BrandonWade/blackblade-batch/models"
	"github.com/jmoiron/sqlx"
)

// CardRepository interface for working with a cardRepository
type CardRepository interface {
	UpsertCards(cards []models.ScryfallCard) error
	GenerateSetNameImageValues() error
}

type cardRepository struct {
	db *sqlx.DB
}

// NewCardRepository create a new CardRepository instance
func NewCardRepository(db *sqlx.DB) CardRepository {
	return &cardRepository{
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
		cardID, err := c.upsertCard(tx, card)
		if err != nil {
			return err
		}

		if cardID == 0 {
			continue
		}

		err = c.upsertCardMultiverseIDs(tx, cardID, card.MultiverseIDs)
		if err != nil {
			return err
		}

		err = c.upsertCardFrameEffects(tx, cardID, card.FrameEffects)
		if err != nil {
			return err
		}

		err = c.upsertCardPrices(tx, cardID, card.Prices)
		if err != nil {
			return err
		}

		cardFaces := c.getCardFaces(card)
		for _, cardFace := range cardFaces {
			cardFaceID, err := c.upsertCardFace(tx, cardID, cardFace)
			if err != nil {
				return err
			}

			err = c.upsertCardFaceColors(tx, cardFaceID, cardFace.Colors)
			if err != nil {
				return err
			}

			err = c.upsertCardFaceColorIndicators(tx, cardFaceID, cardFace.ColorIndicator)
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

func (c *cardRepository) upsertCard(tx *sql.Tx, card models.ScryfallCard) (int64, error) {
	result, err := tx.Exec(`INSERT INTO cards (
		scryfall_id,
		oracle_id,
		tcgplayer_id,
		card_back_id,`+
		"`set`,"+
		`set_name,
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
		?
	)`,
		card.ID,
		card.OracleID,
		card.TCGPlayerID,
		card.CardBackID,
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

	return cardID, nil
}

func (c *cardRepository) upsertCardMultiverseIDs(tx *sql.Tx, cardID int64, multiverseIDs []int) error {
	for _, multiverseID := range multiverseIDs {
		_, err := tx.Exec(`INSERT INTO card_multiverse_ids (
			card_id,
			multiverse_id
		) VALUES (
			?,
			?
		) ON DUPLICATE KEY UPDATE
			multiverse_id = ?
		`,
			cardID,
			multiverseID,
			multiverseID,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *cardRepository) upsertCardFrameEffects(tx *sql.Tx, cardID int64, frameEffects []string) error {
	for _, frameEffect := range frameEffects {
		_, err := tx.Exec(`INSERT INTO card_frame_effects (
			card_id,
			frame_effect
		) VALUES (
			?,
			?
		) ON DUPLICATE KEY UPDATE
			frame_effect = ?
		`,
			cardID,
			frameEffect,
			frameEffect,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *cardRepository) upsertCardPrices(tx *sql.Tx, cardID int64, prices models.ScryfallPrices) error {
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
		if card.Layout == "flip" || card.Layout == "split" || card.Layout == "adventure" {
			for i := range card.CardFaces {
				card.CardFaces[i].ImageURIs.Small = card.ImageURIs.Small
				card.CardFaces[i].ImageURIs.Normal = card.ImageURIs.Normal
				card.CardFaces[i].ImageURIs.Large = card.ImageURIs.Large
				card.CardFaces[i].ImageURIs.PNG = card.ImageURIs.PNG
				card.CardFaces[i].ImageURIs.ArtCrop = card.ImageURIs.ArtCrop
				card.CardFaces[i].ImageURIs.BorderCrop = card.ImageURIs.BorderCrop
			}
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
		TypeLine:        card.TypeLine,
		Watermark:       card.Watermark,
	}

	return []models.ScryfallCardFace{
		cardFace,
	}
}

func (c *cardRepository) upsertCardFace(tx *sql.Tx, cardID int64, cardFace models.ScryfallCardFace) (int64, error) {
	result, err := tx.Exec(`INSERT INTO card_faces (
		card_id,
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
		?
	)`,
		cardID,
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

func (c *cardRepository) upsertCardFaceColors(tx *sql.Tx, cardFaceID int64, colors []string) error {
	for _, color := range colors {
		_, err := tx.Exec(`INSERT INTO card_face_colors (
			card_face_id,
			color
		) VALUES (
			?,
			?
		) ON DUPLICATE KEY UPDATE
			color = ?
		`,
			cardFaceID,
			color,
			color,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *cardRepository) upsertCardFaceColorIndicators(tx *sql.Tx, cardFaceID int64, colorIndicators []string) error {
	for _, colorIndicator := range colorIndicators {
		_, err := tx.Exec(`INSERT INTO card_face_color_indicators (
			card_face_id,
			color_indicator
		) VALUES (
			?,
			?
		) ON DUPLICATE KEY UPDATE
			color_indicator = ?
		`,
			cardFaceID,
			colorIndicator,
			colorIndicator,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *cardRepository) GenerateSetNameImageValues() error {
	_, err := c.db.Exec(`UPDATE cards c
		INNER JOIN (
			SELECT
			c.oracle_id,
			JSON_ARRAYAGG(JSON_OBJECT('set_name', c.set_name, 'image', f.image_normal)) sets
			FROM cards c
			INNER JOIN card_faces f ON c.id = f.card_id
			GROUP BY c.oracle_id
		) a
		SET c.set_name_image_json = a.sets
		WHERE c.oracle_id = a.oracle_id
	`)

	return err
}
