package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/brewinski/home-lab-server/src/metro-volleyball-bot/bot"
	"github.com/brewinski/home-lab-server/src/metro-volleyball-bot/monitor"
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
	// Read the Discord token as an applicatioln flag.
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

	// create new http client
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	// register the bot ready handler
	dg.AddHandler(myBot.ReadyHandler)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		slog.Error("open discord connection", "error", err)
		return
	}

	// Draw Pdf monitor
	pdfMonitor := monitor.New(monitor.Config{
		Client: &httpClient,
	})
	// make initial request for pdf data
	_, _, err = pdfMonitor.CheckForChanges(PageUrl)
	if err != nil {
		slog.Error("pdf initial data request failed", "error", err)
	}

	slog.Info("bot is running. press ctrl-c to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-time.NewTicker(TickSpeed).C:
			// monitor the draw pdf page
			handlePDFChanges(pdfMonitor, myBot, dg)
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

	messages, errs := bot.PdfChangesHandler(s, fmt.Sprintf("PDF changed, go to %s and review the changes", PageUrl))

	// log any errors that occurred
	for _, err := range errs {
		slog.Error("pdf changes handler message failures", "error", err)
	}

	// log the messages that were sent
	for _, message := range messages {
		slog.Info("pdf changes handler message successes", "message", message)
	}
}
