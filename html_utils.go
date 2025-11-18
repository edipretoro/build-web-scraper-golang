package main

import (
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
