package buff

import (
	"fmt"
	"net/url"
	"path"
	"strconv"

	"github.com/mikezzb/steam-trading-crawler/crawler"
	"github.com/mikezzb/steam-trading-crawler/errors"
	"github.com/mikezzb/steam-trading-crawler/types"
	"github.com/mikezzb/steam-trading-crawler/utils"
	shared "github.com/mikezzb/steam-trading-shared"
)

type BuffCrawler struct {
	parser  *BuffParser
	crawler *crawler.Crawler
}

func (c *BuffCrawler) Stop() {
	c.crawler.Stop()
}

func (c *BuffCrawler) GetCookies() (string, error) {
	return c.crawler.GetCookies()
}

func NewCrawler(cookie string) (*BuffCrawler, error) {
	config := &crawler.CrawlerConfig{
		Cookie:      cookie,
		AuthUrls:    []string{BUFF_LISTING_API, BUFF_TRANSACTION_API},
		SleepMinSec: BUFF_SLEEP_TIME_MIN_S,
		SleepMaxSec: BUFF_SLEEP_TIME_MAX_S,
		Headers:     BUFF_HEADERS,
	}
	crawler, err := crawler.NewCrawler(config)
	if err != nil {
		return nil, err
	}

	parser := &BuffParser{}

	return &BuffCrawler{
		crawler: crawler,
		parser:  parser,
	}, nil
}

func (c *BuffCrawler) CrawlItemListingPage(itemName string, buffId, pageNum int, filters map[string]string) (*types.ListingsData, *types.CrawlerControl, error) {
	params := url.Values{}
	params.Add("game", BUFF_CSGO_NAME)
	params.Add("goods_id", strconv.Itoa(buffId))
	params.Add("page_num", strconv.Itoa(pageNum))
	params.Add("page_size", strconv.Itoa(BUFF_LISTING_ITEMS_PER_PAGE))
	params.Add("sort_by", "price.asc")
	params.Add("mode", "")
	params.Add("allow_tradable_cooldown", "1")
	params.Add("_", utils.GetBuffTimestamp())

	buffLog("Crawling page %d for %s\n", pageNum, itemName)

	utils.AddFilters(params, filters)

	savePath := path.Join(BUFF_RAW_RES_DIR, fmt.Sprintf("buff_l_%s_%d_%s.json", itemName, pageNum, shared.GetNowHHMMSS()))
	resData := &BuffListingResponseData{}
	resp, err := c.crawler.DoReqWithSave(BUFF_LISTING_API, params, "GET", savePath, resData, getRefererHeader(buffId))

	if err != nil {
		return nil, nil, err
	}

	// parse response
	data, err := c.parser.ParseItemListings(itemName, resp, resData)

	if err != nil {
		return nil, nil, err
	}

	buffLog("Parsed %d listings for %s\n", len(data.Listings), itemName)

	return data, c.parser.ParseListingControl(resData), nil
}

func (c *BuffCrawler) CrawlItemListings(itemName string, handler types.IHandler, config *types.CrawlTaskConfig) error {
	// reset stop flag
	c.crawler.ResetStop()

	var updatedItem *types.Item
	// round up
	numPages := utils.GetNumPages(config.MaxItems, BUFF_LISTING_ITEMS_PER_PAGE)
	buffLog("Crawling %d pages for %s\n", numPages, itemName)
	// convert name to buff id
	buffId, ok := shared.GetBuffIds()[itemName]
	if !ok {
		return errors.ErrItemNotFound
	}

	for i := 1; i <= numPages; i++ {
		data, control, err := c.CrawlItemListingPage(itemName, buffId, i, config.Filters)

		// handle stop
		if c.crawler.Stopped || err == errors.ErrCrawlerManuallyStopped {
			buffLog("Crawler manually stopped\n")
			handler.OnComplete(
				&types.ItemData{
					Item: updatedItem,
				})
			return nil
		}

		if err != nil {
			handler.OnError(err)
			return err
		}

		handler.OnResult(data)

		// merge item data
		if updatedItem == nil {
			updatedItem = data.Item
		} else {
			// update the price
			if shared.NumStrCmp(data.Item.SteamPrice.Price, updatedItem.SteamPrice.Price) < 0 {
				updatedItem.SteamPrice = data.Item.SteamPrice
			}
		}

		// handle control
		if control != nil {
			if i >= control.TotalPages {
				break
			}
		}

	}

	// only update the item after all pages are crawled
	handler.OnComplete(
		&types.ItemData{
			Item: updatedItem,
		},
	)

	return nil
}

func (c *BuffCrawler) CrawlItemTransactionPage(itemName string, buffId int, filters map[string]string) (*types.TransactionData, *types.CrawlerControl, error) {
	params := url.Values{}
	params.Add("game", BUFF_CSGO_NAME)
	params.Add("goods_id", strconv.Itoa(buffId))
	params.Add("_", utils.GetBuffTimestamp())

	utils.AddFilters(params, filters)

	savePath := path.Join(BUFF_RAW_RES_DIR, fmt.Sprintf("buff_t_%s_%s.json", itemName, shared.GetNowHHMMSS()))
	resData := &BuffTransactionResponseData{}
	resp, err := c.crawler.DoReqWithSave(BUFF_TRANSACTION_API, params, "GET", savePath, resData, getRefererHeader(buffId))

	if err != nil {
		return nil, nil, err
	}

	// parse response
	data, err := c.parser.ParseItemTransactions(itemName, resp, resData)

	if err != nil {
		return nil, nil, err
	}

	buffLog("Parsed %d transactions for %s\n", len(data.Transactions), itemName)

	return data, c.parser.ParseTransactionControl(resData), nil
}

func (c *BuffCrawler) CrawlItemTransactions(itemName string, handler types.IHandler, config *types.CrawlTaskConfig) error {
	// reset stop flag
	c.crawler.ResetStop()

	buffLog("Crawling transactions for %s\n", itemName)
	// convert name to buff id
	buffId, ok := shared.GetBuffIds()[itemName]
	if !ok {
		return errors.ErrItemNotFound
	}

	// only one page
	data, _, err := c.CrawlItemTransactionPage(itemName, buffId, config.Filters)

	if err != nil {
		handler.OnError(err)
		return err
	}

	handler.OnResult(data)

	handler.OnComplete(nil)

	return nil
}
