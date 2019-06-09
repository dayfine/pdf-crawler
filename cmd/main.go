package main

import (
	"flag"

	"github.com/dayfine/pdf-crawler/crawler"
)

func main() {
	inputUrl := flag.String("url", "", "Specify which URL to crawl (e.g. google.com)")
	depth := flag.Int64("depth", 5, "Specify the level of depth to crawl through")
	followForeignHosts := flag.Bool("foreign", false, "Whether to visit hosts other than the one for initally specified URL")
	saveDir := flag.String("dir", "~/Downloads", "Specify where to save the downloaded files")
	slient := flag.Bool("slient", true, "If false, print a list of URL processed at the end")

	flag.Parse()

	pdfCrawler := crawler.NewCrawler(crawler.CrawlerOptions{
		InitialUrl:         *inputUrl,
		FollowForeignHosts: *followForeignHosts,
		SaveDir:            *saveDir,
	})

	pdfCrawler.Crawl(crawler.CrawlParameters{
		InputUrl: *inputUrl,
		Depth:    *depth,
	})
	if !*slient {
		pdfCrawler.PrintAllProcessed()
	}
}
