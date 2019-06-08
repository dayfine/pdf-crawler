package crawler

import (
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type FetchResult struct {
	Body string
	Urls []string
}

type UrlFetcher interface {
	// Returns the body of URL and a slice of URLs on that page
	Fetch(url string) (result *FetchResult, err error)
}

func NewUrlFetcher() UrlFetcher {
	return &UrlFetchImpl{}
}

type UrlFetchImpl struct {
	proccessed map[string]bool
}

func (f *UrlFetchImpl) Fetch(url string) (*FetchResult, error) {
	// _, seen := f.proccessed[url]
	// if seen {
	// 	return nil, nil
	// }
	// f.proccessed[url] = true

	// TODO(dayfine): use channel / go routine
	resp, err := http.Get(url)
	if err != nil || !shouldAnalyze(resp) {
		return nil, err
	}
	defer resp.Body.Close()

	result := &FetchResult{}
	// bodyBytes, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return nil, err
	// }
	// result.Body = string(bodyBytes)

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
				result.Urls = append(result.Urls, buildUrlFromPath(path))
			}
		}
	}
}

func shouldAnalyze(resp *http.Response) bool {
	return isWebpage(resp.Header.Get("Content-Type"))
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

func buildUrlFromPath(path string) string {
	return path
}
