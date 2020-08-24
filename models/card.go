package models

import "github.com/gobuffalo/nulls"

// Card represents a row from the blackblade.cards table
type Card struct {
	ID              int64               `json:"id,omitempty" db:"id"`
	ScryfallID      string              `json:"scryfall_id,omitempty" db:"scryfall_id"`
	OracleID        string              `json:"oracle_id,omitempty" db:"oracle_id"`
	TCGPlayerID     nulls.Int64         `json:"tcgplayer_id,omitempty" db:"tcgplayer_id"`
	CardBackID      string              `json:"card_back_id,omitempty" db:"card_back_id"`
	MultiverseIDs   []CardMultiverseIDs `json:"multiverse_ids,omitempty" db:"multiverse_ids"`
	Set             string              `json:"set,omitempty" db:"set"`
	SetName         string              `json:"set_name,omitempty" db:"set_name"`
	Rarity          string              `json:"rarity,omitempty" db:"rarity"`
	Layout          string              `json:"layout,omitempty" db:"layout"`
	BorderColor     string              `json:"border_color,omitempty" db:"border_color"`
	Frame           string              `json:"frame,omitempty" db:"frame"`
	ReleasedAt      nulls.String        `json:"released_at,omitempty" db:"released_at"`
	Legalities      CardLegalities      `json:"legalities,omitempty" db:"card_legalities"`
	FrameEffects    CardFrameEffects    `json:"frame_effects,omitempty" db:"frame_effects"`
	Prices          CardPrices          `json:"prices,omitempty" db:"card_prices"`
	HasFoil         bool                `json:"has_foil,omitempty" db:"foil"`
	HasNonFoil      bool                `json:"has_non_foil,omitempty" db:"non_foil"`
	IsOversized     bool                `json:"nonfoil,omitempty" db:"is_oversized"`
	IsReserved      bool                `json:"reserved,omitempty" db:"is_reserved"`
	IsBooster       bool                `json:"booster,omitempty" db:"is_booster"`
	IsDigitalOnly   bool                `json:"digital,omitempty" db:"is_digital"`
	IsFullArt       bool                `json:"full_art,omitempty" db:"is_full_art"`
	IsTextless      bool                `json:"textless,omitempty" db:"is_textless"`
	IsReprint       bool                `json:"reprint,omitempty" db:"is_reprint"`
	HasHighresImage bool                `json:"highres_image,omitempty" db:"has_highres_image"`
	RulingsURI      string              `json:"rulings_uri,omitempty" db:"rulings_uri"`
	ScryfallURI     string              `json:"scryfall_uri,omitempty" db:"scryfall_uri"`
}

// CardMultiverseIDs represents a row from the blackblade.multiverse_ids table
type CardMultiverseIDs struct {
	ID           int64 `json:"id,omitempty" db:"id"`
	CardID       int64 `json:"card_id,omitempty" db:"card_id"`
	MultiverseID int64 `json:"multiverse_id,omitempty" db:"multiverse_id"`
}

// CardLegalities represents a row from the blackblade.card_legalities table
type CardLegalities struct {
	ID        int64  `json:"id,omitempty" db:"id"`
	CardID    int64  `json:"card_id,omitempty" db:"card_id"`
	Standard  string `json:"standard,omitempty" db:""`
	Future    string `json:"future,omitempty" db:""`
	Historic  string `json:"historic,omitempty" db:""`
	Pioneer   string `json:"pioneer,omitempty" db:""`
	Modern    string `json:"modern,omitempty" db:""`
	Legacy    string `json:"legacy,omitempty" db:""`
	Pauper    string `json:"pauper,omitempty" db:""`
	Vintage   string `json:"vintage,omitempty" db:""`
	Penny     string `json:"penny,omitempty" db:""`
	Commander string `json:"commander,omitempty" db:""`
	Brawl     string `json:"brawl,omitempty" db:""`
	Duel      string `json:"duel,omitempty" db:""`
	Oldschool string `json:"oldschool,omitempty" db:""`
}

// CardFrameEffects represents a row from the blackblade.card_frame_effects table
type CardFrameEffects struct {
	ID          int64  `json:"id,omitempty" db:"id"`
	CardID      int64  `json:"card_id,omitempty" db:"card_id"`
	FrameEffect string `json:",omitempty" db:""`
}

// CardPrices represents a row from the blackblade.card_prices table
type CardPrices struct {
	ID      int64  `json:"id,omitempty" db:"id"`
	CardID  int64  `json:"card_id,omitempty" db:"card_id"`
	USD     string `json:",omitempty" db:""`
	USDFoil string `json:",omitempty" db:""`
	EUR     string `json:",omitempty" db:""`
	TIX     string `json:",omitempty" db:""`
}
