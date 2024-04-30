package types

import (
	"time"
)

type MarketPrice struct {
	Price     string    `bson:"price" json:"price"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

type Item struct {
	Name    string `bson:"name" json:"name"`
	IconUrl string `bson:"iconUrl" json:"iconUrl"`

	// Market prices
	BuffPrice  *MarketPrice `bson:"buffPrice,omitempty" json:"buffPrice"`
	UUPrice    *MarketPrice `bson:"uuPrice,omitempty" json:"uuPrice"`
	IgxePrice  *MarketPrice `bson:"igxePrice,omitempty" json:"igxePrice"`
	SteamPrice *MarketPrice `bson:"steamPrice,omitempty" json:"steamPrice"`
}

type Listing struct {
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
	CheckedAt time.Time `bson:"checkedAt" json:"checkedAt"`

	Name             string `bson:"name" json:"name"`
	Market           string `bson:"market" json:"market"`
	Price            string `bson:"price" json:"price"`
	PreviewUrl       string `bson:"previewUrl" json:"previewUrl"`
	GoodsId          int    `bson:"goodsId" json:"goodsId"`
	ClassId          string `bson:"classId" json:"classId"`
	AssetId          string `bson:"assetId" json:"assetId"`
	TradableCooldown string `bson:"tradableCooldown" json:"tradableCooldown"`

	PaintWear  string `bson:"paintWear" json:"paintWear"`
	PaintIndex int    `bson:"paintIndex" json:"paintIndex"`
	PaintSeed  int    `bson:"paintSeed" json:"paintSeed"`
	Rarity     string `bson:"rarity" json:"rarity"`

	// Market specific ID
	InstanceId string `bson:"instanceId" json:"instanceId"`
}

// Currently same as Listing
type Transaction struct {
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`

	Name    string `bson:"name" json:"name"`
	Market  string `bson:"market" json:"market"`
	AssetId string `bson:"assetId" json:"assetId"`

	Price            string `bson:"price" json:"price"`
	PreviewUrl       string `bson:"previewUrl" json:"previewUrl"`
	GoodsId          int    `bson:"goodsId" json:"goodsId"`
	ClassId          string `bson:"classId" json:"classId"`
	TradableCooldown string `bson:"tradableCooldown" json:"tradableCooldown"`

	PaintWear  string `bson:"paintWear" json:"paintWear"`
	PaintIndex int    `bson:"paintIndex" json:"paintIndex"`
	PaintSeed  int    `bson:"paintSeed" json:"paintSeed"`

	Rarity string `bson:"rarity" json:"rarity"`

	// market specific unique id
	InstanceId string `bson:"instanceId" json:"instanceId"`
}
