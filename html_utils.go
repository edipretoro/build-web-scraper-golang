package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

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
	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "text/html") {
		log.Printf("Bad Content-Type: %s", resp.Header.Get("Content-Type"))
		return "", errors.New("content fetched is not text/html")
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	
	return string(content), nil
}

type config struct {
	pages              map[string]PageData
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	_, ok := cfg.pages[normalizedURL]
	if ok {
		return true
	}
	return false
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	parsedURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		log.Printf("problem when parsing the URL: %s", err.Error())
		return
	}
	if parsedURL.Hostname() != cfg.baseURL.Hostname() {
		log.Printf("- %s: %s", cfg.baseURL.String(), rawCurrentURL)
		return
	}
	normURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		log.Printf("problem when normalizing `%s`: %s", rawCurrentURL, err.Error())
		return
	}
	if cfg.addPageVisit(normURL) {
		return
	}
	log.Printf("getting content from `%s`", rawCurrentURL)
	content, err := getHTML(rawCurrentURL)
	if err != nil {
		log.Printf("problem when fetching content from `%s`: %s", rawCurrentURL, err.Error())
		return
	}
	cfg.mu.Lock()
	cfg.pages[normURL] = extractPageData(content, rawCurrentURL)
	cfg.mu.Unlock()
	
	for _, url := range cfg.pages[normURL].OutgoingLinks {
		cfg.wg.Add(1)
		go func() {
			defer func() {
				cfg.wg.Done()
				<-cfg.concurrencyControl
			}()
			cfg.concurrencyControl <- struct{}{}
			cfg.crawlPage(url)
		}()
	}
}
