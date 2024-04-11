package buff

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"steam-trading/crawler/errors"
	"steam-trading/shared"
	"strconv"

	"github.com/mikezzb/steam-trading-crawler/types"
	"github.com/mikezzb/steam-trading-crawler/utils"
)

type BuffCrawler struct {
	cookie string
	client *http.Client
	parser *BuffParser
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

func (c *BuffCrawler) DoReq(u string, params url.Values, method string) (*http.Response, error) {
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

	// do request
	resp, err := c.client.Do(req)
	if err != nil {
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

func (c *BuffCrawler) CrawItemListingPage(itemName string, pageNum int, filters map[string]string) (*types.ListingsData, error) {
	// convert name to buff id
	buffId, ok := shared.GetBuffIds()[itemName]
	if !ok {
		return nil, errors.ErrItemNotFound
	}

	params := url.Values{}
	params.Add("game", BUFF_CSGO_NAME)
	params.Add("goods_id", strconv.Itoa(buffId))
	params.Add("page_num", strconv.Itoa(pageNum))
	params.Add("sort_by", "price.asc")
	params.Add("mode", "")
	params.Add("allow_tradable_cooldown", "1")
	params.Add("_", shared.GetTimestampNow())

	fmt.Printf("Crawling page %d for %s\n", pageNum, itemName)

	for k, v := range filters {
		params.Add(k, v)
	}

	resp, err := c.DoReq(BUFF_LISTING_API, params, "GET")

	if err != nil {
		return nil, err
	}

	// save raw response
	bodyBytes, _ := utils.Body2Bytes(resp)

	saveName := fmt.Sprintf("buff_l_%s_%d_%s.json", itemName, pageNum, shared.GetTimestampNow())
	_ = utils.SaveResponseBody(bodyBytes, path.Join(BUFF_RAW_RES_DIR, saveName))

	// parse response
	data, err := c.parser.ParseItemListings(itemName, bodyBytes)

	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *BuffCrawler) CrawlItemListings(itemName string, config types.CrawlerConfig) error {
	maxPages := config.MaxItems / BUFF_LISTING_ITEMS_PER_PAGE
	// maxPages := 1
	fmt.Printf("Crawling %d pages for %s\n", maxPages, itemName)

	for i := 1; i <= maxPages; i++ {
		data, err := c.CrawItemListingPage(itemName, i, config.Filters)
		config.OnResult(data)

		if err != nil {
			if config.OnError != nil {
				config.OnError(err)
			}
			return err
		}

		if i != maxPages {
			shared.RandomSleep(BUFF_SLEEP_TIME_MIN, BUFF_SLEEP_TIME_MAX)
		}
	}

	if config.OnComplete != nil {
		config.OnComplete()
	}

	return nil
}

func (c *BuffCrawler) CrawlItemTransactions(itemName string, config types.CrawlerConfig) error {
	return nil
}
