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
    // Program Information
    Token string
)

const (
    NotificationsChannel = "metro-volleyball-notifications"
    PING_FREQUENCY = 1 * time.Hour
    PAGE_URL = "https://www.vq.org.au/competitions/metro-league/"
)

func main() {
    // set json as the default logger
    slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(
        slog.Duration("tick_speed", PING_FREQUENCY),
        slog.String("channel", NotificationsChannel),
        slog.String("page", PAGE_URL),
    ))

    // Create a new Discord session using the provided bot token.
    flag.StringVar(&Token, "t", "", "Bot Token")
    flag.Parse()

    // Create a new Discord session using the provided bot token.
    dg, err := discordgo.New("Bot " + Token)
    if err != nil {
        slog.Error("new discord session", "error", err)
        return
    }

    // create a new bot implementation
    myBot := bot.New(bot.Config{
        UpdatesChannel: NotificationsChannel,
    })

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

