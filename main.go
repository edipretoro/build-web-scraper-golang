package main

import (
	"fmt"
	"log"
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
	content, err := getHTML(args[0])
	if err != nil {
		log.Fatalf("problem when fetching the content: %s", err.Error())
	}
	fmt.Print(content)
}
