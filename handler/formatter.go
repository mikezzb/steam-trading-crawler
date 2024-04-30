package handler

import (
	"time"

	"github.com/mikezzb/steam-trading-crawler/types"
	shared "github.com/mikezzb/steam-trading-shared"
	"github.com/mikezzb/steam-trading-shared/database/model"
)

// DB formatter (adapter): format the data from crawler to db model
type Formatter struct{}

var instance *Formatter

func (f *Formatter) FormatItem(item *types.Item) *model.Item {
	category, skin, exterior := shared.DecodeItemFullName(item.Name)
	itemModel := &model.Item{
		ID:       shared.GetItemId(item.Name),
		Name:     item.Name,
		Category: category,
		Skin:     skin,
		Exterior: exterior,

		IconUrl: item.IconUrl,
	}
	if item.IgxePrice != nil {
		itemModel.IgxePrice = &model.MarketPrice{
			Price:     shared.GetDecimal128(item.IgxePrice.Price),
			UpdatedAt: item.IgxePrice.UpdatedAt,
		}
	}
	if item.BuffPrice != nil {
		itemModel.BuffPrice = &model.MarketPrice{
			Price:     shared.GetDecimal128(item.BuffPrice.Price),
			UpdatedAt: item.BuffPrice.UpdatedAt,
		}
	}
	if item.SteamPrice != nil {
		itemModel.SteamPrice = &model.MarketPrice{
			Price:     shared.GetDecimal128(item.SteamPrice.Price),
			UpdatedAt: item.SteamPrice.UpdatedAt,
		}
	}
	if item.UUPrice != nil {
		itemModel.UUPrice = &model.MarketPrice{
			Price:     shared.GetDecimal128(item.UUPrice.Price),
			UpdatedAt: item.UUPrice.UpdatedAt,
		}
	}
	return itemModel
}

func (f *Formatter) FormatListing(item *types.Listing) *model.Listing {
	return &model.Listing{
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
		CheckedAt: time.Now(),

		Name:             item.Name,
		Market:           item.Market,
		Price:            shared.GetDecimal128(item.Price),
		PreviewUrl:       item.PreviewUrl,
		GoodsId:          item.GoodsId,
		ClassId:          item.ClassId,
		AssetId:          item.AssetId,
		TradableCooldown: item.TradableCooldown,

		PaintWear:  shared.GetDecimal128(item.PaintWear),
		PaintIndex: item.PaintIndex,
		PaintSeed:  item.PaintSeed,
		Rarity:     item.Rarity,

		InstanceId: item.InstanceId,
	}
}

func (f *Formatter) FormatTransaction(item *types.Transaction) *model.Transaction {
	return &model.Transaction{
		Metadata: model.TransactionMetadata{
			AssetId: item.AssetId,
			Market:  item.Market,
		},

		Name:      item.Name,
		CreatedAt: item.CreatedAt,

		Price:            shared.GetDecimal128(item.Price),
		PreviewUrl:       item.PreviewUrl,
		GoodsId:          item.GoodsId,
		ClassId:          item.ClassId,
		TradableCooldown: item.TradableCooldown,

		PaintWear:  shared.GetDecimal128(item.PaintWear),
		PaintIndex: item.PaintIndex,
		PaintSeed:  item.PaintSeed,
		Rarity:     item.Rarity,

		InstanceId: item.InstanceId,
	}
}

func (f *Formatter) FormatListings(listings []types.Listing) []model.Listing {
	formattedListings := make([]model.Listing, len(listings))
	for i, item := range listings {
		formattedListings[i] = *f.FormatListing(&item)
	}
	return formattedListings
}

func (f *Formatter) FormatTransactions(transactions []types.Transaction) []model.Transaction {
	formattedTransactions := make([]model.Transaction, len(transactions))
	for i, item := range transactions {
		formattedTransactions[i] = *f.FormatTransaction(&item)
	}
	return formattedTransactions
}

func NewFormatter() *Formatter {
	if instance == nil {
		instance = &Formatter{}
	}
	return instance
}
