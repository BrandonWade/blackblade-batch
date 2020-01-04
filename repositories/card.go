package repositories

import (
	"database/sql"

	scryfall "github.com/BlueMonday/go-scryfall"
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

		if card.ImageURIs != nil {
			err = c.upsertCardImageURIs(tx, cardID, card.ImageURIs)
			if err != nil {
				return 0, err
			}
		}

		err = c.upsertCardColors(tx, cardID, card.Colors)
		if err != nil {
			return 0, err
		}

		err = c.upsertCardColorIdentities(tx, cardID, card.ColorIdentity)
		if err != nil {
			return 0, err
		}

		err = c.upsertCardLegalities(tx, cardID, card.Legalities)
		if err != nil {
			return 0, err
		}

		// err = c.upsertCardGames(tx, cardID, card.Games)
		// if err != nil {
		// 	return 0, err
		// }

		// err = c.upsertCardArtistIDs(tx, cardID, card.ArtistIDs)
		// if err != nil {
		// 	return 0, err
		// }

		err = c.upsertCardFrameEffects(tx, cardID, card.FrameEffects)
		if err != nil {
			return 0, err
		}

		// err = c.upsertCardPromoTypes(tx, cardID, card.PromoTypes)
		// if err != nil {
		// 	return 0, err
		// }

		err = c.upsertCardPreview(tx, cardID, card.Preview)
		if err != nil {
			return 0, err
		}

		err = c.upsertCardPrices(tx, cardID, card.Prices)
		if err != nil {
			return 0, err
		}

		err = c.upsertCardRelatedURIs(tx, cardID, card.RelatedURIs)
		if err != nil {
			return 0, err
		}

		err = c.upsertCardPurchaseURIs(tx, cardID, card.PurchaseURIs)
		if err != nil {
			return 0, err
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

func (c *cardRepository) upsertCard(tx *sql.Tx, card scryfall.Card) (sql.Result, error) {
	// TODO: Add missing fields
	return tx.Exec(`INSERT INTO cards (
		scryfall_id,
		oracle_id,
		name,
		lang,
		uri,
		scryfall_uri,
		layout,
		has_highres_image,
		mana_cost,
		cmc,
		type_line,
		oracle_text,
		power,
		toughness,
		loyalty,
		is_reserved,
		has_foil,
		has_nonfoil,
		is_oversized,
		is_promo,
		is_reprint,`+
		"`set`,"+
		`set_name,
		set_uri,
		set_search_uri,
		scryfall_set_uri,
		rulings_uri,
		prints_search_uri,
		collector_number,
		is_digital,
		rarity,
		artist,
		illustration_id,
		border_color,
		frame,
		is_full_art
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
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?
	) ON DUPLICATE KEY UPDATE
		name = ?,
		uri = ?,
		scryfall_uri = ?,
		mana_cost = ?,
		cmc = ?,
		oracle_text = ?,
		set_uri = ?,
		set_search_uri = ?,
		scryfall_set_uri = ?,
		rulings_uri = ?,
		prints_search_uri = ?,
		is_digital = ?
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
		card.Name,
		card.URI,
		card.ScryfallURI,
		card.ManaCost,
		card.CMC,
		card.OracleText,
		card.SetURI,
		card.SetSearchURI,
		card.ScryfallSetURI,
		card.RulingsURI,
		card.PrintsSearchURI,
		card.Digital,
	)
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

func (c *cardRepository) upsertCardImageURIs(tx *sql.Tx, cardID int64, imageURIs *scryfall.ImageURIs) error {
	upsertHelper := func(imageType, uri string) error {
		_, err := tx.Exec(`INSERT INTO card_image_uris (
			card_id,
			image_type,
			uri
		) VALUES (
			?,
			?,
			?
		) ON DUPLICATE KEY UPDATE
			image_type = ?,
			uri = ?`,
			cardID,
			imageType,
			uri,
			imageType,
			uri,
		)

		return err
	}

	err := upsertHelper("png", imageURIs.PNG)
	if err != nil {
		return err
	}

	err = upsertHelper("border_crop", imageURIs.BorderCrop)
	if err != nil {
		return err
	}

	err = upsertHelper("art_crop", imageURIs.ArtCrop)
	if err != nil {
		return err
	}

	err = upsertHelper("large", imageURIs.Large)
	if err != nil {
		return err
	}

	err = upsertHelper("normal", imageURIs.Normal)
	if err != nil {
		return err
	}

	err = upsertHelper("small", imageURIs.Small)
	if err != nil {
		return err
	}

	return nil
}

func (c *cardRepository) upsertCardColors(tx *sql.Tx, cardID int64, colors []scryfall.Color) error {
	for _, color := range colors {
		_, err := tx.Exec(`INSERT INTO card_colors (
			card_id,
			color
		) VALUES (
			?,
			?
		) ON DUPLICATE KEY UPDATE
			color = ?
		`,
			cardID,
			color,
			color,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *cardRepository) upsertCardColorIdentities(tx *sql.Tx, cardID int64, colorIdentities []scryfall.Color) error {
	for _, colorIdentity := range colorIdentities {
		_, err := tx.Exec(`INSERT INTO card_color_identities (
			card_id,
			color
		) VALUES (
			?,
			?
		) ON DUPLICATE KEY UPDATE
			color = ?
		`,
			cardID,
			colorIdentity,
			colorIdentity,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *cardRepository) upsertCardLegalities(tx *sql.Tx, cardID int64, legalities scryfall.Legalities) error {
	upsertHelper := func(format, legality scryfall.Legality) error {
		_, err := tx.Exec(`INSERT INTO card_legalities (
			card_id,
			format,
			legality
		) VALUES (
			?,
			?,
			?
		) ON DUPLICATE KEY UPDATE
			format = ?,
			legality = ?
		`,
			cardID,
			format,
			legality,
			format,
			legality,
		)

		return err
	}

	err := upsertHelper("standard", legalities.Standard)
	if err != nil {
		return err
	}

	err = upsertHelper("future", legalities.Future)
	if err != nil {
		return err
	}

	// err = upsertHelper("historic", legalities.Historic)
	// if err != nil {
	// 	return err
	// }

	err = upsertHelper("pioneer", legalities.Pioneer)
	if err != nil {
		return err
	}

	err = upsertHelper("modern", legalities.Modern)
	if err != nil {
		return err
	}

	err = upsertHelper("legacy", legalities.Legacy)
	if err != nil {
		return err
	}

	err = upsertHelper("pauper", legalities.Pauper)
	if err != nil {
		return err
	}

	// err = upsertHelper("vintage", legalities.Vintage)
	// if err != nil {
	// 	return err
	// }

	err = upsertHelper("penny", legalities.Penny)
	if err != nil {
		return err
	}

	err = upsertHelper("commander", legalities.Commander)
	if err != nil {
		return err
	}

	// err = upsertHelper("brawl", legalities.Brawl)
	// if err != nil {
	// 	return err
	// }

	err = upsertHelper("duel", legalities.Duel)
	if err != nil {
		return err
	}

	// err = upsertHelper("oldschool", legalities.OldSchool)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (c *cardRepository) upsertCardGames(tx *sql.Tx, cardID int64, games []string) error {
	for _, game := range games {
		_, err := tx.Exec(`INSERT INTO card_games (
			card_id,
			game
		) VALUES (
			?,
			?
		) ON DUPLICATE KEY UPDATE
			game = ?
		`,
			cardID,
			game,
			game,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *cardRepository) upsertCardArtistIDs(tx *sql.Tx, cardID int64, artistIDs []string) error {
	for _, artistID := range artistIDs {
		_, err := tx.Exec(`INSERT INTO card_artist_ids (
			card_id,
			artist_id
		) VALUES (
			?,
			?
		) ON DUPLICATE KEY UPDATE
			artist_id = ?
		`,
			cardID,
			artistID,
			artistID,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *cardRepository) upsertCardFrameEffects(tx *sql.Tx, cardID int64, frameEffects []scryfall.FrameEffect) error {
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

func (c *cardRepository) upsertCardPromoTypes(tx *sql.Tx, cardID int64, promoTypes []string) error {
	for _, promoType := range promoTypes {
		_, err := tx.Exec(`INSERT INTO card_promo_types (
			card_id,
			promo_type
		) VALUES (
			?,
			?
		) ON DUPLICATE KEY UPDATE
			promo_type = ?
		`,
			cardID,
			promoType,
			promoType,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *cardRepository) upsertCardPreview(tx *sql.Tx, cardID int64, preview scryfall.Preview) error {
	_, err := tx.Exec(`INSERT INTO card_previews (
		card_id,
		source,
		source_uri
	) VALUES (
		?,
		?,
		?
	) ON DUPLICATE KEY UPDATE
		source = ?,
		source_uri = ?`,
		cardID,
		preview.Source,
		preview.SourceURI,
		// preview.PreviewedAt,
		preview.Source,
		preview.SourceURI,
		// preview.PreviewedAt,
	)

	return err
}

func (c *cardRepository) upsertCardPrices(tx *sql.Tx, cardID int64, prices scryfall.Prices) error {
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

func (c *cardRepository) upsertCardRelatedURIs(tx *sql.Tx, cardID int64, relatedURIs scryfall.RelatedURIs) error {
	_, err := tx.Exec(`INSERT INTO card_related_uris (
		card_id,
		tcgplayer_decks,
		edhrec,
		mtgtop8
	) VALUES (
		?,
		?,
		?,
		?
	) ON DUPLICATE KEY UPDATE
		tcgplayer_decks = ?,
		edhrec = ?,
		mtgtop8 = ?
	`,
		cardID,
		relatedURIs.TCGPlayerDecks,
		relatedURIs.EDHREC,
		relatedURIs.MTGTop8,
		relatedURIs.TCGPlayerDecks,
		relatedURIs.EDHREC,
		relatedURIs.MTGTop8,
	)

	return err
}

func (c *cardRepository) upsertCardPurchaseURIs(tx *sql.Tx, cardID int64, purchaseURIs scryfall.PurchaseURIs) error {
	_, err := tx.Exec(`INSERT INTO card_purchase_uris (
		card_id,
		tcgplayer,
		cardmarket,
		cardhoarder
	) VALUES (
		?,
		?,
		?,
		?
	) ON DUPLICATE KEY UPDATE
		tcgplayer = ?,
		cardmarket = ?,
		cardhoarder = ?
	`,
		cardID,
		purchaseURIs.TCGPlayer,
		purchaseURIs.CardMarket,
		purchaseURIs.CardHoarder,
		purchaseURIs.TCGPlayer,
		purchaseURIs.CardMarket,
		purchaseURIs.CardHoarder,
	)

	return err
}
