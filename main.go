package main

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"sync"
)

func init() {
    log.SetOutput(io.Discard)
}

func main() {
	maxConcurrency := 10
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	}
	if len(args) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}
	fmt.Println("starting crawl of:", args[0])
	parsedURL, err := url.Parse(args[0])
	if err != nil {
		log.Printf("problem when parsing the URL: %s", err.Error())
		os.Exit(1)
	}
	cfg := config{
		pages: make(map[string]PageData),
		baseURL: parsedURL,
		mu: &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg: &sync.WaitGroup{},
	}
	cfg.crawlPage(args[0])
	cfg.wg.Wait()
	fmt.Println("printing result:")
	for k, _ := range cfg.pages {
		fmt.Printf("- %s\n", k)
	}
}
