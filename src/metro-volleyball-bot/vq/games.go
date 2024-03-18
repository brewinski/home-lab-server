package vq

import (
	"fmt"
	"time"
)

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
