package main

import (
	"flag"

	"github.com/dayfine/pdf-crawler/crawler"
)

func main() {
	initialUrl := flag.String("url", "", "Specify which URL to crawl (e.g. google.com)")
	saveDir := flag.String("dir", "~/Downloads", "Specify where to save the downloaded files")
	flag.Parse()

	crawler.Crawl(*initialUrl, *saveDir)
}
