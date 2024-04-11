package utils

import (
	"bytes"
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
