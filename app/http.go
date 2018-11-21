package app

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func download(file *os.File, url string) error {
	fmt.Printf("Downloading graph xml: %s\n", url)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Invalid status downloading xml: %v", resp.Status))
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func isValidUrl(toTest string) bool {
	u, err := url.Parse(toTest)

	if err != nil {
		return false
	}

	if u.Hostname() == "" {
		return false
	}

	return true
}
