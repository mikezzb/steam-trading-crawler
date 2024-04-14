package buff

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/mikezzb/steam-trading-crawler/errors"
	"github.com/mikezzb/steam-trading-crawler/types"
	"github.com/mikezzb/steam-trading-crawler/utils"
	shared "github.com/mikezzb/steam-trading-shared"
)

type BuffCrawler struct {
	cookie      string
	client      *http.Client
	parser      *BuffParser
	lastReqTime time.Time
}

func (c *BuffCrawler) Init(cookie string) error {
	c.cookie = cookie
	client, err := utils.NewClientWithCookie(cookie, []string{BUFF_LISTING_API, BUFF_TRANSACTION_API})
	if err != nil {
		return err
	}
	c.client = client
	c.parser = &BuffParser{}

	return nil
}

func (c *BuffCrawler) SetHeaders(req *http.Request) {
	for k, v := range BUFF_HEADERS {
		req.Header.Set(k, v)
	}
}

func (c *BuffCrawler) SleepForSafe() {
	timeSinceLastReq := time.Since(c.lastReqTime)

	if timeSinceLastReq < BUFF_SLEEP_TIME_MIN {
		sleepDuration := shared.GetRandomSleepDuration(
			BUFF_SLEEP_TIME_MIN_S, BUFF_SLEEP_TIME_MAX_S)
		sleepTime := sleepDuration - timeSinceLastReq
		log.Printf("Sleeping for %s\n", sleepTime)
		time.Sleep(sleepTime)
	}

	c.lastReqTime = time.Now()
}

func (c *BuffCrawler) DoReq(u string, params url.Values, method string) (*http.Response, error) {
	c.SleepForSafe()

	// encode params
	baseUrl, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	baseUrl.RawQuery = params.Encode()

	// make request
	req, err := http.NewRequest(method, baseUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	// set headers
	c.SetHeaders(req)

	return nil, nil

	// do request
	// resp, err := c.client.Do(req)
	// if err != nil {
	// 	return nil, err
	// }

	// return resp, nil
}

func (c *BuffCrawler) DoReqWithSave(u string, params url.Values, method, savePath string, resData interface{}) (*http.Response, error) {
	resp, err := c.DoReq(u, params, method)
	if err != nil {
		return nil, err
	}

	// save raw response
	bodyBytes, _ := utils.Body2Bytes(resp)

	defer resp.Body.Close()

	err = utils.SaveResponseBody(bodyBytes, savePath)

	if err != nil {
		return nil, err
	}

	// decode response
	decodedReader, err := utils.ReadBytes(bodyBytes)
	if err != nil {
		return nil, err
	}
	defer decodedReader.Close()

	// unmarshal response
	if err := json.NewDecoder(decodedReader).Decode(&resData); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *BuffCrawler) GetCookies() (string, error) {
	parsedUrl, _ := url.Parse(BUFF_LISTING_API)
	cookies := c.client.Jar.Cookies(parsedUrl)
	shared.PrintCookies(cookies, "Saved Cookies")
	return utils.StringifyCookies(cookies), nil
}

func (c *BuffCrawler) CrawlItemListingPage(itemName string, buffId, pageNum int, filters map[string]string) (*types.ListingsData, error) {
	params := url.Values{}
	params.Add("game", BUFF_CSGO_NAME)
	params.Add("goods_id", strconv.Itoa(buffId))
	params.Add("page_num", strconv.Itoa(pageNum))
	params.Add("sort_by", "price.asc")
	params.Add("mode", "")
	params.Add("allow_tradable_cooldown", "1")
	params.Add("_", shared.GetTimestampNow())

	log.Printf("Crawling page %d for %s\n", pageNum, itemName)

	for k, v := range filters {
		params.Add(k, v)
	}

	savePath := path.Join(BUFF_RAW_RES_DIR, fmt.Sprintf("buff_l_%s_%d_%s.json", itemName, pageNum, shared.GetTimestampNow()))
	resData := &BuffListingResponseData{}
	resp, err := c.DoReqWithSave(BUFF_LISTING_API, params, "GET", savePath, resData)

	if err != nil {
		return nil, err
	}

	// parse response
	data, err := c.parser.ParseItemListings(itemName, resp, resData)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *BuffCrawler) CrawlItemListings(itemName string, handler *types.Handler, config *types.CrawlerConfig) error {
	maxPages := config.MaxItems / BUFF_LISTING_ITEMS_PER_PAGE
	// maxPages := 1
	log.Printf("Crawling %d pages for %s\n", maxPages, itemName)
	// convert name to buff id
	buffId, ok := shared.GetBuffIds()[itemName]
	if !ok {
		return errors.ErrItemNotFound
	}

	for i := 1; i <= maxPages; i++ {
		data, err := c.CrawlItemListingPage(itemName, buffId, i, config.Filters)
		handler.OnResult(data)

		if err != nil {
			if handler.OnError != nil {
				handler.OnError(err)
			}
			return err
		}
	}

	if handler.OnComplete != nil {
		handler.OnComplete()
	}

	return nil
}

func (c *BuffCrawler) CrawlItemTransactionPage(itemName string, buffId int, filters map[string]string) (*types.TransactionData, error) {
	params := url.Values{}
	params.Add("game", BUFF_CSGO_NAME)
	params.Add("goods_id", strconv.Itoa(buffId))
	params.Add("_", shared.GetTimestampNow())

	for k, v := range filters {
		params.Add(k, v)
	}

	savePath := path.Join(BUFF_RAW_RES_DIR, fmt.Sprintf("buff_t_%s_%s.json", itemName, shared.GetTimestampNow()))
	resData := &BuffTransactionResponseData{}
	resp, err := c.DoReqWithSave(BUFF_TRANSACTION_API, params, "GET", savePath, resData)

	if err != nil {
		return nil, err
	}

	// parse response
	data, err := c.parser.ParseItemTransactions(itemName, resp, resData)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *BuffCrawler) CrawlItemTransactions(itemName string, handler *types.Handler, config *types.CrawlerConfig) error {
	log.Printf("Crawling transactions for %s\n", itemName)
	// convert name to buff id
	buffId, ok := shared.GetBuffIds()[itemName]
	if !ok {
		return errors.ErrItemNotFound
	}

	// only one page
	data, err := c.CrawlItemTransactionPage(itemName, buffId, config.Filters)

	if err != nil {
		if handler.OnError != nil {
			handler.OnError(err)
		}
		return err
	}

	handler.OnResult(data)

	if handler.OnComplete != nil {
		handler.OnComplete()
	}

	return nil
}
