package buff

import "time"

// BUFF CONFIGS

const (
	BUFF_LISTING_API            = "https://buff.163.com/api/market/goods/sell_order"
	BUFF_TRANSACTION_API        = "https://buff.163.com/api/market/goods/bill_order"
	BUFF_CSGO_NAME              = "csgo"
	BUFF_LISTING_ITEMS_PER_PAGE = 20

	BUFF_SLEEP_TIME_MIN_S = 15
	BUFF_SLEEP_TIME_MAX_S = 25

	BUFF_SLEEP_TIME_MIN = BUFF_SLEEP_TIME_MIN_S * time.Second
	BUFF_SLEEP_TIME_MAX = BUFF_SLEEP_TIME_MAX_S * time.Second
	BUFF_RAW_RES_DIR    = "output"
)

const (
	BUFF_SORTING_CREATED_AT = "created.desc"
	BUFF_SORTING_PRICE_ASC  = "price.asc"
)

var BUFF_HEADERS = map[string]string{
	"Accept":             "application/json, text/javascript, */*; q=0.01",
	"Accept-Encoding":    "gzip, deflate, br, zstd",
	"Accept-Language":    "en-HK,en;q=0.9,zh-HK;q=0.8,zh;q=0.7,en-GB;q=0.6,en-US;q=0.5,zh-CN;q=0.4,zh-TW;q=0.3",
	"Connection":         "keep-alive",
	"Host":               "buff.163.com",
	"Sec-Fetch-Dest":     "empty",
	"Sec-Fetch-Mode":     "cors",
	"Sec-Fetch-Site":     "same-origin",
	"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
	"X-Requested-With":   "XMLHttpRequest",
	"sec-ch-ua":          "\"Google Chrome\";v=\"123\", \"Not:A-Brand\";v=\"8\", \"Chromium\";v=\"123\"",
	"sec-ch-ua-mobile":   "?0",
	"sec-ch-ua-platform": "\"Windows\"",
}
