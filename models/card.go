package models

// Card represents a row from the blackblade.cards table
type Card struct {
	ID              int64   `json:"id" db:"id"`
	ScryfallID      string  `json:"scryfall_id,omitempty" db:"scryfall_id"`
	OracleID        string  `json:"oracle_id,omitempty" db:"oracle_id"`
	Name            string  `json:"name,omitempty" db:"name"`
	Lang            string  `json:"lang,omitempty" db:"lang"`
	ReleasedAt      string  `json:"released_at,omitempty" db:"released_at"`
	URI             string  `json:"uri,omitempty" db:"uri"`
	ScryfallURI     string  `json:"scryfall_uri,omitempty" db:"scryfall_uri"`
	Layout          string  `json:"layout,omitempty" db:"layout"`
	HighResImage    bool    `json:"highres_image,omitempty" db:"highres_image"`
	ManaCost        string  `json:"mana_cost,omitempty" db:"mana_cost"`
	CMC             float64 `json:"cmc,omitempty" db:"cmc"`
	TypeLine        string  `json:"type_line,omitempty" db:"type_line"`
	OracleText      string  `json:"oracle_text,omitempty" db:"oracle_text"`
	Power           string  `json:"power,omitempty" db:"power"`
	Toughness       string  `json:"toughness,omitempty" db:"toughness"`
	Loyalty         string  `json:"loyalty,omitempty" db:"loyalty"`
	Reserved        bool    `json:"reserved,omitempty" db:"reserved"`
	Foil            bool    `json:"foil,omitempty" db:"foil"`
	NonFoil         bool    `json:"nonfoil,omitempty" db:"nonfoil"`
	Oversized       bool    `json:"oversized,omitempty" db:"oversized"`
	Promo           bool    `json:"promo,omitempty" db:"promo"`
	Reprint         bool    `json:"reprint,omitempty" db:"reprint"`
	Variation       bool    `json:"variation,omitempty" db:"variation"`
	Set             string  `json:"set,omitempty" db:"set"`
	SetName         string  `json:"set_name,omitempty" db:"set_name"`
	SetType         string  `json:"set_type,omitempty" db:"set_type"`
	SetURI          string  `json:"set_uri,omitempty" db:"set_uri"`
	SetSearchURI    string  `json:"set_search_uri,omitempty" db:"set_search_uri"`
	ScryfallSetURI  string  `json:"scryfall_set_uri,omitempty" db:"scryfall_set_uri"`
	RulingsURI      string  `json:"rulings_uri,omitempty" db:"rulings_uri"`
	PrintsSearchURI string  `json:"prints_search_uri,omitempty" db:"prints_search_uri"`
	CollectorNumber string  `json:"collector_number,omitempty" db:"collector_number"`
	Digital         bool    `json:"digital,omitempty" db:"digital"`
	Rarity          string  `json:"rarity,omitempty" db:"rarity"`
	CardBackID      string  `json:"card_back_id,omitempty" db:"card_back_id"`
	Artist          string  `json:"artist,omitempty" db:"artist"`
	IllustrationID  string  `json:"illustration_id,omitempty" db:"illustration_id"`
	BorderColor     string  `json:"border_color,omitempty" db:"border_color"`
	Frame           string  `json:"frame,omitempty" db:"frame"`
	FullArt         bool    `json:"full_art,omitempty" db:"full_art"`
	Textless        bool    `json:"textless,omitempty" db:"textless"`
	Booster         bool    `json:"booster,omitempty" db:"booster"`
	StorySpotlight  bool    `json:"story_spotlight,omitempty" db:"story_spotlight"`
}
