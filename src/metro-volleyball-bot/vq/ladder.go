package vq

import (
	"fmt"
	"strings"
)

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

func (ladder GetLadderResponseBody) ToString() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("Ladder: %s\n\n", "https://vqmetro23s3.softr.app/ladder-m1 ```"))

	for _, team := range ladder.Records {
		fields := team.Fields
		sb.WriteString(fmt.Sprintf("%s. %s\n", fields.Rank, team.Fields.TeamName))
		sb.WriteString(fmt.Sprintf("\tpoints: %s\n", fields.CompetitionPoints))
		sb.WriteString(fmt.Sprintf("\tnext game: %s\n", fields.NextMatch_Detail))
		sb.WriteString("\n")
	}

	sb.WriteString("```")

	return sb.String()
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
