package buff

type BuffGoodsInfo struct {
	AppID      int    `json:"appid"`
	IconUrl    string `json:"icon_url"`
	SteamPrice string `json:"steam_price"`
}

type BuffItemAssetInfoInfo struct {
	PaintIndex int `json:"paintindex"`
	PaintSeed  int `json:"paintseed"`
}

type BuffItemAssetInfo struct {
	Info             BuffItemAssetInfoInfo `json:"info"`
	ClassId          string                `json:"classid"`
	AssetId          string                `json:"assetid"` // for steam preview
	PaintWear        string                `json:"paintwear"`
	TradableCooldown string                `json:"tradable_cooldown_text"`
}

type BuffItem struct {
	Price      string            `json:"price"`
	AssetInfo  BuffItemAssetInfo `json:"asset_info"`
	CreatedAt  int               `json:"created_at"`
	UpdatedAt  int               `json:"updated_at"`
	PreviewUrl string            `json:"img_src"`
	GoodsId    int               `json:"goods_id"`
}

type BuffListingResponseData struct {
	Code string `json:"code"`
	Data struct {
		GoodsInfos map[string]BuffGoodsInfo `json:"goods_infos"`
		Items      []BuffItem               `json:"items"`
	} `json:"data"`
}

type BuffTransactionResponseData struct {
	Code string `json:"code"`
	Data struct {
		Items []BuffItem `json:"items"`
	} `json:"data"`
}
