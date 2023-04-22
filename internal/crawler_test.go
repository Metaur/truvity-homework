package internal

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestCrawler_HappyPath(t *testing.T) {
	assertions := assert.New(t)
	dir := t.TempDir()
	url := "http://test.com"
	crawler := NewCrawler(url, 2, dir,
		&mockLoader{},
		NewFilePersister(dir),
		NewLinkExtractor(),
		NewMemoryFilter(url))

	err := crawler.Start()

	assertions.Nil(err)
	assertions.FileExists(filepath.Join(dir, "test.com", "index.html"))
	assertions.FileExists(filepath.Join(dir, "test.com", "test1.html"))
	assertions.FileExists(filepath.Join(dir, "test.com", "test2.html"))
	assertions.FileExists(filepath.Join(dir, "test.com", "test3.html"))
}

type mockLoader struct {
}

func (m mockLoader) Load(url string) (*[]byte, error) {
	bytes := []byte(`
		<html>
		<body>
			<a href="/test1"></a>
			<a href="/test2"></a>
			<a href="/test3"></a>
		<body>
		</html>
	`)
	return &bytes, nil
}
