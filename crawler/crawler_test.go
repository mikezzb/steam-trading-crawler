package crawler_test

import (
	"net/url"
	"testing"

	"github.com/mikezzb/steam-trading-crawler/crawler"
)

func TestThrottle(t *testing.T) {
	t.Run("Sleep", func(t *testing.T) {
		crawler, _ := crawler.NewCrawler(&crawler.CrawlerConfig{
			Cookie:      "cookie",
			SleepMinSec: 1,
			SleepMaxSec: 2,
		})
		// no throttle
		crawler.DoReq("localhost:8000", url.Values{}, "GET", nil)
		// throttled
		crawler.DoReq("localhost:8000", url.Values{}, "GET", nil)
	})
}
