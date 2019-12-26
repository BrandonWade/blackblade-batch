package repositories

import (
	"github.com/BlueMonday/go-scryfall"
	"github.com/jmoiron/sqlx"
)

// CardRepository interface for working with a cardRepository
type CardRepository interface {
	UpsertCards([]scryfall.Card) (int64, error)
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
func (c *cardRepository) UpsertCards(cards []scryfall.Card) (int64, error) {
	total := int64(0)

	tx, err := c.db.Begin()
	if err != nil {
		return 0, err
	}

	for _, card := range cards {
		// TODO: Add missing fields
		result, err := tx.Exec(`INSERT INTO cards (
			scryfall_id,
			oracle_id,
			name,
			lang,
			uri,
			scryfall_uri,
			layout,
			highres_image,
			mana_cost,
			cmc,
			type_line,
			oracle_text,
			power,
			toughness,
			loyalty,
			reserved,
			foil,
			nonfoil,
			oversized,
			promo,
			reprint,
			set,
			set_name,
			set_uri,
			set_search_uri,
			scryfall_set_uri,
			rulings_uri,
			prints_search_uri,
			collector_number,
			digital,
			rarity,
			artist,
			illustration_id,
			border_color,
			frame,
			full_art
		) values (
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
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?
		)
		ON DUPLICATE KEY UPDATE
		`,
			card.ID,
			card.OracleID,
			card.Name,
			card.Lang,
			// card.ReleasedAt,
			card.URI,
			card.ScryfallURI,
			card.Layout,
			card.HighresImage,
			card.ManaCost,
			card.CMC,
			card.TypeLine,
			card.OracleText,
			card.Power,
			card.Toughness,
			card.Loyalty,
			card.Reserved,
			card.Foil,
			card.NonFoil,
			card.Oversized,
			card.Promo,
			card.Reprint,
			// card.Variation,
			card.Set,
			card.SetName,
			// card.SetType,
			card.SetURI,
			card.SetSearchURI,
			card.ScryfallSetURI,
			card.RulingsURI,
			card.PrintsSearchURI,
			card.CollectorNumber,
			card.Digital,
			card.Rarity,
			// card.CardBackID,
			card.Artist,
			card.IllustrationID,
			card.BorderColor,
			card.Frame,
			card.FullArt,
			// card.Textless,
			// card.Booster,
			// card.StorySpotlight,
		)
		if err != nil {
			return 0, err
		}

		// TODO: Insert into additional tables

		count, err := result.RowsAffected()
		if err != nil {
			return 0, err
		}

		total += count
	}

	err = tx.Commit()
	if err != nil {

	}

	return total, nil
}
