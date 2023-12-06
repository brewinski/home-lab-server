package main

import (
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const (
	PING_FREQUENCY = 5 * time.Second
	PAGE_URL = "https://www.vq.org.au/competitions/metro-league/"
)

var (
	lastData string
)

func main() {
	for range time.Tick(PING_FREQUENCY) {
		slog.Info("checking for changes", "url", PAGE_URL)
		response, err := http.Get(PAGE_URL)
		if err != nil {
			slog.Error("api request failed", "error", err)
			os.Exit(1);
		} 
		
		respBytes, err := io.ReadAll(response.Body)
		if err != nil {
			slog.Error("api response read failed", "error", err)
			os.Exit(1);
		}

		slog.Info(
			"response", 
			"status", response.Status, 
			"content-length", response.ContentLength, 
		)

		if lastData == "" {
			slog.Info("initial data set")
			lastData = string(respBytes)
		} 
		
		if lastData != string(respBytes) {
			slog.Info("data changed")
			lastData = string(respBytes)
		} else {
			slog.Info("data unchanged")
		}
			
		response.Body.Close()
	}
}

// 1168314085151100988

// MTE2ODMxNDA4NTE1MTEwMDk4OA.G5pIeG.hfqGGjZuQDOWoA2yYatNfQy0oLe0YJuaGLPnpY

// https://discord.com/api/oauth2/authorize?client_id=1168314085151100988&permissions=8&scope=bot