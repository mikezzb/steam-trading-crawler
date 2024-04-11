package utils

import (
	"net/http"
	"net/url"
)

type BrowserCookieJar struct {
	cookies map[string][]*http.Cookie
}

// browser cookie jar simulates the browser's set-cookie behavior implementing the golang http.CookieJar interface

func (j *BrowserCookieJar) Cookies(u *url.URL) []*http.Cookie {
	return j.cookies[u.Host]
}

// SetCookies implements the http.CookieJar interface, and it merges cookies with existing cookies
func (j *BrowserCookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	oldCookies := j.Cookies(u)
	if oldCookies == nil {
		j.cookies[u.Host] = cookies
		return
	}

	// merge cookies with existing cookies
	for _, cookie := range cookies {
		// search for existing cookie with the same name
		found := false
		for _, existingCookie := range oldCookies {
			if existingCookie.Name == cookie.Name {
				existingCookie.Value = cookie.Value
				found = true
				break
			}
		}
		if !found {
			j.cookies[u.Host] = append(j.cookies[u.Host], cookie)
		}
	}

}

func NewBrowserCookieJar() *BrowserCookieJar {
	return &BrowserCookieJar{cookies: make(map[string][]*http.Cookie)}
}
