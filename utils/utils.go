package utils

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"steam-trading/shared"
)

func PostFormatListing(name string, listings []shared.Listing) []shared.Listing {
	for i := range listings {
		// add name to listings
		listings[i].Name = name
		listings[i].Rarity = shared.GetListingTier(listings[i])
	}
	return listings
}

// decodes a gzip-encoded reader and returns a decoded reader.
func DecodeReader(reader io.Reader) (io.ReadCloser, error) {
	// Attempt to create a gzip reader
	gzReader, err := gzip.NewReader(reader)
	if err != nil {
		// If the input is not gzip-encoded, return the original reader
		if err == gzip.ErrHeader {
			return io.NopCloser(reader), nil
		}
		return nil, err
	}
	return gzReader, nil
}

func SaveRawResponse(resp *http.Response, path string) error {
	// Open the file with write access, create if not exists, truncate otherwise
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Decode the response body if gzip-encoded
	decodedReader, err := DecodeReader(resp.Body)
	if err != nil {
		return err
	}

	defer decodedReader.Close()

	// Copy the decoded response body to the file
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
