package igxe

import (
	"fmt"
	"log"
	"path"
	"strconv"

	shared "github.com/mikezzb/steam-trading-shared"
)

// consts

const (
	// {IGXE_LISTING_API}/{itemId}
	IGXE_LISTING_API = "https://www.igxe.cn/product/trade/730"
	// {IGXE_TRANSACTION_API}/{itemId}
	IGXE_TRANSACTION_API = "https://www.igxe.cn/product/get_product_sales_history/730"

	IGXE_LISTING_ITEMS_PER_PAGE = 10
)

// configs
const (
	IGXE_SLEEP_TIME_MIN_S = 13
	IGXE_SLEEP_TIME_MAX_S = 22

	IGXE_RAW_RES_DIR = "output"
)

var IGXE_HEADERS = map[string]string{
	"Accept":             "*/*",
	"Accept-Encoding":    "gzip, deflate, br, zstd",
	"Accept-Language":    "en-HK,en;q=0.9",
	"Content-Length":     "0",
	"Origin":             "https://www.igxe.cn",
	"Priority":           "u=1, i",
	"Sec-Ch-Ua":          "\"Chromium\";v=\"124\", \"Google Chrome\";v=\"124\", \"Not-A.Brand\";v=\"99\"",
	"Sec-Ch-Ua-Mobile":   "?0",
	"Sec-Ch-Ua-Platform": "\"macOS\"",
	"Sec-Fetch-Dest":     "empty",
	"Sec-Fetch-Mode":     "cors",
	"Sec-Fetch-Site":     "same-origin",
	"User-Agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
}

func getIgxeListingUrl(itemId int) string {
	return IGXE_LISTING_API + "/" + strconv.Itoa(itemId)
}

func getIgxeTransactionUrl(itemId int) string {
	return IGXE_TRANSACTION_API + "/" + strconv.Itoa(itemId)
}

func getIgxeSavePath(itemName string, pageNum int, label string) string {
	return path.Join(IGXE_RAW_RES_DIR, fmt.Sprintf("igxe_%s_%s_%d_%s.json", label, itemName, pageNum, shared.GetTimestampNow()))
}

func igxeLog(format string, v ...interface{}) {
	log.Printf("[igxe] "+format, v...)
}

func getRefererHeader(itemId int) map[string]string {
	itemUrl := fmt.Sprintf("https://www.igxe.cn/product/730/%d", itemId)
	return map[string]string{
		"Referer": itemUrl,
	}
}
