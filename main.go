package main

import (
	"fmt"
	// "io"
	"log"
	"net/url"
	"os"
	"strconv"
	"sync"
)

// func init() {
//     log.SetOutput(io.Discard)
// }

func main() {
	args := os.Args[1:]
	if len(args) < 3 {
		fmt.Println("not enough arguments: website maxConcurrency maxPages")
		os.Exit(1)
	}
	if len(args) > 3 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}
	baseURL := args[0]
	maxConcurrency, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("problem when parsing maxConcurrency")
		os.Exit(1)
	}
	maxPages, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("problem when parsing maxPages")
		os.Exit(1)
	}
	fmt.Println("starting crawl of:", baseURL)
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		log.Printf("problem when parsing the URL: %s", err.Error())
		os.Exit(1)
	}
	cfg := config{
		pages: make(map[string]PageData),
		maxPages: maxPages,
		baseURL: parsedURL,
		mu: &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg: &sync.WaitGroup{},
	}
	cfg.wg.Add(1)
	go cfg.crawlPage(baseURL)
	cfg.wg.Wait()
	fmt.Println("Saving crawl into `report.csv`")
	writeCSVReport(cfg.pages, "report.csv")
}

