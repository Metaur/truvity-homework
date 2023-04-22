package internal

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"io"
)

type LinkExtractor interface {
	Extract(content *[]byte) ([]string, error)
}

type linkExtractor struct {
}

func NewLinkExtractor() LinkExtractor {
	return &linkExtractor{}
}

func (e linkExtractor) Extract(content *[]byte) ([]string, error) {
	result := make([]string, 0)

	tokenizer := html.NewTokenizer(bytes.NewReader(*content))

	for {
		tokenType := tokenizer.Next()

		if tokenType == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				//end of the file, break out of the loop
				break
			}
			return nil, fmt.Errorf("failed to tokenize html")
		}

		token := tokenizer.Token()

		if isATag(tokenType, token) {
			if cl, ok := extractLinkFromToken(token); ok {
				result = append(result, cl)
			}
		}
	}

	return result, nil
}

func isATag(tokenType html.TokenType, token html.Token) bool {
	return tokenType == html.StartTagToken && token.DataAtom.String() == "a"
}

func extractLinkFromToken(token html.Token) (string, bool) {
	for _, attr := range token.Attr {
		if attr.Key == "href" {
			return attr.Val, true
		}
	}
	return "", false
}
