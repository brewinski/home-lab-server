package vq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	apiUrl       = "https://vqmetro23s3.softr.app/v1/integrations/airtable/67a0cea2-90f1-4d07-8903-89cda40f4264/appdBNmBQcBRBqB3P/Competition%20Manager/records?block_id=77226c67-17e8-4238-b983-db7105c48dfe"
	ladderApiUrl = "https://vqmetro23s3.softr.app/v1/integrations/airtable/67a0cea2-90f1-4d07-8903-89cda40f4264/appdBNmBQcBRBqB3P/Ladder/records?block_id=582e2243-42bb-4233-a84a-ddf9fc5361e1"
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
}

type ClientConfig struct {
	Client *http.Client
}

type GetGameResponseBody struct {
	Records []GameRecord `json:"records"`
	Offset  string       `json:"offset"`
}

type GetGamesRequestBody struct {
	PageSize                   int `json:"page_size"`
	AirtableResponseFormatting struct {
		Format string `json:"format"`
	} `json:"airtable_response_formatting"`
	View            string `json:"view"`
	FilterByFormula string `json:"filter_by_formula"`
	Offset          string `json:"offset"`
}

func NewClient(config ClientConfig) *Client {
	return &Client{
		config.Client,
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
		FilterByFormula: fmt.Sprintf("AND(OR(SEARCH(\"%s\", LOWER(ARRAYJOIN({LinkTmD})))))", team)}

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

type GetLadderRequestBody struct {
	PageSize                   int `json:"page_size"`
	AirtableResponseFormatting struct {
		Format string `json:"format"`
	} `json:"airtable_response_formatting"`
	View            string `json:"view"`
	FilterByFormula string `json:"filter_by_formula"`
	Rows            int    `json:"rows"`
	Offset          string `json:"offset"`
}

type GetLadderResponseBody struct {
	Records []LadderRecord `json:"records"`
	Offset  string         `json:"offset"`
}

type LadderRecord struct {
	ID     string       `json:"id"`
	Fields LadderFields `json:"fields"`
}

type LadderFields struct {
	LinkedTeam               string `json:"Linked Team"`
	MatchesA                 string `json:"MatchesA"`
	MatchesB                 string `json:"MatchesB"`
	LinkedTmDuty             string `json:"LinkedTmDuty"`
	Rank                     string `json:"Rank"`
	TeamName                 string `json:"TeamName"`
	DutyCount                string `json:"DutyCount"`
	TmA630                   string `json:"TmA-6:30"`
	TmB630                   string `json:"TmB-6:30"`
	Count630                 string `json:"Count-6:30"`
	TmA745                   string `json:"TmA-7:45"`
	TmB745                   string `json:"TmB-7:45"`
	Count745                 string `json:"Count-7:45pm"`
	TmA900                   string `json:"TmA-9:00"`
	TmB900                   string `json:"TmB-9:00"`
	Count900                 string `json:"Count-9:00"`
	Commit_IndexD            string `json:"Commit_IndexD"`
	Commit_IndexA            string `json:"Commit_IndexA"`
	NextCommitment_Type      string `json:"NextCommitment_Type"`
	Commit_RemainingD        string `json:"Commit_RemainingD"`
	Commit_DetailsD          string `json:"Commit_DetailsD"`
	Commit_RemainingA        string `json:"Commit_RemainingA"`
	Commit_DetailsA          string `json:"Commit_DetailsA"`
	Commit_DetailsB          string `json:"Commit_DetailsB"`
	NextCommitment_Detail    string `json:"NextCommitment_Detail"`
	Commit_RemainD           string `json:"Commit_RemainD"`
	Commit_RemainA           string `json:"Commit_RemainA"`
	NextCommitment_VenueTime string `json:"NextCommitment_Venue/Time"`
	DivisionClassification   string `json:"Division Classification"`
	DivisionName             string `json:"Division Name"`
	MatchesWonA              string `json:"MatchesWonA"`
	MatchesWonB              string `json:"MatchesWonB"`
	MatchesWon               string `json:"Matches Won"`
	MatchesPlayedA           string `json:"MatchesPlayedA"`
	MatchesPlayedB           string `json:"MatchesPlayedB"`
	MatchesPlayed            string `json:"Matches Played"`
	MatchesDrawnA            string `json:"MatchesDrawnA"`
	MatchesDrawnB            string `json:"MatchesDrawnB"`
	MatchesDrawn             string `json:"Matches Drawn"`
	MatchesForfeitA          string `json:"MatchesForfeitA"`
	MatchesForfeitB          string `json:"MatchesForfeitB"`
	MatchesForfeit           string `json:"MatchesForfeit"`
	MatchesDisqualifiedA     string `json:"MatchesDisqualifiedA"`
	MatchesDisqualifiedB     string `json:"MatchesDisqualifiedB"`
	MatchesDisqualified      string `json:"MatchesDisqualified"`
	TotalMatchesForfeitDQ    string `json:"Total_MatchesForfeit/DQ"`
	MatchesLost              string `json:"Matches Lost"`
	DisplayWLD               string `json:"Display_WL-DQ"`
	SetsForA                 string `json:"SetsForA"`
	SetsForB                 string `json:"SetsForB"`
	SetsFor                  string `json:"Sets For"`
	TotalSetsA               string `json:"TotalSetsA"`
	TotalSetsB               string `json:"TotalSetsB"`
	TotalSetsPlayed          string `json:"TotalSetsPlayed"`
	SetsDrawnA               string `json:"SetsDrawnA"`
	SetsDrawnB               string `json:"SetsDrawnB"`
	SetsAgainst              string `json:"Sets Against"`
	DisplaySetsWLD           string `json:"Display_SetsWLD"`
	PointsWonA               string `json:"PointsWonA"`
	PointsWonB               string `json:"PointsWonB"`
	PointsFor                string `json:"Points For"`
	PointsPlayedA            string `json:"PointsPlayedA"`
	PointsPlayedB            string `json:"PointsPlayedB"`
	TotalPointsPlayed        string `json:"TotalPointsPlayed"`
	PointsAgainst            string `json:"Points Against"`
	DisplayPtsWLD            string `json:"Display_PtsWLD"`
	SetRatio                 string `json:"SetRatio"`
	DisplaySetRatio          string `json:"Display_SetRatio"`
	PointRatio               string `json:"PointRatio"`
	DisplayPtRatio           string `json:"Display_PtRatio"`
	PtsWin                   string `json:"Pts_Win"`
	PtsLoss                  string `json:"Pts_Loss"`
	PtsDraw                  string `json:"Pts_Draw"`
	PtsSets                  string `json:"Pts_Sets"`
	PtsForfeit               string `json:"Pts_Forfeit"`
	PtsDisqualify            string `json:"Pts_Disqualify"`
	PenaltyTeamA             string `json:"Penalty_TeamA"`
	PenaltyTeamB             string `json:"Penalty_TeamB"`
	PenaltyDutyTeam          string `json:"Penalty_DutyTeam"`
	Penalties                string `json:"Penalties"`
	CompetitionPoints        string `json:"Competition Points"`
	AggregatedPoints         string `json:"Aggregated Points"`
	PoolRemaining            string `json:"PoolRemaining"`
	MaxPrediction            string `json:"Max Prediction"`
	MinPrediction            string `json:"Min Prediction"`
	PoolStatus               string `json:"PoolStatus"`
	Division                 string `json:"Division"`
	PositionalRanking        string `json:"PositionalRanking"`
	Pool                     string `json:"Pool"`
	CrossoverRanking         string `json:"CrossoverRanking"`
	NextMatch_VenueTime      string `json:"NextMatch_Venue/Time"`
	NextMatch_Detail         string `json:"NextMatch_Detail"`
	CountA                   string `json:"CountA"`
	CountB                   string `json:"CountB"`
	TotalScheduled           string `json:"TotalScheduled"`
	TeamNameLookup           string `json:"TeamNameLookup"`
}

func (c *Client) GetLadder() (GetLadderResponseBody, error) {
	ladderRequestBody, err := json.Marshal(GetLadderRequestBody{
		PageSize: 100,
		AirtableResponseFormatting: struct {
			Format string `json:"format"`
		}{
			Format: "string",
		},
		View:            "Admin - Grouped",
		FilterByFormula: "(LOWER(\"M1\") = LOWER(ARRAYJOIN({Division})))",
		Rows:            0,
		Offset:          "",
	})
	if err != nil {
		return GetLadderResponseBody{}, fmt.Errorf("GetLadder() request failed, got: %w", err)
	}

	request, err := http.NewRequest(http.MethodPost, ladderApiUrl, bytes.NewBuffer(ladderRequestBody))
	if err != nil {
		return GetLadderResponseBody{}, fmt.Errorf("GetLadder() request failed, got: %w", err)
	}

	request.Header = defaultRequestHeaders
	request.Header.Set("softr-page-id", "811945dd-84f4-44c3-b263-6fb9114ba8c3")

	res, err := c.client.Do(request)
	if err != nil {
		return GetLadderResponseBody{}, fmt.Errorf("GetLadder() request failed, got: %w", err)
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

type GameFields struct {
	TeamA                string `json:"LinkTmA"`
	TeamB                string `json:"LinkTmB"`
	DutyTeam             string `json:"LinkTmD"`
	GameTime             string `json:"Time"`
	GameDay              string `json:"Date"`
	Round                string `json:"Round"`
	Venue                string `json:"Venue"`
	Court                string `json:"Court"`
	MatchID              string `json:"Match"`
	CurrentTimeReference string `json:"CurrentTimeReference"`
	Competition          string `json:"Competition"`
}

type GameRecord struct {
	ID     string     `json:"id"`
	Fields GameFields `json:"fields"`
}

func (g GameRecord) ParseGameDayTime() (time.Time, error) {
	gameDayTime := fmt.Sprintf("%s %s", g.Fields.GameDay, g.Fields.GameTime)
	return time.Parse("2/1/2006 3:04pm", gameDayTime)
}
