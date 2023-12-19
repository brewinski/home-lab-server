package vq_test

import (
	"fmt"
	"net/http"
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
			want:    vq.GetGameResponseBody{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := vq.NewClient(vq.ClientConfig{
				Client: tt.fields.client,
			})
			got, err := c.GetGames(100, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetGames() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// TODO: bring this back for real tests
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("Client.GetGames() = %v, want %v", got, tt.want)
			// }

			teamGamesMap := map[string][]vq.Game{}

			for _, game := range got.Records {
				teamGamesMap[game.Fields.TeamA] = append(teamGamesMap[game.Fields.TeamA], game)
				teamGamesMap[game.Fields.TeamB] = append(teamGamesMap[game.Fields.TeamB], game)
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

			futureGames := []vq.Game{}
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
			want:    vq.GetGameResponseBody{},
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
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("Client.GetGamesByTeamAndDuty() = %v, want %v", got, tt.want)
			// }

			futureGames := []vq.Game{}
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
