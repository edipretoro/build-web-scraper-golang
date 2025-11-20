package main

import (
	"encoding/csv"
	"os"
	"strings"
)

func writeCSVReport(pages map[string]PageData, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return nil
	}
	w := csv.NewWriter(file)

	firstLine := []string{"page_url","h1","first_paragraph","outgoing_link_urls","image_urls"}
	if err = w.Write(firstLine); err != nil {
		return err
	}
	for _, page := range pages {
		record := []string{
			page.URL,
			page.H1,
			page.FirstParagraph,
			strings.Join(page.OutgoingLinks, ";"),
			strings.Join(page.ImageURLs, ";"),
		}
		if err := w.Write(record); err != nil {
			return err
		}
	}

	// Write any buffered data to the underlying writer (standard output).
	w.Flush()

	if err := w.Error(); err != nil {
		return err
	}
	return nil
}
