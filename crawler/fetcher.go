package crawler

import (
	"net/http"
)

type UrlFetcher interface {
	// Returns slice of URLs on the page of requested URL.
	Fetch(url string) (*Resource, error)
}

func NewUrlFetcher() UrlFetcher {
	return &PlainUrlFetcher{}
}

// A UrlFetcher that simply finds URLs on the HTML page.
type PlainUrlFetcher struct {
}

func (f *PlainUrlFetcher) Fetch(url string) (*Resource, error) {
	// TODO(dayfine): use channel / go routine
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return &Resource{url, resp}, nil
}

// A UrlFetcher that finds URLs by rendering any JavaScript.
type RenderedUrlFetcher struct {
	// TODO(dayfine): implement this
}

// A UrlFetcher that finds URLs by rendering any JavaScript and click around.
type ClickUrlFetcher struct {
	// TODO(dayfine): implement this
}
