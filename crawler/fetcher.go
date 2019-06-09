package crawler

import (
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type FetchResult struct {
	PageUrl     string
	ContentType string
	Urls        []string
}

type UrlFetcher interface {
	// Returns slice of URLs on the page of requested URL.
	Fetch(url string) (*FetchResult, error)
}

func NewUrlFetcher() UrlFetcher {
	return &PlainUrlFetcher{}
}

// A UrlFetcher that simply finds URLs on the HTML page.
type PlainUrlFetcher struct {
}

func (f *PlainUrlFetcher) Fetch(url string) (*FetchResult, error) {
	// TODO(dayfine): use channel / go routine
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result := &FetchResult{
		PageUrl:     url,
		ContentType: getContetType(resp),
	}
	if !shouldAnalyze(result) {
		return result, nil
	}

	tokenizer := html.NewTokenizer(resp.Body)
	for {
		token_type := tokenizer.Next()
		if token_type == html.ErrorToken {
			if err := tokenizer.Err(); err != io.EOF {
				return result, err
			}
			return result, nil
		}

		switch token_type {
		case html.StartTagToken, html.SelfClosingTagToken:
			if path := getPath(tokenizer.Token()); path != "" {
				result.Urls = append(result.Urls, toUrl(url, path))
			}
		}
	}
}

// A UrlFetcher that finds URLs by rendering any JavaScript.
type RenderedUrlFetcher struct {
	// TODO(dayfine): implement this
}

// A UrlFetcher that finds URLs by rendering any JavaScript and click around.
type ClickUrlFetcher struct {
	// TODO(dayfine): implement this
}

func shouldAnalyze(result *FetchResult) bool {
	return isWebpage(result.ContentType)
}

func getContetType(resp *http.Response) string {
	ct := resp.Header.Get("Content-Type")
	return strings.Split(ct, ";")[0]
}

func isWebpage(contentType string) bool {
	return contentType == "text/html" || strings.Contains(contentType, "text/html")
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
