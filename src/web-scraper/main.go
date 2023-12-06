package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

const HttpClientTimeout = "10s"

type ProductRedirects struct {
	Redirects []*url.URL
	Product Product
}

func main() {
	headers := http.Header{}// accept: application/json, text/plain, */*'
	headers.Add("Accept", "application/json, text/plain, */*")
	headers.Add("content-type", "application/json")

	timeout, err := time.ParseDuration(HttpClientTimeout)
	if err != nil {
		slog.Error("Main(): parsing duration", "error", err)
		return
	}

	client := http.Client{
		Timeout: timeout,	
	}
	
	productReader := New(client)

	products, err := productReader.GetProducts()
	if err != nil {
		slog.Error("Main() get products", "error", err)
	}

	productRedirects := []ProductRedirects{}

	for _, product := range products {
		url, err := url.Parse(product.Link)
		if err != nil {
			slog.Error("parse link", "error", err, "product", product)
			continue
		}

		redirects, err := RedirectReader(client, url)
		if err != nil {
			slog.Error("reading redirects", "error", err, "product", product)
			continue
		}

		productRedirects = append(productRedirects, ProductRedirects{
			Redirects: redirects,
			Product: product,
		})
	}

	fmt.Println(productRedirects)
}

func RedirectReader(client http.Client, link *url.URL) ([]*url.URL, error) {
	redirects := []*url.URL{
		link,
	}

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		slog.Warn("redirect response", "url", req.URL)
		redirects = append(redirects, req.URL)
		return nil
	}

	req, err := http.NewRequest("GET", link.String(), nil)
	if err != nil {
		slog.Error("building request", "error", err)
	}

	_, err = client.Do(req)
	if err != nil {
		slog.Error("link request", "link", link, "error", err)
		return []*url.URL{}, nil
	}

	return redirects, nil
}