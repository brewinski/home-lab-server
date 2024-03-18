package vq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/brewinski/home-lab-server/src/metro-volleyball-bot/cfg"
)

const (
	// TODO: replace api url with vq api url / env var
	apiUrl = "https://vqmetro23s3.softr.app/v1/integrations/airtable/67a0cea2-90f1-4d07-8903-89cda40f4264/appdBNmBQcBRBqB3P/Competition%20Manager/records?block_id=77226c67-17e8-4238-b983-db7105c48dfe"
	// TODO: use these vars combined to document the api and paths.
	vqAPIUrl        = "https://vqmetro23s3.softr.app/v1/integrations/airtable/67a0cea2-90f1-4d07-8903-89cda40f4264/appdBNmBQcBRBqB3P"
	competitionPath = "/Competition%20Manager/records?block_id=77226c67-17e8-4238-b983-db7105c48dfe"
)

var (
	defaultRequestHeaders = http.Header{
		"accept":             {"application/json, text/plain, */*"},
		"accept-language":    {"en-GB,en-US;q=0.9,en;q=0.8"},
		"content-type":       {"application/json"},
		"sec-ch-ua":          {"\"Not_A Brand\";v=\"8\", \"Chromium\";v=\"120\", \"Google Chrome\";v=\"120\""},
		"sec-ch-ua-mobile":   {"?0"},
		"sec-ch-ua-platform": {"\"macOS\""},
		"sec-fetch-dest":     {"empty"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-site":     {"same-origin"},
	}
)

type Client struct {
	client *http.Client
	apiUrl string
}

type ClientConfig struct {
	Client *http.Client
	ApiUrl string
}

func NewClient(config ClientConfig) *Client {
	if config.Client == nil {
		config.Client = &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	return &Client{
		config.Client,
		config.ApiUrl,
	}
}

func (c *Client) GetGamesByTeamAndDuty(limit int, offset, team string) (GetGameResponseBody, error) {
	gameReqBody := GetGamesRequestBody{
		PageSize: limit,
		AirtableResponseFormatting: struct {
			Format string `json:"format"`
		}{
			Format: "string",
		},
		View:            "RMS - Timeslot",
		Offset:          offset,
		FilterByFormula: fmt.Sprintf("AND(OR(SEARCH(LOWER(\"%s\"), LOWER(ARRAYJOIN({LinkTmD})))))", team)}

	return c.getGames(gameReqBody)
}

// get games filtered by team with a limit and an offset. maximum limit is 100.
func (c *Client) GetGamesByTeam(limit int, offset, team string) (GetGameResponseBody, error) {
	gameReqBody := GetGamesRequestBody{
		PageSize: limit,
		AirtableResponseFormatting: struct {
			Format string `json:"format"`
		}{
			Format: "string",
		},
		View:            "RMS - Timeslot",
		Offset:          offset,
		FilterByFormula: fmt.Sprintf("AND(OR(SEARCH(\"%s\", LOWER(ARRAYJOIN({Display_Identifier})))))", team),
	}

	return c.getGames(gameReqBody)
}

// get all games with a limit and an offset. maximum limit is 100.
func (c *Client) GetGames(limit int, offset string) (GetGameResponseBody, error) {
	return c.getGames(GetGamesRequestBody{
		PageSize: limit,
		AirtableResponseFormatting: struct {
			Format string `json:"format"`
		}{
			Format: "string",
		},
		View:   "RMS - Timeslot",
		Offset: offset,
	})
}

func (c *Client) getGames(reqBody GetGamesRequestBody) (GetGameResponseBody, error) {
	gameRequestBody, err := json.Marshal(reqBody)
	if err != nil {
		return GetGameResponseBody{}, fmt.Errorf("GetGames() request failed, got: %w", err)
	}

	request, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(gameRequestBody))
	if err != nil {
		return GetGameResponseBody{}, fmt.Errorf("GetGames() request failed, got: %w", err)
	}

	request.Header = defaultRequestHeaders
	request.Header.Set("softr-page-id", "f787ee5d-5f9d-4fe1-95f4-22a2ceeb4d01")

	res, err := c.client.Do(request)
	if err != nil {
		return GetGameResponseBody{}, fmt.Errorf("GetGames() request failed, got: %w", err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return GetGameResponseBody{}, fmt.Errorf("GetGames() request failed, got: %w", err)
	}

	defer res.Body.Close()

	var games GetGameResponseBody

	err = json.Unmarshal(body, &games)
	if err != nil {
		return GetGameResponseBody{}, fmt.Errorf("GetGames() request failed, got: %w", err)
	}
	return games, nil
}

func (c *Client) GetLadder() (GetLadderResponseBody, error) {
	ladderRequestBody, err := json.Marshal(GetLadderRequestBody{
		PageSize: 100,
		AirtableResponseFormatting: struct {
			Format string `json:"format"`
		}{
			Format: "string",
		},
		View:            "Division Ranking",
		FilterByFormula: cfg.VQLadderFilterByFormula,
		Rows:            0,
		Offset:          "",
	})

	if err != nil {
		return GetLadderResponseBody{}, fmt.Errorf("GetLadder() request failed, got: %w", err)
	}

	request, err := http.NewRequest(http.MethodPost, c.apiUrl+cfg.VQLadderPath, bytes.NewBuffer(ladderRequestBody))
	if err != nil {
		return GetLadderResponseBody{}, fmt.Errorf("GetLadder() request failed, got: %w", err)
	}

	request.Header = defaultRequestHeaders
	request.Header.Set("softr-page-id", cfg.VQLadderPageID)

	res, err := c.client.Do(request)
	if err != nil {
		return GetLadderResponseBody{}, fmt.Errorf("GetLadder() request failed, got: %w", err)
	}

	if res.StatusCode > 399 {
		body, _ := io.ReadAll(res.Body)
		return GetLadderResponseBody{}, fmt.Errorf("GetLadder() request failed with status code %d, response body: %s", res.StatusCode, body)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return GetLadderResponseBody{}, fmt.Errorf("GetLadder() request failed, got: %w", err)
	}

	defer res.Body.Close()

	var ladder GetLadderResponseBody

	err = json.Unmarshal(body, &ladder)
	if err != nil {
		return GetLadderResponseBody{}, fmt.Errorf("GetLadder() request failed, got: %w", err)
	}

	return ladder, nil
}
