package internal

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Loader interface {
	Load(url string) (*[]byte, error)
}

type htmlLoader struct {
}

func NewHtmlLoader() Loader {
	return &htmlLoader{}
}

func (l htmlLoader) Load(url string) (*[]byte, error) {
	fmt.Println(fmt.Sprintf("Processing %s", url))
	resp, err := http.Get(url)

	if err != nil {
		return nil, fmt.Errorf("could not fetch html for %s: %v", url, err)
	}
	defer resp.Body.Close()

	if !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		return nil, fmt.Errorf("invalid content-type, skipping %s", url)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %v", err)
	}

	return &bytes, nil
}
