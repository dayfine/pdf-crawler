package crawler

import (
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
}

type Resource struct {
	url  string
	resp *http.Response
}

type CrawlerImpl struct {
	options           CrawlerOptions
	fetcher           UrlFetcher
	processed         map[string]bool
	ignoredFileExtMap map[string]bool
	handlers          map[string]Handler
}

func NewCrawler(options CrawlerOptions) Crawler {
	processedMap := make(map[string]bool)
	ignoredFileExtMap := make(map[string]bool)
	for _, ext := range ignoredFileExts {
		ignoredFileExtMap[ext] = true
	}
	handlerMap := make(map[string]Handler)
	handlerMap["pdf"] = NewPdfFileDownloadHandler(options.SaveDir)

	return &CrawlerImpl{
		options,
		NewUrlFetcher(),
		processedMap,
		ignoredFileExtMap,
		handlerMap,
	}
}

func (c *CrawlerImpl) Crawl(param CrawlParameters) {
	if param.Depth < 0 || !c.shouldProcessUrl(param.InputUrl) {
		return
	}

	res, err := c.fetcher.Fetch(param.InputUrl)
	if err != nil {
		log.Println(err)
		return
	}

	c.handleResource(res)
	if isWebpage(res) {
		for _, url := range getUrls(res) {
			crawlParam := CrawlParameters{
				InputUrl: url,
				Depth:    param.Depth - 1,
			}
			c.Crawl(crawlParam)
		}
	}
}

func (c *CrawlerImpl) handleResource(res *Resource) {
	if isPdf(res) {
		err := c.handlers["pdf"].Handle(res)
		if err != nil {
			log.Println(err)
			return
		}
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

func getUrls(res *Resource) []string {
	tokenizer := html.NewTokenizer(res.resp.Body)
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
				urls = append(urls, toUrl(res.url, path))
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

func getContetType(res *Resource) string {
	ct := res.resp.Header.Get("Content-Type")
	return strings.Split(ct, ";")[0]
}

func isWebpage(res *Resource) bool {
	contentType := getContetType(res)
	return contentType == "text/html" || strings.Contains(contentType, "text/html")
}

func isPdf(res *Resource) bool {
	contentType := getContetType(res)
	return contentType == "application/pdf" || strings.Contains(res.url, ".pdf")
}
