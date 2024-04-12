package buff

import (
	"encoding/json"
	"fmt"
	"net/http"
	"steam-trading/shared"
	"steam-trading/shared/database/model"

	"github.com/mikezzb/steam-trading-crawler/types"
	"github.com/mikezzb/steam-trading-crawler/utils"
)

type BuffParser struct {
}

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

func itemToListing(item BuffItem) model.Listing {
	return model.Listing{
		Price:            item.Price,
		CreatedAt:        item.CreatedAt,
		UpdatedAt:        item.UpdatedAt,
		PreviewUrl:       item.PreviewUrl,
		GoodsId:          item.GoodsId,
		ClassId:          item.AssetInfo.ClassId,
		AssetId:          item.AssetInfo.AssetId,
		TradableCooldown: item.AssetInfo.TradableCooldown,
		PaintWear:        item.AssetInfo.PaintWear,
		PaintIndex:       item.AssetInfo.Info.PaintIndex,
		PaintSeed:        item.AssetInfo.Info.PaintSeed,
	}
}

func (p *BuffParser) formatItemListings(data BuffListingResponseData) ([]model.Listing, error) {
	items := data.Data.Items
	listing := make([]model.Listing, len(items))
	for i, item := range items {
		listing[i] = itemToListing(item)
	}
	return listing, nil
}

func getFirstValue(data map[string]BuffGoodsInfo) BuffGoodsInfo {
	for _, v := range data {
		return v
	}
	return BuffGoodsInfo{}
}

func ExtractLowestPrice(listing []model.Listing) string {
	lowestPrice := listing[0].Price
	for _, item := range listing {
		if item.Price < lowestPrice {
			lowestPrice = item.Price
		}
	}
	return lowestPrice
}

func (p *BuffParser) formatItem(data BuffListingResponseData, listings []model.Listing) (*model.Item, error) {
	item := getFirstValue(data.Data.GoodsInfos)

	formattedItems := model.Item{}
	formattedItems.IconUrl = item.IconUrl
	formattedItems.SteamPrice = item.SteamPrice
	formattedItems.LowestMarketPrice = ExtractLowestPrice(listings)
	formattedItems.LowestMarketName = shared.MARKET_NAME_BUFF

	return &formattedItems, nil
}

func (p *BuffParser) ParseItemListings(name string, bodyBytes []byte) (*types.ListingsData, error) {
	decodedReader, err := utils.ReadBytes(bodyBytes)
	if err != nil {
		fmt.Printf("Failed to decode response: %v\n", err)
		return nil, err
	}
	defer decodedReader.Close()

	// unmarshal response
	var data BuffListingResponseData
	if err := json.NewDecoder(decodedReader).Decode(&data); err != nil {
		fmt.Printf("Failed to unmarshal response: %v\n", err)
		return nil, err
	}

	// format data
	if listings, err := p.formatItemListings(data); err != nil {
		return nil, err
	} else if item, err := p.formatItem(data, listings); err != nil {
		return nil, err
	} else {
		formattedListings := utils.PostFormatListing(name, listings)
		return &types.ListingsData{
			Item:     item,
			Listings: formattedListings,
		}, nil
	}
}

func (p *BuffParser) ParseItemTransactions(resp *http.Response) ([]model.Transaction, error) {
	return nil, nil
}
