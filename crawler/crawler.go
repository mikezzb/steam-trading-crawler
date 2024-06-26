package crawler

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/mikezzb/steam-trading-crawler/errors"
	"github.com/mikezzb/steam-trading-crawler/utils"
	shared "github.com/mikezzb/steam-trading-shared"
)

// Base crawler providing basic functions of a crawler, use as embeded struct
type Crawler struct {
	client      *http.Client
	lastReqTime time.Time
	Stopped     bool
	config      *CrawlerConfig
}

func (c *Crawler) Stop() {
	c.Stopped = true
}

func (c *Crawler) ResetStop() {
	c.Stopped = false
}

type CrawlerConfig struct {
	Cookie      string
	AuthUrls    []string
	SleepMinSec int
	SleepMaxSec int
	Headers     map[string]string

	// for internal use
	sleepMin time.Duration
	sleepMax time.Duration
}

func NewCrawler(config *CrawlerConfig) (*Crawler, error) {
	// format configs
	config.sleepMin = time.Duration(config.SleepMinSec) * time.Second
	config.sleepMax = time.Duration(config.SleepMaxSec) * time.Second

	c := &Crawler{
		config: config,
	}
	client, err := utils.NewClientWithCookie(config.Cookie, config.AuthUrls)
	if err != nil {
		return nil, err
	}
	c.client = client
	// init last req time so the first req will do immediately
	c.lastReqTime = time.Now().Add(-c.config.sleepMax)

	return c, nil
}

func (c *Crawler) SetHeaders(req *http.Request, headers map[string]string) {
	for k, v := range headers {
		req.Header.Set(k, v)
	}
}

func (c *Crawler) sleepForSafe() {
	timeSinceLastReq := time.Since(c.lastReqTime)

	if timeSinceLastReq < c.config.sleepMin {
		sleepDuration := shared.GetRandomSleepDuration(
			c.config.SleepMinSec, c.config.SleepMaxSec)
		sleepTime := sleepDuration - timeSinceLastReq
		log.Printf("Sleeping for %s\n", sleepTime)
		time.Sleep(sleepTime)
	}

	c.lastReqTime = time.Now()
}

func (c *Crawler) DoReq(u string, params url.Values, method string, headers map[string]string) (*http.Response, error) {
	c.sleepForSafe()

	if c.Stopped {
		return nil, errors.ErrCrawlerManuallyStopped
	}

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
	c.SetHeaders(req, c.config.Headers)
	c.SetHeaders(req, headers)

	// do request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Crawler) DoReqWithSave(u string, params url.Values, method, savePath string, resData interface{}, headers map[string]string) (*http.Response, error) {
	resp, err := c.DoReq(u, params, method, headers)
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

func (c *Crawler) GetCookies() (string, error) {
	if c.client.Jar == nil || len(c.config.AuthUrls) == 0 {
		return "", errors.ErrNoCookies
	}
	authUrl := c.config.AuthUrls[0]
	parsedUrl, _ := url.Parse(authUrl)
	cookies := c.client.Jar.Cookies(parsedUrl)
	return utils.StringifyCookies(cookies), nil
}
