package main

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getH1FromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}

	h1 := doc.Find("h1").First()
	return strings.TrimSpace(h1.Text())
}

func getFirstParagraphFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}

	main := doc.Find("main").First()
	var p *goquery.Selection
	if main.Length() > 0 {
		p = main.Find("p").First()
	} else {
		p = doc.Find("p").First()
	}
	return strings.TrimSpace(p.Text())
}

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	var urls []string
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			parsedURL, err := url.Parse(href)
			if err == nil {
				resolvedURL := baseURL.ResolveReference(parsedURL)
				urls = append(urls, resolvedURL.String())
			}
		}
	})

	return urls, nil
}

func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	var images []string
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists {
			parsedURL, err := url.Parse(src)
			if err == nil {
				resolvedURL := baseURL.ResolveReference(parsedURL)
				images = append(images, resolvedURL.String())
			}
		}
	})

	return images, nil
}

type PageData struct {
	URL            string
	H1             string
	FirstParagraph string
	OutgoingLinks  []string
	ImageURLs      []string
}

func extractPageData(html, pageURL string) PageData {
	parsedURL, err := url.Parse(strings.ToLower(pageURL))
	if err != nil {
		return PageData{}
	}
	outgoingLinks, err := getURLsFromHTML(html, parsedURL)
	if err != nil {
		outgoingLinks = []string{}
	}
	imageURLs, err := getImagesFromHTML(html, parsedURL)
	if err != nil {
		imageURLs = []string{}
	}
	return PageData{
		URL: pageURL,
		H1: getH1FromHTML(html),
		FirstParagraph: getFirstParagraphFromHTML(html),
		OutgoingLinks: outgoingLinks,
		ImageURLs: imageURLs,
	}
}

func getHTML(rawURL string) (string, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("User-Agent", "BootCrawler/1.0")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", errors.New("problem when fetching the URL")
	}
	if resp.Header.Get("Content-Type") != "text/html" {
		return "", errors.New("content fetched is not text/html")
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	
	return string(content), nil
}
