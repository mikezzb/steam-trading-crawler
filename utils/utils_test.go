package utils_test

import (
	"os"
	"testing"

	"github.com/mikezzb/steam-trading-crawler/types"
	"github.com/mikezzb/steam-trading-crawler/utils"
	shared "github.com/mikezzb/steam-trading-shared"
)

func TestImageDownload(t *testing.T) {
	t.Run("ImageDownload", func(t *testing.T) {
		testUrl := "https://market.fp.ps.netease.com/file/65f8313923c06400e90f6fdb6Nsdj02W05"
		err := utils.DownloadImage(testUrl, "test.png")
		if err != nil {
			t.Error(err)
		}

		// check if the file exists
		_, err = os.Stat("test.png")
		if err != nil {
			t.Error(err)
		}

		// clean up
		err = os.Remove("test.png")
		if err != nil {
			t.Error(err)
		}
	})
}

func TestPostFormat(t *testing.T) {
	t.Run("GetPrice", func(t *testing.T) {
		listings := []types.Listing{
			{
				Price: "900",
			},
			{
				Price: "2400",
			},
			{
				Price: "10000",
			},
			{
				Price: "10002",
			},
			{
				Price: "10400",
			},
		}

		lowestPrice := utils.ExtractLowestPrice(listings)

		if shared.NumStrCmp(lowestPrice, listings[0].Price) != 0 {
			t.Errorf("Expected %v, got %v", listings[0].Price, lowestPrice)
		}
	})
}
