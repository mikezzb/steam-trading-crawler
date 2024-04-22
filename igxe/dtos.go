package igxe

type IgxePageInfo struct {
	PageCount int `json:"page_count"`
	Total     int `json:"total"`
	PageNo    int `json:"page_no"`
}

type IgxeListing struct {
	LastUpdated string `json:"last_updated"`
	// for share and inspect
	ID              int    `json:"id"`
	PaintType       int    `json:"paint_type"`
	PaintIndex      int    `json:"paint_index"`
	ExteriorWear    string `json:"exterior_wear"`
	Price           string `json:"unit_price"`
	SteamPid        string `json:"steam_pid"` // asset id?
	PaintSeed       int    `json:"paint_seed"`
	PaintSeedName   string `json:"paint_seed_name"` // rarity
	ReferencePrice  string `json:"reference_price"`
	InspectImgSmall string `json:"inspect_img_small"`
	LockSpanStr     string `json:"lock_span_str"`
}

type IgxeListingResponseData struct {
	Succ     bool          `json:"succ"`
	Page     IgxePageInfo  `json:"page"`
	Listings []IgxeListing `json:"d_list"`
}

// missing some fields compared to IgxeListing
type IgxeTransaction struct {
	LastUpdated string `json:"last_updated"`
	// for share and inspect
	ID            int    `json:"id"`
	PaintType     int    `json:"paint_type"`
	PaintIndex    int    `json:"paint_index"`
	ExteriorWear  string `json:"exterior_wear"`
	Price         string `json:"unit_price"`
	PaintSeed     int    `json:"paint_seed"`
	PaintSeedName string `json:"paint_seed_name"`
}

type IgxeTransactionResponseData struct {
	Succ bool              `json:"succ"`
	Data []IgxeTransaction `json:"data"`
}
