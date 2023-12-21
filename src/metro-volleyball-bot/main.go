package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/brewinski/home-lab-server/src/metro-volleyball-bot/bot"
	"github.com/brewinski/home-lab-server/src/metro-volleyball-bot/monitor"
	"github.com/brewinski/home-lab-server/src/metro-volleyball-bot/vq"
	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token     string
	TickSpeed time.Duration
	PageUrl   string
)

const (
	NotificationsChannel = "metro-volleyball-notifications"
)

func main() {
	// Read the Discord token as an application flag.
	flag.StringVar(&Token, "t", "", "The token for the specific discord application.")
	// Read the page check frequency duration. Parsed as "1ms", "1ns", "1s", "1m", or "1h"
	flag.DurationVar(&TickSpeed, "ts", 1*time.Hour, "Page ping frequency as a string duration")
	// Page URL to monitor
	flag.StringVar(&PageUrl, "url", PageUrl, "The URL to monitor for changes")
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

	// create monitor deps and start monitoring
	// Draw Pdf monitor
	pdfMonitor := monitor.New(monitor.Config{
		Client: &httpClient,
	})
	// make initial request for pdf data
	_, _, err = pdfMonitor.CheckForChanges(PageUrl)
	if err != nil {
		slog.Error("pdf initial data request failed", "error", err)
	}

	// volleyball qld client
	vqClient := vq.NewClient(vq.ClientConfig{
		Client: &httpClient,
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
			handlePDFChanges(pdfMonitor, myBot, dg)
			handleLadderChanges()
		case <-sc:
			// Wait until CTRL-C or other term signal is received.
			slog.Info("termination signal received, bot stopping...")
			return
		}
	}

}

func handlePDFChanges(m *monitor.DataSourceMonitor, bot *bot.Bot, s *discordgo.Session) {
	pdfUpdate, _, err := m.CheckForChanges(PageUrl)
	if err != nil {
		slog.Error("pdf monitor check for changes", "error", err)
	}

	if !pdfUpdate {
		slog.Info("pdf monitor check no changes", "update", pdfUpdate)
		return
	}

	slog.Info("pdf monitor check for changes", "update", pdfUpdate)

	messages, errs := bot.ChangeHandler(s, fmt.Sprintf("PDF changed, go to %s and review the changes", PageUrl))

	// log any errors that occurred
	for _, err := range errs {
		slog.Error("pdf changes handler message failures", "error", err)
	}

	// log the messages that were sent
	for _, message := range messages {
		slog.Info("pdf changes handler message successes", "message", message)
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

		if LadderHasChanged(ladderUpdate, currentLadder) {
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

func LadderHasChanged(new, old vq.GetLadderResponseBody) bool {
	newRecords, oldRecords := new.Records, old.Records

	if len(new.Records) != len(old.Records) {
		return true
	}

	for i, newRecord := range newRecords {
		oldRecord := oldRecords[i]

		if newRecord.Fields.Rank != oldRecord.Fields.Rank {
			slog.Info("rank changed", "new", newRecord.Fields.Rank, "old", oldRecord.Fields.Rank)
			return true
		}

		if newRecord.Fields.TeamNameLookup != oldRecord.Fields.TeamNameLookup {
			slog.Info("team name changed", "new", newRecord.Fields.TeamNameLookup, "old", oldRecord.Fields.TeamNameLookup)
			return true
		}

		if newRecord.Fields.CompetitionPoints != oldRecord.Fields.CompetitionPoints {
			slog.Info("competition points changed", "new", newRecord.Fields.CompetitionPoints, "old", oldRecord.Fields.CompetitionPoints)
			return true
		}
	}

	return false
}
