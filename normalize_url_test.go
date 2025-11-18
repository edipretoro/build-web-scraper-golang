package main

import (
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		expected string
	}{
		{
			name:     "normalize_url_protocol",
			inputURL: "https://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "normalize_url_slash",
			inputURL: "https://blog.boot.dev/path/",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "normalize_url_capitals",
			inputURL: "https://BLOG.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "normalize_url_http",
			inputURL: "http://BLOG.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if err != nil {
				t.Errorf("FAIL: unexpected error: %v", err)
				return
			}
			if actual != tc.expected {
				t.Errorf("FAIL: expected URL: %v, actual: %v", tc.expected, actual)
			}
		})
	}
}
