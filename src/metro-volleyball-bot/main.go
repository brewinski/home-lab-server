package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"time"

	"github.com/brewinski/home-lab-server/src/metro-volleyball-bot/bot"
	"github.com/brewinski/home-lab-server/src/metro-volleyball-bot/vq"
	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token                string
	TickSpeed            time.Duration
	PageUrl              string
	NotificationsChannel string
	VQClientUrl          string
)

func main() {
	// Read the Discord token as an application flag.
	flag.StringVar(&Token, "t", "", "The token for the specific discord application.")
	// Read the page check frequency duration. Parsed as "1ms", "1ns", "1s", "1m", or "1h"
	flag.DurationVar(&TickSpeed, "ts", 1*time.Hour, "Page ping frequency as a string duration")
	// Page URL to monitor
	flag.StringVar(&PageUrl, "url", PageUrl, "The URL to monitor for changes")
	// VQ Metro Draw Data Base Url
	flag.StringVar(
		&VQClientUrl,
		"draw-url",
		"https://vqmetro23s3.softr.app/v1/integrations/airtable/67a0cea2-90f1-4d07-8903-89cda40f4264/appdBNmBQcBRBqB3P",
		"The channel to send notifications",
	)
	// channel to publish notifications to
	flag.StringVar(&NotificationsChannel, "channel", "volleybot-notifications", "The channel to send notifications")
	// Parse the flags from the command line
	flag.Parse()

	// set json as the default logger
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(
		slog.String("tick_speed", TickSpeed.String()),
		slog.String("channel", NotificationsChannel),
		slog.String("page", PageUrl),
	))

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		slog.Error("new discord session", "error", err)
		return
	}
	// Cleanly close down the Discord session.
	defer dg.Close()

	myBot := bot.New(bot.Config{
		UpdatesChannel: NotificationsChannel,
		TickSpeed:      TickSpeed,
		MonitorUrl:     PageUrl,
	})

	// register the bot ready handler
	dg.AddHandler(myBot.ReadyHandler)
	// commands handler
	commandHandler := bot.OnCommandHandlerFactory(func(command string) (string, error) {
		switch command {
		case "vb-help":
			return "no action registered for this command: " + command, nil
		case "vb-ladder":
			return "no action registered for this command: " + command, nil
		default:
			// Create the response object
			return "no action registered for this command" + command, nil
		}
	})
	dg.AddHandler(commandHandler)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// create new http client
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		slog.Error("open discord connection", "error", err)
		return
	}

	// register volleybot commands
	bot.RegisterCommands(dg)

	// volleyball qld client
	vqClient := vq.NewClient(vq.ClientConfig{
		Client: &httpClient,
		ApiUrl: VQClientUrl,
	})

	// create ladder changes handler
	handleLadderChanges := handleLadderChangesFactory(vqClient, myBot, dg)

	slog.Info("bot is running. press ctrl-c to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-time.NewTicker(TickSpeed).C:
			// monitor the draw pdf page
			handleLadderChanges()
		case <-sc:
			// Wait until CTRL-C or other term signal is received.
			slog.Info("removing registered commands")
			err := bot.RemoveCommands(dg)
			if err != nil {
				slog.Error("remove commands", "error", err)
			}

			slog.Info("termination signal received, bot stopping...")
			return
		}
	}

}

func handleLadderChangesFactory(vqClient *vq.Client, bot *bot.Bot, s *discordgo.Session) func() {
	currentLadder, err := vqClient.GetLadder()
	if err != nil {
		slog.Error("unable to request initial ladder data from server", "error", err)
	}

	return func() {
		slog.Info("checking for ladder changes")
		ladderUpdate, err := vqClient.GetLadder()
		if err != nil {
			slog.Error("unable to request ladder data from server", "error", err)
			return
		}

		if reflect.DeepEqual(ladderUpdate, currentLadder) {
			slog.Info("ladder monitor check no changes", "update", ladderUpdate)
			return
		}

		slog.Info("ladder monitor detected changes", "update", ladderUpdate)

		// set the current ladder to the new ladder
		currentLadder = ladderUpdate

		sb := strings.Builder{}

		sb.WriteString(fmt.Sprintf("Ladder changed: %s\n\n", "https://vqmetro23s3.softr.app/ladder-m1 ```"))

		for _, team := range ladderUpdate.Records {
			fields := team.Fields
			sb.WriteString(fmt.Sprintf("%s. %s\n", fields.Rank, team.Fields.TeamNameLookup))
			sb.WriteString(fmt.Sprintf("\tpoints: %s\n", fields.CompetitionPoints))
			sb.WriteString(fmt.Sprintf("\tnext game: %s\n", fields.NextMatch_Detail))
			sb.WriteString("\n")
		}

		sb.WriteString("```")

		fmt.Println(sb.String())

		messages, errs := bot.ChangeHandler(s, sb.String())

		// log any errors that occurred
		for _, err := range errs {
			slog.Error("ladder changes handler message failures", "error", err)
		}

		// log the messages that were sent
		for _, message := range messages {
			slog.Info("ladder changes handler message successes", "message", message)
		}
	}
}
