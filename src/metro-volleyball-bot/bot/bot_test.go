package bot_test

import (
	"testing"
	"time"

	"github.com/brewinski/home-lab-server/src/metro-volleyball-bot/bot"
	"github.com/bwmarrin/discordgo"
)

func TestBot_ReadyHandler(t *testing.T) {
	type fields struct {
		lastPageResponse string
		config           Config
	}
	type args struct {
		s     *discordgo.Session
		event *discordgo.Ready
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "TestBot_ReadyHandler will create a channel for sending updates to the server",
			fields: fields{
				lastPageResponse: "",
				config: bot.Config{
					MonitorUrl: "testing",
					UpdatesChannel: "testing", 
					TickSpeed: 1 * time.Second,
				},
			},
			args: args{
				&discordgo.Session{},
				
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bot{
				lastPageResponse: tt.fields.lastPageResponse,
				config:           tt.fields.config,
			}
			b.ReadyHandler(tt.args.s, tt.args.event)
		})
	}
}
