package utils

import (
	"net/http"
	"net/url"
	"steam-trading/shared"
	"strings"

	"github.com/mikezzb/steam-trading-crawler/types"
)

func ParseCookieString(cookieStr string) []*http.Cookie {
	var cookies []*http.Cookie

	cookieParts := strings.Split(cookieStr, ";")

	for _, part := range cookieParts {
		part = strings.TrimSpace(part)

		keyValue := strings.SplitN(part, "=", 2)

		if len(keyValue) == 2 {
			cookie := &http.Cookie{
				Name:  keyValue[0],
				Value: keyValue[1],
			}

			cookies = append(cookies, cookie)
		}
	}

	return cookies
}

func NewClientWithCookie(cookieStr string, apiURLs []string) (*http.Client, error) {
	jar := NewBrowserCookieJar()
	cookies := ParseCookieString(cookieStr)
	for _, apiURL := range apiURLs {
		parsedURL, err := url.Parse(apiURL)
		if err != nil {
			return nil, err
		}
		jar.SetCookies(parsedURL, cookies)
	}

	client := &http.Client{
		Jar: jar,
	}

	return client, nil
}

func StringifyCookies(cookies []*http.Cookie) string {
	cookieStr := ""
	for _, cookie := range cookies {
		cookieStr += cookie.String() + "; "
	}
	cookieStr = strings.TrimRight(cookieStr, "; ")
	return cookieStr
}

func UpdateSecrets(crawler types.Crawler, store shared.PersisitedStore, label string) {
	cookieStr, err := crawler.GetCookies()
	if err != nil {
		return
	}
	store.Set(label, cookieStr)
	store.Save()
}
