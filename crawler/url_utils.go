package crawler

import (
	"net/url"
	"path"
	"strings"
)

func shouldPathBeIgnored(urlPath string) bool {
	return urlPath == "#" || urlPath == "/" || urlPath == "?"
}

func toUrl(base string, urlPath string) string {
	baseUrl, _ := url.Parse(base)
	pathUrl, _ := url.Parse(urlPath)

	var result *url.URL
	if strings.HasPrefix(urlPath, "//") || pathUrl.Scheme != "" {
		result = pathUrl
		if shouldPathBeIgnored(result.Path) {
			result.Path = ""
		}
	} else {
		result = baseUrl
		if shouldPathBeIgnored(result.Path) || strings.HasPrefix(urlPath, "/") {
			result.Path = ""
		}
		if !shouldPathBeIgnored(urlPath) {
			result.Path = path.Join(result.Path, pathUrl.Path)
			result.RawQuery = pathUrl.RawQuery
		}
	}

	if result.Scheme == "" {
		result.Scheme = "http"
	}
	if result.RawQuery == "" {
		result.ForceQuery = false
	}
	return result.String()
}

func isSameHost(curr string, next string) bool {
	currUrl, _ := url.Parse(toUrl(curr, ""))
	nextUrl, _ := url.Parse(toUrl(next, ""))
	currDomain := strings.Split(currUrl.Host, ":")[0]
	nextDomain := strings.Split(nextUrl.Host, ":")[0]
	return currDomain == nextDomain
}
