package buff

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"steam-trading/crawler/errors"
	"steam-trading/shared"
	"strconv"
	"time"

	"github.com/mikezzb/steam-trading-crawler/utils"
)

type BuffCrawler struct {
	secret string
	client *http.Client
	parser *BuffParser
}

func (c *BuffCrawler) Init(secret string) error {
	c.secret = secret
	c.client = &http.Client{}
	c.parser = &BuffParser{}

	return nil
}

func (c *BuffCrawler) SetHeaders(req *http.Request) {
	for k, v := range BUFF_HEADERS {
		req.Header.Set(k, v)
	}
	req.Header.Set("Cookie", c.secret)

	fmt.Println(req.URL.String())
	fmt.Println(req.Method)
	fmt.Println(req.Header)
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

	// return nil, errors.ErrItemNotFound

	// do request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *BuffCrawler) CrawItemListingPage(itemName string, pageNum int, filters map[string]string) error {
	// convert name to buff id
	buffId, ok := shared.GetBuffIds()[itemName]
	if !ok {
		return errors.ErrItemNotFound
	}

	params := url.Values{}
	params.Add("game", BUFF_CSGO_NAME)
	params.Add("goods_id", strconv.Itoa(buffId))
	params.Add("page_num", strconv.Itoa(pageNum))
	params.Add("sort_by", "price.asc")
	params.Add("mode", "")
	params.Add("allow_tradable_cooldown", "1")
	// timestamp in BJT (ms)
	// params.Add("_", shared.GetTimestampNow())

	fmt.Printf("Crawling page %d for %s\n", pageNum, itemName)

	for k, v := range filters {
		params.Add(k, v)
	}

	resp, err := c.DoReq(BUFF_LISTING_API, params, "GET")

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// save raw response
	saveName := fmt.Sprintf("buff_l_%s_%d_%s.json", itemName, pageNum, shared.GetTimestampNow())
	_ = utils.SaveRawResponse(resp, path.Join(BUFF_RAW_RES_DIR, saveName))

	// parse response
	item, listings, err := c.parser.ParseItemListings(itemName, resp)

	if err != nil {
		return err
	}

	// save to DB

	fmt.Println(item)
	fmt.Println(listings)

	return nil
}

func (c *BuffCrawler) CrawlItemListings(itemName string, maxListings int, filters map[string]string) error {
	// maxPages := maxListings / BUFF_LISTING_ITEMS_PER_PAGE
	maxPages := 1

	for i := 1; i <= maxPages; i++ {
		err := c.CrawItemListingPage(itemName, i, filters)
		if err != nil {
			return err
		}

		if i != maxPages {
			// sleep for buff sleep time seconds
			time.Sleep(time.Duration(BUFF_SLEEP_TIME) * time.Second)
		}
	}

	return nil
}

func (c *BuffCrawler) CrawlItemTransactions(itemName string) error {
	return nil
}
