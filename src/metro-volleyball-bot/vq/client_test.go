package vq_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/brewinski/home-lab-server/src/metro-volleyball-bot/cfg"
	"github.com/brewinski/home-lab-server/src/metro-volleyball-bot/vq"
)

// TODO: with fake responses test get games for team.

func TestClient_GetGames_Integration(t *testing.T) {
	type fields struct {
		client *http.Client
	}
	tests := []struct {
		name    string
		fields  fields
		want    vq.GetGameResponseBody
		wantErr bool
	}{
		{
			name: "TestClient_GetGames will return a list of games",
			fields: fields{
				client: &http.Client{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := vq.NewClient(vq.ClientConfig{
				Client: tt.fields.client,
			})
			got, err := c.GetGames(2, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetGames() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got.Records) > 2 {
				t.Errorf("Client.GetGames() length is greater than requested limit limit = 2, got %d", len(got.Records))
				return
			}

			// if we didn't get an offset we can't test it.
			if got.Offset == "" {
				return
			}

			// test the offset works.
			gotOffset, err := c.GetGames(2, got.Offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetGames() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if reflect.DeepEqual(got, gotOffset) {
				t.Errorf("Client.GetGames() should be different got = %v, want %v", got, gotOffset)
				return
			}

		})
	}
}

func TestClient_GetGamesByTeam_Integration(t *testing.T) {
	type fields struct {
		client *http.Client
	}
	type args struct {
		limit  int
		offset string
		team   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    vq.GetGameResponseBody
		wantErr bool
	}{
		{
			name: "TestClient_GetGames will return a list of games",
			fields: fields{
				client: &http.Client{},
			},
			args: args{
				limit:  5,
				offset: "",
				team:   "aces",
			},
			want:    vq.GetGameResponseBody{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := vq.NewClient(vq.ClientConfig{
				Client: tt.fields.client,
			})
			got, err := c.GetGamesByTeam(tt.args.limit, tt.args.offset, tt.args.team)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetGamesByTeam() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// for each game we expect the team to be playing
			for _, game := range got.Records {
				lowerTeam := strings.ToLower(tt.args.team)
				lowerTeamA := strings.ToLower(game.Fields.TeamA)
				lowerTeamB := strings.ToLower(game.Fields.TeamB)

				if !strings.Contains(lowerTeamA, lowerTeam) && !strings.Contains(lowerTeamB, lowerTeam) {
					t.Errorf("Client.GetGamesByTeam() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
		})
	}
}

func TestClient_GetGamesByTeamAndDuty_Integration(t *testing.T) {
	type fields struct {
		client *http.Client
	}
	type args struct {
		limit  int
		offset string
		team   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    vq.GetGameResponseBody
		wantErr bool
	}{
		{
			name: "TestClient_GetGames will return a list of games",
			fields: fields{
				client: &http.Client{},
			},
			args: args{
				limit:  10,
				offset: "",
				team:   "apg",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := vq.NewClient(vq.ClientConfig{
				Client: tt.fields.client,
			})

			got, err := c.GetGamesByTeamAndDuty(tt.args.limit, tt.args.offset, tt.args.team)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetGamesByTeamAndDuty() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for _, game := range got.Records {
				lowerTeam := strings.ToLower(tt.args.team)
				lowerDutyTeam := strings.ToLower(game.Fields.DutyTeam)
				lowerTeamA := strings.ToLower(game.Fields.TeamA)
				lowerTeamB := strings.ToLower(game.Fields.TeamB)
				// for every game we expect the duty team to be the requested team
				if !strings.Contains(lowerDutyTeam, lowerTeam) {
					t.Errorf("Client.GetGamesByTeamAndDuty() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				// for every game we expect the requested team not to be playing
				if lowerTeamA == lowerTeam || lowerTeamB == lowerTeam {
					t.Errorf("Client.GetGamesByTeamAndDuty() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}

		})
	}
}

var ladderResponseBody = vq.GetLadderResponseBody{
	Records: []vq.LadderRecord{
		{
			ID: "1",
			Fields: vq.LadderFields{
				Rank:     "1",
				TeamName: "Aces",
			},
		},
		{
			ID: "2",
			Fields: vq.LadderFields{
				Rank:     "1",
				TeamName: "APG",
			},
		},
	},
	Offset: "",
}

func TestClient_GetLadder_Integration(t *testing.T) {
	type fields struct {
		// client config
		client  *http.Client
		baseUrl string
		// test server config. (optional) if this is set the test server will be used instead of the base url.
		testServer *httptest.Server
	}
	tests := []struct {
		name    string
		fields  fields
		want    vq.GetLadderResponseBody
		wantErr bool
	}{
		{
			name: "TestClient_GetLadder will return a list of teams",
			fields: fields{
				client: &http.Client{},
				testServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if pageId := r.Header.Get("softr-page-id"); pageId == "" {
						w.WriteHeader(http.StatusBadRequest)
						w.Write([]byte("missing softr-page-id header"))
						return
					}

					ladderBytes, err := json.Marshal(ladderResponseBody)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte("error getting ladder"))
						return
					}

					w.WriteHeader(http.StatusOK)
					w.Write(ladderBytes)
				})),
			},
			want:    vq.GetLadderResponseBody{},
			wantErr: false,
		},
		{
			name: "TestClient_GetLadder will return a list of teams",
			fields: fields{
				client:  &http.Client{},
				baseUrl: cfg.VQBaseUrl,
			},
			want:    vq.GetLadderResponseBody{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testServer := tt.fields.testServer
			// we need to close the test server after we are done with it.
			if testServer != nil {
				defer testServer.Close()
				// set the testing url to the url of the test server
				tt.fields.baseUrl = testServer.URL
			}

			c := vq.NewClient(vq.ClientConfig{
				Client: tt.fields.client,
				ApiUrl: tt.fields.baseUrl,
			})

			got, err := c.GetLadder()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetLadder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for _, team := range got.Records {
				fmt.Printf("rank[%s], team[%s], points[%s], next_game[%s]\n\n", team.Fields.Rank, team.Fields.TeamNameLookup, team.Fields.CompetitionPoints, team.Fields.NextMatch_Detail)
			}
		})
	}
}
