package vq_test

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/brewinski/home-lab-server/src/metro-volleyball-bot/vq"
)

func TestClient_GetGames(t *testing.T) {
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

func TestClient_GetGamesByTeam(t *testing.T) {
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
				limit:  100,
				offset: "",
				team:   "redbacks",
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
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("Client.GetGamesByTeam() = %v, want %v", got, tt.want)
			// }

			futureGames := []vq.GameRecord{}
			for _, game := range got.Records {
				// fmt.Printf("%s vs %s, at %s on %s\n", game.Fields.TeamA, game.Fields.TeamB, game.Fields.GameTime, game.Fields.GameDay)
				gameTime, err := game.ParseGameDayTime()
				if err != nil {
					t.Errorf("Client.GetGamesByTeam() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if gameTime.Before(time.Now()) {
					continue
				}

				futureGames = append(futureGames, game)
				fmt.Printf("game:[%s vs %s] time:[%s %s]\n", game.Fields.TeamA, game.Fields.TeamB, game.Fields.GameDay, game.Fields.GameTime)
			}

			fmt.Printf("future games: %v\n", futureGames)
		})
	}
}

func TestClient_GetGamesByTeamAndDuty(t *testing.T) {
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
				limit:  100,
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

func TestClient_GetLadder(t *testing.T) {
	type fields struct {
		client *http.Client
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
			},
			want:    vq.GetLadderResponseBody{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := vq.NewClient(vq.ClientConfig{
				Client: tt.fields.client,
			})

			got, err := c.GetLadder()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetLadder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("Client.GetLadder() = %v, want %v", got, tt.want)
			// }

			for _, team := range got.Records {
				fmt.Printf("rank[%s], team[%s], points[%s], next_game[%s]\n\n", team.Fields.Rank, team.Fields.TeamNameLookup, team.Fields.CompetitionPoints, team.Fields.NextMatch_Detail)
			}
		})
	}
}
