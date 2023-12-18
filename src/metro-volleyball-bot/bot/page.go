package bot

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type Web struct {
	client *http.Client
}

type WebConfig struct {
	Client *http.Client
}

func NewWeb(config WebConfig) *Web {
	return &Web{
		config.Client,
	}
}

// MonitorPage will monitor a specific page for changes.
func (w *Web) MonitorPage(pageUrl string) (string, error) {
	response, err := w.client.Get(pageUrl)
	if err != nil {
		return "", fmt.Errorf("MonitorPage() request failed, got: %w", err)
	}

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return "", fmt.Errorf("MonitorPage() response status code, got: %d", response.StatusCode)
	}

	defer response.Body.Close()

	respBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	slog.Info("page response", "status", response.Status, "content-length", response.ContentLength)

	return string(respBytes), nil
}

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
