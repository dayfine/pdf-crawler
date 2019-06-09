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
		if shouldPathBeIgnored(result.Path) {
			result.Path = ""
		}
		if !shouldPathBeIgnored(urlPath) {
			result.Path = path.Join(result.Path, urlPath)
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