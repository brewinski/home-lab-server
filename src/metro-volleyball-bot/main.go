package main

import (
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/brewinski/home-lab-server/src/metro-volleyball-bot/bot"
	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
    Token string
    TickSpeed time.Duration
)

const (
    NotificationsChannel = "metro-volleyball-notifications"
    PageUrl = "https://www.vq.org.au/competitions/metro-league/"
)

func main() {
    // Read the Discord token as an applicatioln flag.
    flag.StringVar(&Token, "t", "", "The token for the specific discord application.")
    // Read the page check frequency duration. Parsed as "1ms", "1ns", "1s", "1m", or "1h"
    flag.DurationVar(&TickSpeed, "ts", 1 * time.Hour, "Page ping frequency as a string duration")
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

    myBot := bot.New(bot.Config{
        UpdatesChannel: NotificationsChannel,
        TickSpeed: TickSpeed,
        MonitorUrl: PageUrl,
    })

    // get initial page data
    // myBot.

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

    // monitor the page 
    go myBot.MonitorPageHandler(dg)

    // Wait here until CTRL-C or other term signal is received.
    slog.Info("bot is running. press ctrl-c to exit.")

    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)
    <-sc
    // Cleanly close down the Discord session.
    dg.Close()

    // notify that a 
    slog.Info("termination signal recieved, bot stopping...")
}

