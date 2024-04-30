package handler

import (
	shared "github.com/mikezzb/steam-trading-shared"
	"github.com/mikezzb/steam-trading-shared/database/model"
)

type Formatter struct{}

func (f *Formatter) FormatItem(item *model.Item) {
	item.ID = shared.GetItemId(item.Name)
}

func NewFormatter() *Formatter {
	return &Formatter{}
}
