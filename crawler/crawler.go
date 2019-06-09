package crawler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

var ignoredFileExts = []string{".mp3", ".wav", ".mkv", ".flv", ".vob", ".ogv", ".ogg", ".png", ".jpg", ".gif", ".avi", ".mov", ".wmv", ".mp4", ".mp3", ".mpg"}

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

	resp, err := c.fetcher.Fetch(param.InputUrl)
	if err != nil {
		log.Println(err)
		return
	}
	contentType := getContetType(resp)
	fmt.Printf("URL: [%s], CT: [%s]\n", param.InputUrl, contentType)
	if contentType == "application/pdf" {
		fmt.Println("Should download it here / Call handler")
		return
	} else if !isWebpage(contentType) {
		return
	}

	for _, url := range getUrls(param.InputUrl, resp) {
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

func getUrls(url string, resp *http.Response) []string {
	tokenizer := html.NewTokenizer(resp.Body)
	var urls []string
	for {
		token_type := tokenizer.Next()
		if token_type == html.ErrorToken {
			if err := tokenizer.Err(); err != io.EOF {
				log.Print(err)
			}
			return urls
		}

		switch token_type {
		case html.StartTagToken, html.SelfClosingTagToken:
			if path := getPath(tokenizer.Token()); path != "" {
				urls = append(urls, toUrl(url, path))
			}
		}
	}
}

func getPath(token html.Token) string {
	for _, attr := range token.Attr {
		switch attr.Key {
		case "href", "src":
			return strings.TrimSpace(attr.Val)
		}
	}
	return ""
}

func getContetType(resp *http.Response) string {
	ct := resp.Header.Get("Content-Type")
	return strings.Split(ct, ";")[0]
}

func isWebpage(contentType string) bool {
	return contentType == "text/html" || strings.Contains(contentType, "text/html")
}
