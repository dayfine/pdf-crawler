package crawler

import (
	"fmt"
	"log"
	"path/filepath"
)

var ignoredFileExts = []string{".mp3", ".wav", ".mkv", ".flv", ".vob", ".ogv", ".ogg", ".gif", ".avi", ".mov", ".wmv", ".mp4", ".mp3", ".mpg"}

type CrawlParameters struct {
	InputUrl string
	Depth    int64
}

type CrawlerOptions struct {
	InitialUrl         string
	FollowForeignHosts bool
	SaveDir            string
}

type Crawler interface {
	Crawl(param CrawlParameters)
	PrintAllProcessed()
}

func NewCrawler(options CrawlerOptions) Crawler {
	processedMap := make(map[string]bool)
	ignoredFileExtMap := make(map[string]bool)
	for _, ext := range ignoredFileExts {
		ignoredFileExtMap[ext] = true
	}
	return &CrawlerImpl{
		options,
		NewUrlFetcher(),
		processedMap,
		ignoredFileExtMap,
	}
}

type CrawlerImpl struct {
	options           CrawlerOptions
	fetcher           UrlFetcher
	processed         map[string]bool
	ignoredFileExtMap map[string]bool
}

func (c *CrawlerImpl) Crawl(param CrawlParameters) {
	if param.Depth < 0 || !c.shouldProcessUrl(param.InputUrl) {
		return
	}

	fetchResult, err := c.fetcher.Fetch(param.InputUrl)
	if err != nil {
		log.Println(err)
		return
	}

	for _, url := range fetchResult.Urls {
		crawlParam := CrawlParameters{
			InputUrl: url,
			Depth:    param.Depth - 1,
		}
		c.Crawl(crawlParam)
	}
}

func (c *CrawlerImpl) shouldProcessUrl(url string) bool {
	if _, seen := c.processed[url]; seen {
		return false
	}
	if ext := filepath.Ext(url); c.ignoredFileExtMap[ext] {
		return false
	}
	if !c.options.FollowForeignHosts && !isSameHost(c.options.InitialUrl, url) {
		return false
	}
	c.processed[url] = true
	return true
}

func (c *CrawlerImpl) PrintAllProcessed() {
	fmt.Println("Printing All URLs Processed")
	for url, _ := range c.processed {
		fmt.Println(url)
	}
}
