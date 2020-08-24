package models

// ScryfallCard represents a scryfall card
type ScryfallCard struct {
	Object          string              `json:"object"`
	ID              string              `json:"id"`
	OracleID        string              `json:"oracle_id"`
	MultiverseIDs   []int               `json:"multiverse_ids"`
	MtgoID          int64               `json:"mtgo_id"`
	MtgoFoilID      int64               `json:"mtgo_foil_id"`
	TCGPlayerID     int64               `json:"tcgplayer_id"`
	Name            string              `json:"name"`
	Lang            string              `json:"lang"`
	ReleasedAt      string              `json:"released_at"`
	URI             string              `json:"uri"`
	ScryfallURI     string              `json:"scryfall_uri"`
	Layout          string              `json:"layout"`
	HighresImage    bool                `json:"highres_image"`
	ImageURIs       ScryfallImageURIs   `json:"image_uris"`
	ManaCost        string              `json:"mana_cost"`
	CMC             float64             `json:"cmc"`
	TypeLine        string              `json:"type_line"`
	OracleText      string              `json:"oracle_text"`
	Power           string              `json:"power"`
	Toughness       string              `json:"toughness"`
	Loyalty         string              `json:"loyalty"`
	Colors          []string            `json:"colors"`
	ColorIdentity   []string            `json:"color_identity"`
	Keywords        []string            `json:"keywords"`
	Legalities      ScryfallLegalities  `json:"legalities"`
	Games           []string            `json:"games"`
	Reserved        bool                `json:"reserved"`
	Foil            bool                `json:"foil"`
	Nonfoil         bool                `json:"nonfoil"`
	Oversized       bool                `json:"oversized"`
	Promo           bool                `json:"promo"`
	Reprint         bool                `json:"reprint"`
	Variation       bool                `json:"variation"`
	Set             string              `json:"set"`
	SetName         string              `json:"set_name"`
	SetType         string              `json:"set_type"`
	SetURI          string              `json:"set_uri"`
	SetSearchURI    string              `json:"set_search_uri"`
	ScryfallSetURI  string              `json:"scryfall_set_uri"`
	RulingsURI      string              `json:"rulings_uri"`
	PrintsSearchURI string              `json:"prints_search_uri"`
	CollectorNumber string              `json:"collector_number"`
	Digital         bool                `json:"digital"`
	Rarity          string              `json:"rarity"`
	FlavorText      string              `json:"flavor_text"`
	CardBackID      string              `json:"card_back_id"`
	Artist          string              `json:"artist"`
	ArtistIDs       []string            `json:"artist_ids"`
	IllustrationID  string              `json:"illustration_id"`
	BorderColor     string              `json:"border_color"`
	Frame           string              `json:"frame"`
	FrameEffects    []string            `json:"frame_effects"`
	FullArt         bool                `json:"full_art"`
	Textless        bool                `json:"textless"`
	Booster         bool                `json:"booster"`
	StorySpotlight  bool                `json:"story_spotlight"`
	EDHRecRank      int64               `json:"edhrec_rank"`
	Prices          ScryfallPrices      `json:"prices"`
	RelatedURIs     ScryfallRelatedURIs `json:"related_uris"`
}

// ScryfallImageURIs represents a scryfall card's images
type ScryfallImageURIs struct {
	Small      string `json:"small"`
	Normal     string `json:"normal"`
	Large      string `json:"large"`
	PNG        string `json:"png"`
	ArtCrop    string `json:"art_crop"`
	BorderCrop string `json:"border_crop"`
}

// ScryfallLegalities represents a scryfall card's legalities
type ScryfallLegalities struct {
	Standard  string `json:"standard"`
	Future    string `json:"future"`
	Historic  string `json:"historic"`
	Pioneer   string `json:"pioneer"`
	Modern    string `json:"modern"`
	Legacy    string `json:"legacy"`
	Pauper    string `json:"pauper"`
	Vintage   string `json:"vintage"`
	Penny     string `json:"penny"`
	Commander string `json:"commander"`
	Brawl     string `json:"brawl"`
	Duel      string `json:"duel"`
	Oldschool string `json:"oldschool"`
}

// ScryfallPrices represents a scryfall card's prices
type ScryfallPrices struct {
	USD     string `json:"usd"`
	USDFoil string `json:"usd_foil"`
	EUR     string `json:"eur"`
	Tix     string `json:"tix"`
}

// ScryfallRelatedURIs represents a scryfall card's related URIs
type ScryfallRelatedURIs struct {
	Gatherer       string `json:"gatherer"`
	TCGPlayerDecks string `json:"tcgplayer_decks"`
	EDHRec         string `json:"edhrec"`
	MTGTop8        string `json:"mtgtop8"`
}
