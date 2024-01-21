package main

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
)

// Fetch returns the content of a url as a string
func Fetch(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	var bodyString string
	if resp.StatusCode == 200 {
		buf := new(bytes.Buffer)
		_, err := buf.ReadFrom(resp.Body)
		if err != nil {
			return "", err
		}
		bodyString = buf.String()
	} else {
		slog.Warn(fmt.Sprintf("%s returned status code %d", url, resp.StatusCode))
		return "", nil
	}

	return bodyString, nil
}
