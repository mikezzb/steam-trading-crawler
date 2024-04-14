package utils_test

import (
	"os"
	"testing"

	"github.com/mikezzb/steam-trading-crawler/utils"
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
