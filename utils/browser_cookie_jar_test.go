package utils_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/mikezzb/steam-trading-crawler/utils"
	shared "github.com/mikezzb/steam-trading-shared"
)

// const tragetUrl = "https://buff.163.com/api/market/goods/sell_order?game=csgo&goods_id=43018&page_num=1&sort_by=default&mode=&allow_tradable_cooldown=1&_=1712867980716"

// var parsedUrl, _ = url.Parse(tragetUrl)

func TestBrowserCookieJar(t *testing.T) {
	t.Run("Cookie", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// print the cookies
			cookies := r.Cookies()
			log.Printf("Recevied request")
			shared.PrintCookies(cookies, "Request")
			// mock set-cookie response
			w.Header().Set("Set-Cookie", "session=Updated; Path=/")
			w.Header().Set("Set-Cookie", "csrf_token=Updated; Path=/")
		}))

		defer server.Close()

		parsedUrl, _ := url.Parse(server.URL)

		log.Printf("Parsed URL: %v\n", parsedUrl.Host)

		// init default cookie
		keyStore, _ := shared.NewPersisitedStore("../secrets.json")

		client, _ := utils.NewClientWithCookie(keyStore.Get("buff_secret").(string), []string{parsedUrl.String()})
		oldCookies := client.Jar.Cookies(parsedUrl)
		shared.PrintCookies(oldCookies, "Old")

		// set cookie from resp 1
		resp, err := client.Get(server.URL)
		resp.Request.URL = parsedUrl
		if err != nil {
			t.Errorf("Failed to get response: %v", err)
		}

		shared.PrintCookies(client.Jar.Cookies(parsedUrl), "Updated")

		// make resp 2
		_, err = client.Get(server.URL)

		shared.PrintCookies(client.Jar.Cookies(parsedUrl), "Save")

	})
}
