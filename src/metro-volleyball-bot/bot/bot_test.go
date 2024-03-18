package bot_test

import (
	"testing"
	"time"

	"github.com/brewinski/home-lab-server/src/metro-volleyball-bot/bot"
	"github.com/bwmarrin/discordgo"
)

// func TestBot_SessionTesting(t *testing.T) {

// 	testHandler := func(s *discordgo.Session, event *discordgo.Ready) {
// 		// no-op handler for testing
// 		fmt.Println("test handler")
// 	}

// 	// Create a new Discord session using the provided bot token.
// 	d := discordgo.Session{}

// 	d.AddHandler(testHandler)

// 	d.Handle

// }

func TestBot_ReadyHandler(t *testing.T) {
	type fields struct {
		lastPageResponse string
		config           bot.Config
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
					MonitorUrl:     "testing",
					UpdatesChannel: "testing",
					TickSpeed:      1 * time.Second,
				},
			},
			args: args{
				&discordgo.Session{},
				&discordgo.Ready{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := bot.New(bot.Config{
				MonitorUrl:     tt.fields.config.MonitorUrl,
				UpdatesChannel: tt.fields.config.UpdatesChannel,
				TickSpeed:      tt.fields.config.TickSpeed,
			})

			b.ReadyHandler(tt.args.s, tt.args.event)
		})
	}
}
