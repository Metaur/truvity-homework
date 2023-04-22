package internal

import "fmt"

type Crawler interface {
	Start() error
	Terminate()
}

type crawler struct {
	startUrl  string
	depth     int
	targetDir string

	loader    Loader
	writer    Persister
	extractor LinkExtractor
	filter    Filter

	shouldStop bool
}

func NewCrawler(startUrl string, depth int, targetDir string,
	loader Loader, writer Persister, extractor LinkExtractor, filter Filter) Crawler {
	return &crawler{
		startUrl:   startUrl,
		depth:      depth,
		targetDir:  targetDir,
		loader:     loader,
		writer:     writer,
		extractor:  extractor,
		filter:     filter,
		shouldStop: false,
	}
}

func (c crawler) Start() error {
	return c.crawl(c.startUrl, 0)
}

func (c crawler) Terminate() {
	c.shouldStop = true
}

func (c crawler) crawl(url string, currentDepth int) error {
	if currentDepth > c.depth || c.shouldStop {
		return nil
	}

	response, err := c.loader.Load(url)
	if err != nil {
		fmt.Println("got error during load", err)
		return err
	}

	err = c.writer.Write(url, response)
	if err != nil {
		fmt.Println("got error during write", err)
	}

	links, err := c.extractor.Extract(response)
	if err != nil {
		fmt.Println("got error during extract", err)
		return err
	}

	result := c.filter.Filter(links)

	fmt.Println(fmt.Sprintf("Extracted urls %v, count %d", result, len(result)))

	for _, nextUrl := range result {
		err := c.crawl(nextUrl, currentDepth+1)
		if err != nil {
			continue
		}
	}

	return nil
}
