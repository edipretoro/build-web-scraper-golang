package main

import (
	"net/url"
	"strings"
)

func normalizeURL(page_url string) (string, error) {
	parsedURL, err := url.Parse(strings.ToLower(page_url))
	if err != nil {
		return "", err
	}

	normalizedURL := parsedURL.Host + parsedURL.Path
	if len(normalizedURL) > 0 && normalizedURL[len(normalizedURL)-1] == '/' {
		normalizedURL = normalizedURL[:len(normalizedURL)-1]
	}

	return normalizedURL, nil
}
