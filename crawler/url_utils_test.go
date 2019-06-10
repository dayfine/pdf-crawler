package crawler

import (
	"testing"
)

func TestToUrl(t *testing.T) {
	tests := []struct {
		baseUrl  string
		path     string
		expected string
	}{
		{"http://google.com", "", "http://google.com"},
		{"http://google.com", "/", "http://google.com"},
		{"http://google.com", "#", "http://google.com"},
		{"http://google.com", "?", "http://google.com"},
		{"google.com", "", "http://google.com"},
		{"https://google.com", "?", "https://google.com"},
		{"https://google.com", "//abc.xyz", "http://abc.xyz"},
		{"https://google.com", "https://abc.xyz", "https://abc.xyz"},
		{"http://google.com/", "", "http://google.com"},
		{"http://google.com/#", "", "http://google.com"},
		{"http://google.com/?", "", "http://google.com"},
		{"http://google.com", "/privacy", "http://google.com/privacy"},
		{"google.com", "privacy", "http://google.com/privacy"},
		{"https://google.com/internal", "policy", "https://google.com/internal/policy"},
		{"https://google.com/internal", "//abc.xyz/investor", "http://abc.xyz/investor"},
		{"https://google.com/internal", "/internal/investor", "https://google.com/internal/investor"},
		{"https://google.com/internal", "/internal/investor?q=a", "https://google.com/internal/investor?q=a"},
		{"ftp://google.com/file", "dir", "ftp://google.com/file/dir"},
	}

	for _, test := range tests {
		if actual := toUrl(test.baseUrl, test.path); actual != test.expected {
			t.Errorf("toUrl(%s, %s): expected [%s], but actual [%s]",
				test.baseUrl, test.path, test.expected, actual)
		}
	}
}
