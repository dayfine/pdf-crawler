package crawler

import (
	"fmt"
	"log"
)

var fetcher UrlFetcher

func Crawl(initialUrl string, saveDir string) {
	fetchResult, err := fetcher.Fetch(initialUrl)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Fetched: %+v\n", fetchResult)
}

func init() {
	fetcher = NewUrlFetcher()
}
