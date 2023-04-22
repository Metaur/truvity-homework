package main

import (
	"flag"
	"fmt"
	"github.com/Metaur/truvity-homework/internal"
	"os"
	"os/signal"
	"strings"
)

func main() {
	outputDir := flag.String("output", "", "Output directory")
	url := flag.String("url", "", "Crawler url entrypoint")
	depth := flag.Int("depth", 5, "Depth")

	flag.Parse()

	fmt.Println("URL: ", *url)
	fmt.Println("Output dir: ", *outputDir)
	fmt.Println("Depth: ", *depth)

	if url == nil || len(*url) == 0 {
		panic(fmt.Errorf("url is not provided"))
	}

	signalsChan := make(chan os.Signal, 1)
	defer close(signalsChan)
	signal.Notify(signalsChan, os.Interrupt)

	sanitizedUrl := strings.TrimSuffix(*url, "/")

	loader := internal.NewHtmlLoader()
	writer := internal.NewFilePersister(*outputDir)
	extractor := internal.NewLinkExtractor()
	filter := internal.NewMemoryFilter(sanitizedUrl)
	crawler := internal.NewCrawler(sanitizedUrl, *depth, *outputDir, loader, writer, extractor, filter)

	go func() {
		for {
			s := <-signalsChan
			if s == os.Interrupt {
				crawler.Terminate()
				fmt.Println("Received termination signal, exiting...")
				return
			}
		}
	}()

	crawler.Start()

	fmt.Println("Done.")
}
