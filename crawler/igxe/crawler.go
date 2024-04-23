package igxe

import (
	"net/url"
	"strconv"

	"github.com/mikezzb/steam-trading-crawler/crawler"
	"github.com/mikezzb/steam-trading-crawler/errors"
	"github.com/mikezzb/steam-trading-crawler/types"
	"github.com/mikezzb/steam-trading-crawler/utils"
	shared "github.com/mikezzb/steam-trading-shared"
	"github.com/mikezzb/steam-trading-shared/database/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IgxeCrawler struct {
	crawler *crawler.Crawler
	parser  *IgxeParser
}

func NewCrawler(cookie string) (*IgxeCrawler, error) {
	crawler, err := crawler.NewCrawler(&crawler.CrawlerConfig{
		Cookie:      cookie,
		AuthUrls:    nil,
		SleepMinSec: IGXE_SLEEP_TIME_MIN_S,
		SleepMaxSec: IGXE_SLEEP_TIME_MAX_S,
		Headers:     IGXE_HEADERS,
	})

	if err != nil {
		return nil, err
	}

	parser := &IgxeParser{}

	return &IgxeCrawler{
		crawler: crawler,
		parser:  parser,
	}, nil
}

func (c *IgxeCrawler) Stop() {
	c.crawler.Stop()
}

func (c *IgxeCrawler) GetCookies() (string, error) {
	return c.crawler.GetCookies()
}

func (c *IgxeCrawler) getItemWithPrice(name string, price primitive.Decimal128) *model.Item {
	return &model.Item{
		Name: name,
		IgxePrice: model.MarketPrice{
			Price:     price,
			UpdatedAt: shared.GetNow(),
		},
	}
}

func (c *IgxeCrawler) crawlItemListingPage(itemName string, igxeId, pageNum int, filters map[string]string) (*types.ListingsData, *types.CrawlerControl, error) {
	igxeLog("Crawling page %d for %s\n", pageNum, itemName)

	params := url.Values{}
	params.Add("page_no", strconv.Itoa(pageNum))
	params.Add("product_id", strconv.Itoa(igxeId))
	// default sort by price
	params.Add("sort", "0")
	params.Add("sort_rule", "0")

	utils.AddFilters(params, filters)

	productUrl := getIgxeListingUrl(igxeId)
	savePath := getIgxeSavePath(itemName, pageNum, "l")

	resData := &IgxeListingResponseData{}
	resp, err := c.crawler.DoReqWithSave(productUrl, params, "GET", savePath, resData, getRefererHeader(igxeId))

	if err != nil {
		return nil, nil, err
	}

	// parse data
	data, err := c.parser.ParseItemListings(itemName, resp, resData)

	igxeLog("Parsed %d listings for %s\n", len(data.Listings), itemName)

	if err != nil {
		return nil, nil, err
	}

	return data, c.parser.ParseListingControl(resData), nil
}

func (c *IgxeCrawler) CrawlItemListings(itemName string, handler types.IHandler, config *types.CrawlTaskConfig) error {
	// reset stop flag
	c.crawler.ResetStop()

	igxeId, ok := shared.GetIgxeIds()[itemName]
	if !ok {
		err := errors.ErrItemNotFound
		handler.OnError(err)
		return err
	}

	var updatedItem *model.Item = c.getItemWithPrice(itemName, shared.MAX_DECIMAL128)
	numPages := utils.GetNumPages(config.MaxItems, IGXE_LISTING_ITEMS_PER_PAGE)

	for i := 1; i <= numPages; i++ {
		data, control, err := c.crawlItemListingPage(itemName, igxeId, i, config.Filters)

		// handle stop
		if c.crawler.Stopped || err == errors.ErrCrawlerManuallyStopped {
			handler.OnComplete(&types.ItemData{
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
			if shared.DecCompareTo(data.Item.IgxePrice.Price, updatedItem.IgxePrice.Price) < 0 {
				updatedItem.IgxePrice.Price = data.Item.IgxePrice.Price
				updatedItem.IgxePrice.UpdatedAt = data.Item.IgxePrice.UpdatedAt
			}
		}

		// handle control
		if control != nil {
			if i >= control.TotalPages {
				break
			}
		}
	}

	handler.OnComplete(&types.ItemData{
		Item: updatedItem,
	})

	return nil
}

func (c *IgxeCrawler) crawlItemTransactionPage(itemName string, igxeId, pageNum int, filters map[string]string) (*types.TransactionData, *types.CrawlerControl, error) {
	params := url.Values{}
	utils.AddFilters(params, filters)

	transUrl := getIgxeTransactionUrl(igxeId)
	savePath := getIgxeSavePath(itemName, pageNum, "t")

	resData := &IgxeTransactionResponseData{}
	resp, err := c.crawler.DoReqWithSave(transUrl, params, "GET", savePath, resData, getRefererHeader(igxeId))

	if err != nil {
		return nil, nil, err
	}

	// parse data
	data, err := c.parser.ParseItemTransactions(itemName, resp, resData)

	igxeLog("Parsed %d transactions for %s\n", len(data.Transactions), itemName)

	if err != nil {
		return nil, nil, err
	}

	return data, c.parser.ParseTransactionControl(resData), nil
}

func (c *IgxeCrawler) CrawlItemTransactions(itemName string, handler types.IHandler, config *types.CrawlTaskConfig) error {
	// reset stop flag
	c.crawler.ResetStop()

	igxeId, ok := shared.GetIgxeIds()[itemName]
	if !ok {
		return errors.ErrItemNotFound
	}

	// only one page
	data, _, err := c.crawlItemTransactionPage(itemName, igxeId, 1, config.Filters)

	if err != nil {
		handler.OnError(err)
		return err
	}

	handler.OnResult(data)

	handler.OnComplete(nil)

	return nil
}
