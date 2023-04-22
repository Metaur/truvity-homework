package internal

import (
	"strings"
)

type Filter interface {
	Filter(urls []string) []string
}

type memoryFilter struct {
	baseUrl  string
	urlCache map[string]bool
}

func NewMemoryFilter(baseUrl string) Filter {
	return &memoryFilter{
		baseUrl:  strings.TrimSuffix(baseUrl, "/"),
		urlCache: make(map[string]bool, 0),
	}
}

func (f memoryFilter) Filter(urls []string) []string {
	result := make([]string, 0)
	f.urlCache[f.baseUrl] = true

	for _, rawUrl := range urls {
		url := formatUrl(rawUrl, f.baseUrl)
		isAllowedUrl := strings.HasPrefix(url, f.baseUrl)
		_, alreadyProcessed := f.urlCache[url]

		if isAllowedUrl && !alreadyProcessed {
			result = append(result, url)
			f.urlCache[url] = true
		}
	}

	return result
}

func formatUrl(url string, baseUrl string) string {
	formatted := strings.TrimSuffix(url, "/")
	switch {
	case strings.HasPrefix(formatted, baseUrl):
		return formatted
	case strings.HasPrefix(formatted, "/"):
		return baseUrl + formatted
	}

	return ""
}
