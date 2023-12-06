package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type TableResponse struct {
	Data struct{
		Table Table 
	}
}

type Table struct {
	Products []Product
}

type Product struct {
	Link string
}

type GraphQLProductReader struct {
	httpClient http.Client
}

func New(http http.Client) *GraphQLProductReader {
	return &GraphQLProductReader{
		httpClient: http,
	}
}

func (gpr GraphQLProductReader) GetProducts() ([]Product, error) {
	// build the GQL request object
	request, err := http.NewRequest(
		"GET", 
		"https://graph.canstar.com.au/graphql?operationName=Table&variables=%7B%22vertical%22:%22mortgages%22,%22selectorFields%22:%5B%7B%22name%22:%22Loan%20amount%22,%22type%22:%22string%22,%22value%22:%22500000%22%7D,%7B%22name%22:%22Loan%20purpose%22,%22type%22:%22string%22,%22value%22:%22Refinance%22%7D,%7B%22name%22:%22State%22,%22type%22:%22string%22,%22value%22:%22New%2520South%2520Wales%22%7D%5D,%22filterFields%22:%5B%7B%22name%22:%22Online%20Partner%22,%22type%22:%22string%22,%22value%22:%22true%22%7D,%7B%22name%22:%22Loan%20type%22,%22type%22:%22string%22,%22value%22:%22Variable%22%7D,%7B%22name%22:%22LVR%22,%22type%22:%22string%22,%22value%22:%2280%2525%22%7D,%7B%22name%22:%22LVR%22,%22type%22:%22string%22,%22value%22:%2270%2525%22%7D,%7B%22name%22:%22LVR%22,%22type%22:%22string%22,%22value%22:%2260%2525%2520or%2520less%22%7D,%7B%22name%22:%22Repayment%20type%22,%22type%22:%22string%22,%22value%22:%22Principal%2520%2526%2520Interest%22%7D%5D,%22sort%22:%5B%5D,%22pagination%22:null,%22featureFlags%22:%5B%5D%7D&extensions=%7B%22persistedQuery%22:%7B%22version%22:1,%22sha256Hash%22:%22c4d755441e43406ca02b28e9884f72fb2aeec993cdc5108681560fa6f89a883a%22%7D%7D",
		nil,
	)
	if err != nil {
		return []Product{}, fmt.Errorf("GetProducts: %w", err)
	}

	request.Header.Add("Accept", "application/json, text/plain, */*")
	request.Header.Add("content-type", "application/json")

	resp, err := gpr.httpClient.Do(request)
	if err != nil {
		return []Product{}, fmt.Errorf("GetProducts: %w", err)
	}

	if resp.StatusCode > 300 {
		return []Product{}, fmt.Errorf("GetProducts, %v", resp)
	}

	tresp := TableResponse{}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return []Product{}, fmt.Errorf("GetProducts: %w", err)
	}
	
	err = json.NewDecoder(bytes.NewReader(respBody)).Decode(&tresp)
	if err != nil {
		return []Product{}, fmt.Errorf("GetProducts:  %w", err)
	}	

	return tresp.Data.Table.Products, nil
}

// func ReadProducts() {
// 		// creating a new Colly instance 
// 		c := colly.NewCollector() 

// 		c.OnResponse(func (r *colly.Response) {
// 			tresp := TableResponse{}
	
// 			err := json.NewDecoder(bytes.NewReader(r.Body)).Decode(&tresp)
// 			if err != nil {
// 				slog.Error("decodeing table response", "error", err)
// 				return
// 			}
	
// 			slog.Info("decoded response data", "data", tresp)
// 			products := tresp.Data.Table.Products
	
// 			client := http.Client{
// 				CheckRedirect: func(req *http.Request, via []*http.Request) error {
// 					slog.Warn("redirect response", "url", req.URL)
// 					return nil
// 				},
				
// 			}
	
// 			for _, product := range products {
// 				slog.Info("processing", "product", product)
// 				if product.Link != "" {
// 					resp, err := client.Get(product.Link)
// 					if err != nil {
// 						slog.Error("link request", "product", product, "error", err)
// 					}
	
// 					slog.Info("request response", "response", resp.Header.Get("Location"), "status_code", resp.StatusCode, "status", resp.Status)
// 				}
	
// 				if product.Link != "" {
// 					body := bytes.NewReader([]byte{})
// 					requestR, err := http.NewRequest("GET", product.Link, body)
// 					if err != nil {
// 						slog.Error("building request", "error", err)
// 					}
	
// 					resp, err := client.Do(requestR)
// 					if err != nil {
// 						slog.Error("link request", "product", product, "error", err)
// 					}
	
// 					slog.Info("request response", "response", resp.Header.Get("Location"), "status_code", resp.StatusCode, "status", resp.Status)
// 				}
// 			}
	
// 		})
	
// 		reader := bytes.NewReader([]byte("test"))
	
// 		headers := http.Header{}// accept: application/json, text/plain, */*'
// 		headers.Add("Accept", "application/json, text/plain, */*")
// 		headers.Add("content-type", "application/json")
	
// 		err := c.Request(
// 			"GET", 
// 			"https://graph.canstar.com.au/graphql?operationName=Table&variables=%7B%22vertical%22:%22mortgages%22,%22selectorFields%22:%5B%7B%22name%22:%22Loan%20amount%22,%22type%22:%22string%22,%22value%22:%22500000%22%7D,%7B%22name%22:%22Loan%20purpose%22,%22type%22:%22string%22,%22value%22:%22Refinance%22%7D,%7B%22name%22:%22State%22,%22type%22:%22string%22,%22value%22:%22New%2520South%2520Wales%22%7D%5D,%22filterFields%22:%5B%7B%22name%22:%22Online%20Partner%22,%22type%22:%22string%22,%22value%22:%22true%22%7D,%7B%22name%22:%22Loan%20type%22,%22type%22:%22string%22,%22value%22:%22Variable%22%7D,%7B%22name%22:%22LVR%22,%22type%22:%22string%22,%22value%22:%2280%2525%22%7D,%7B%22name%22:%22LVR%22,%22type%22:%22string%22,%22value%22:%2270%2525%22%7D,%7B%22name%22:%22LVR%22,%22type%22:%22string%22,%22value%22:%2260%2525%2520or%2520less%22%7D,%7B%22name%22:%22Repayment%20type%22,%22type%22:%22string%22,%22value%22:%22Principal%2520%2526%2520Interest%22%7D%5D,%22sort%22:%5B%5D,%22pagination%22:null,%22featureFlags%22:%5B%5D%7D&extensions=%7B%22persistedQuery%22:%7B%22version%22:1,%22sha256Hash%22:%22c4d755441e43406ca02b28e9884f72fb2aeec993cdc5108681560fa6f89a883a%22%7D%7D",
// 			reader,
// 			colly.NewContext(),
// 			headers,
// 		)
	
// 		if err != nil {
// 			slog.Error("request failed", "error", err)
// 			return
// 		}
// }