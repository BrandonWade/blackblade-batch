package repositories

import (
	"database/sql"

	"github.com/BrandonWade/blackblade-batch/models"
	"github.com/jmoiron/sqlx"
)

// CardRepository interface for working with a cardRepository
type CardRepository interface {
	UpsertCards(cards []models.ScryfallCard) (int64, error)
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
func (c *cardRepository) UpsertCards(cards []models.ScryfallCard) (int64, error) {
	total := int64(0)

	tx, err := c.db.Begin()
	if err != nil {
		return 0, err
	}

	for _, card := range cards {
		result, err := c.upsertCard(tx, card)
		if err != nil {
			return 0, err
		}

		cardID, err := result.LastInsertId()
		if err != nil {
			return 0, err
		}

		if cardID == 0 {
			continue
		}

		err = c.upsertCardMultiverseIDs(tx, cardID, card.MultiverseIDs)
		if err != nil {
			return 0, err
		}

		err = c.upsertCardFrameEffects(tx, cardID, card.FrameEffects)
		if err != nil {
			return 0, err
		}

		err = c.upsertCardPrices(tx, cardID, card.Prices)
		if err != nil {
			return 0, err
		}

		cardFaces := c.getCardFaces(card)
		for _, cardFace := range cardFaces {
			cardFaceID, err := c.upsertCardFace(tx, cardFace)
			if err != nil {
				return 0, err
			}

			err = c.upsertCardFaceColors(tx, cardFaceID, cardFace.Colors)
			if err != nil {
				return 0, err
			}

			err = c.upsertCardFaceColorIndicators(tx, cardFaceID, cardFace.ColorIndicator)
			if err != nil {
				return 0, err
			}
		}

		count, err := result.RowsAffected()
		if err != nil {
			return 0, err
		}

		// MySQL returns 1 row affected for an insert, and 2 for an update
		if count > 0 {
			total++
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (c *cardRepository) upsertCard(tx *sql.Tx, card models.ScryfallCard) (sql.Result, error) {
	return tx.Exec(`INSERT INTO cards (
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
}

func (c *cardRepository) upsertCardMultiverseIDs(tx *sql.Tx, cardID int64, multiverseIDs []int) error {
	// TODO: Optimize
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
	// TODO: Optimize
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
	// TODO: Get card faces (use either card.CardFaces if it exists, or build it from various card fields)
	return []models.ScryfallCardFace{}
}

func (c *cardRepository) upsertCardFace(tx *sql.Tx, cardFace models.ScryfallCardFace) (int64, error) {
	// TODO: Implement
	_, err := tx.Exec(``)

	return 0, err
}

func (c *cardRepository) upsertCardFaceColors(tx *sql.Tx, cardFaceID int64, colors []string) error {
	// TODO: Optimize
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
	// TODO: Optimize
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
