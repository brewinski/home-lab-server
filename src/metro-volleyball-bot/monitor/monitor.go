package monitor

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
)

type DataSourceMonitor struct {
	client       *http.Client
	prevResponse string
	mu           sync.Mutex
}

type Config struct {
	Client   *http.Client
	InitData string
}

func New(config Config) *DataSourceMonitor {
	return &DataSourceMonitor{
		config.Client,
		config.InitData,
		sync.Mutex{},
	}
}

func (w *DataSourceMonitor) CheckForChanges(url string) (bool, string, error) {
	response, err := w.Monitor(url)
	if err != nil {
		return false, w.prevResponse, fmt.Errorf("CheckPageForChanges() request failed, got: %w", err)
	}

	// lock the mutex before reading / writing to shared memory.
	// used for cases where this object might be used in a concurrent environment.
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.prevResponse == "" {
		w.prevResponse = response
		return false, w.prevResponse, nil
	}

	if w.prevResponse == response {
		return false, w.prevResponse, nil
	}

	w.prevResponse = response
	return true, w.prevResponse, nil
}

// Monitor will monitor a specific page for changes.
func (w *DataSourceMonitor) Monitor(pageUrl string) (string, error) {
	response, err := w.client.Get(pageUrl)
	if err != nil {
		return "", fmt.Errorf("MonitorPage() request failed, got: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return "", fmt.Errorf("MonitorPage() response status code, got: %d", response.StatusCode)
	}

	respBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	slog.Info("page response", "status", response.Status, "content-length", response.ContentLength)

	return string(respBytes), nil
}
