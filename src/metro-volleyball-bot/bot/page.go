package bot

import (
	"io"
	"log/slog"
	"net/http"
)

// ping the page and check if anything has changed.
// it will return false if the page hasn't been checked in the past
func MonitorPage(pageUrl string) (string, error) {
	response, err := http.Get(pageUrl)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	respBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	slog.Info("page response", "status", response.Status, "content-length", response.ContentLength)

	return string(respBytes), nil
}
