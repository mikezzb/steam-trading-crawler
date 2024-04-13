package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/mikezzb/steam-trading-crawler/types"
	shared "github.com/mikezzb/steam-trading-shared"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/mikezzb/steam-trading-shared/database/model"
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

func PostFormatListing(name string, listings []model.Listing) []model.Listing {
	for i := range listings {
		// add name to listings
		listings[i].Name = name
		listings[i].Rarity = shared.GetListingTier(listings[i])
	}
	return listings
}

// decodes a gzip-encoded reader and returns a decoded reader.
// DecodeReader decodes a gzip-encoded reader and returns a decoded reader.

// DecodeReader reads from the given reader and decodes JSON or gzipped JSON.
func ReadBytes(b []byte) (io.ReadCloser, error) {
	// Check if the byte slice starts with the gzip magic numbers.
	if len(b) >= 2 && b[0] == 0x1f && b[1] == 0x8b {
		// The content is gzipped. Create a gzip reader.
		gzReader, err := gzip.NewReader(bytes.NewReader(b))
		if err != nil {
			return nil, err
		}
		return gzReader, nil
	}

	// If it's not gzipped, return a reader for the plain byte slice.
	// We wrap it in a NopCloser to match the io.ReadCloser return type.
	return io.NopCloser(bytes.NewReader(b)), nil
}

func Body2Bytes(resp *http.Response) ([]byte, error) {
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	return bytes, nil
}

func SaveResponseBody(b []byte, path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	decodedReader, err := ReadBytes(b)
	if err != nil {
		return err
	}

	defer decodedReader.Close()

	_, err = io.Copy(file, decodedReader)
	if err != nil {
		return err
	}

	return nil
}

// WriteJSONToFile writes data to a JSON file.
func WriteJSONToFile(data interface{}, filePath string) error {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return err
	}

	return nil
}

// read image from url and save it to the specified path
func DownloadImage(imageURL, path string) error {
	resp, err := http.Get(imageURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
