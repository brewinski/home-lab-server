package main

import (
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
    Token string
    lastData string
)

const (
    NotificationsChannel = "metro-volleyball-notifications"
    PING_FREQUENCY = 1 * time.Hour
    PAGE_URL = "https://www.vq.org.au/competitions/metro-league/"
)

func main() {
    // set json as the default logger
    slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

    // Create a new Discord session using the provided bot token.
    flag.StringVar(&Token, "t", "", "Bot Token")
    flag.Parse()

    // Create a new Discord session using the provided bot token.
    dg, err := discordgo.New("Bot " + Token)
    if err != nil {
        fmt.Println("error creating Discord session,", err)
        return
    }

    // register the bot ready handler
    dg.AddHandler(ready)

    // In this example, we only care about receiving message events.
    dg.Identify.Intents = discordgo.IntentsGuildMessages
    // Open a websocket connection to Discord and begin listening.
    err = dg.Open()
    if err != nil {
        fmt.Println("error opening connection,", err)
        return
    }

    go CheckPageLoop(dg)

    // Wait here until CTRL-C or other term signal is received.
    fmt.Println("Bot is now running. Press CTRL-C to exit.")
    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)
    <-sc

    // Cleanly close down the Discord session.
    dg.Close()
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
    // Set the playing status.
    // s.UpdateStatus(0, "!gopher")
    slog.Info("Metro Volleyball Bot is ready!")
    for _, guild := range event.Guilds {
        channel, err := createChannelIfNotExists(s, guild.ID, NotificationsChannel)
        if err != nil {
            slog.Error("Could not create channel", "error", err, "guild_id", guild.ID)
            continue;
        }
        s.ChannelMessageSend(channel.ID, "Metro Volleyball Bot is ready! ") // add some emojis
    }
}

func guildCreate(s *discordgo.Session, m *discordgo.GuildCreate) {
    slog.Info("guild joined to app", "guild_id", m.ID)
    // create a specific channel if it doesn't already exist
    channel, err := createChannelIfNotExists(s, m.ID, NotificationsChannel)
    if err != nil {
        slog.Error("Could not create channel", "error", err, "guild_id", m.ID)
        return;
    }

    // Send a message to that channel.
    _, err = s.ChannelMessageSend(channel.ID, "Metro Volleyball Bot as been added to server") // add some emojis
    if err != nil {
        slog.Error("Could not send message", "error", err, "guild_id", m.ID)
        return;
    }
}

func createChannelIfNotExists(s *discordgo.Session, guildId string, channelName string) (channel *discordgo.Channel, err error) {
    slog.Info("setting up channel", "guild_id", guildId)
    // Check if the channel already exists
    channels, err := s.GuildChannels(guildId)
    if err != nil {
        slog.Error("Could not get channels", "error", err, "guild_id", guildId)
        return nil, fmt.Errorf("createChannelIfNotExists: Could not get channels: %w", err)
    }

    for _, channel := range channels {
        if channel.Name == NotificationsChannel {
            slog.Info(
                "channel already exists", 
                "channel_id", channel.ID, 
                "channel_name", channel.Name, 
                "guild_id", guildId,
            )
            return channel, nil
        }
    }

    // Greate a specific channel if it doesn't already exist
    channel, err = s.GuildChannelCreate(guildId, NotificationsChannel, discordgo.ChannelTypeGuildText)
    if err != nil {
        slog.Error("Could not create channel", "error", err, "guild_id", guildId)
        return nil, fmt.Errorf("createChannelIfNotExists: Could not create channel: %w", err)
    }

    slog.Info("channel created", "channel_id", channel.ID, "channel_name", channel.Name, "guild_id", guildId)
    return channel, nil
}

func CheckPageLoop(s *discordgo.Session) {
    for range time.Tick(PING_FREQUENCY) {
        status := CheckPage()

        // if status == "" {
        //     continue;
        // }

        guilds, err := s.UserGuilds(100, "", "")
        if err != nil {
            slog.Error("Could not get guilds", "error", err)
            continue;
        }

        for _, guild := range guilds {
            channel, err := createChannelIfNotExists(s, guild.ID, NotificationsChannel)
            if err != nil {
                slog.Error("Could not create channel", "error", err, "guild_id", guild.ID)
                continue;
            }
            slog.Info("sending message", "channel_id", channel.ID, "channel_name", channel.Name, "guild_id", guild.ID)
            s.ChannelMessageSend(channel.ID, status) // add some emojis
        }
    }
}   


func CheckPage() string {
		slog.Info("checking for changes", "url", PAGE_URL)
		response, err := http.Get(PAGE_URL)
		if err != nil {
			slog.Error("api request failed", "error", err)
            return "failed request"
		} 

        defer response.Body.Close()
		
		respBytes, err := io.ReadAll(response.Body)
		if err != nil {
			slog.Error("api response read failed", "error", err)
			return "failed read"
		}

		slog.Info(
			"response", 
			"status", response.Status, 
			"content-length", response.ContentLength, 
		)

		if lastData == "" {
			slog.Info("initial data set")
			lastData = string(respBytes)
		} 
		
		if lastData != string(respBytes) {
			slog.Info("data changed")
			lastData = string(respBytes)
            return "data changed"
		} else {
			slog.Info("data unchanged")
            return "unchanged"
		}
}

