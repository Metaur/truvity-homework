package internal

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
)

type Persister interface {
	Write(url string, content *[]byte) error
}

type filePersister struct {
	targetFolder string
}

func NewFilePersister(targetFolder string) Persister {
	return &filePersister{targetFolder: targetFolder}
}

func (p filePersister) Write(rawUrl string, content *[]byte) error {
	parsed, err := url.Parse(rawUrl)
	if err != nil {
		return fmt.Errorf("could not parse url %s: %v", rawUrl, err)
	}

	folder := filepath.Join(p.targetFolder, parsed.Host, filepath.Dir(parsed.Path))
	var filename string
	if len(parsed.Path) == 0 {
		filename = filepath.Join(folder, "index.html")
	} else {
		// contains path to a target file without extension
		filename = filepath.Join(folder, filepath.Base(parsed.Path)+".html")
	}
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		err := os.MkdirAll(folder, 0700)
		if err != nil {
			return fmt.Errorf("could not create folder: %v", err)
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not create a file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(*content)
	if err != nil {
		return fmt.Errorf("could not write a file: %v", err)
	}

	return nil
}
