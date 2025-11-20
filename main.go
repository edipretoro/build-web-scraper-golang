package main

import (
	"fmt"
	"os"
)

func main() {
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
	pages := make(map[string]int)
	crawlPage(args[0], args[0], pages)
	fmt.Println("printing result:")
	for k, v := range pages {
		fmt.Printf("- %s: %d\n", k, v)
	}
}
