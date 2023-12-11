package bot

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	config Config
}

type Config struct {
	MonitorUrl string
	UpdatesChannel string
	TickSpeed time.Duration
}

func New(cfg Config) *Bot {
	return &Bot{
		cfg,
	}
}

// Ready handler will be called when the session is ready.
func (b *Bot) ReadyHandler(s *discordgo.Session, event *discordgo.Ready) {
    // Set the playing status.
    slog.Info("metro volleyball bot ready.")

    for _, guild := range event.Guilds {
        channel, err := createChannelIfNotExists(s, guild.ID, b.config.UpdatesChannel)
        if err != nil {
            slog.Error("Could not create channel", "error", err, "guild_id", guild.ID)
            continue;
        }

        s.ChannelMessageSend(channel.ID, "metro bot ready") // add some emojis
    }
}

// Monitor a specific page for changes
func (b *Bot) MonitorPageHandler(s *discordgo.Session) {
    for range time.Tick(b.config.TickSpeed) {
        status := monitorPage(b.config.MonitorUrl)

        if status == "" {
            slog.Info("nothin has changed")
            continue;
        }

        guilds, err := s.UserGuilds(100, "", "")
        if err != nil {
            slog.Error("Could not get guilds", "error", err)
            continue;
        }

        for _, guild := range guilds {
            channel, err := createChannelIfNotExists(s, guild.ID, b.config.UpdatesChannel)
            if err != nil {
                slog.Error("Could not create channel", "error", err, "guild_id", guild.ID)
                continue;
            }
            slog.Info("sending message", "channel_id", channel.ID, "channel_name", channel.Name, "guild_id", guild.ID)
            s.ChannelMessageSend(channel.ID, status) // add some emojis
        }
    }
}

var lastData string

func monitorPage(pageUrl string) string {
	slog.Info("monitoring page for changes", "url", pageUrl)

	response, err := http.Get(pageUrl)
	if err != nil {
		slog.Error("page request failed", "error", err)
		return "failed request"
	}

	defer response.Body.Close()
	
	respBytes, err := io.ReadAll(response.Body)
	if err != nil {
		slog.Error("page response unreadable", "error", err)
		return "failed read"
	}

	slog.Info("page response", "status", response.Status, "content-length", response.ContentLength)

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

// creates the updates channel if it doesn't already exist
func createChannelIfNotExists(s *discordgo.Session, guildId string, channelName string) (channel *discordgo.Channel, err error) {
    // Check if the channel already exists
    channels, err := s.GuildChannels(guildId)
    if err != nil {
        slog.Error("reading channels", "error", err, "guild_id", guildId)
        return nil, fmt.Errorf("createChannelIfNotExists: Could not get channels: %w", err)
    }

    // get the existing channel if it already exists
    for _, channel := range channels {
        if channel.Name == channelName {
            slog.Info("channel exists, skipping create", "channel_id", channel.ID, "channel_name", channel.Name, "guild_id", guildId,)
            return channel, nil
        }
    }

    slog.Info("channel not found, creating channel...", "guild_id", guildId)

    // create the channel if it doesn't already exist.
    channel, err = s.GuildChannelCreate(guildId, channelName, discordgo.ChannelTypeGuildText)
    if err != nil {
        slog.Error("create channel", "error", err, "guild_id", guildId)
        return nil, fmt.Errorf("createChannelIfNotExists: Could not create channel: %w", err)
    }

    slog.Info("channel created", "channel_id", channel.ID, "channel_name", channel.Name, "guild_id", guildId)
    return channel, nil
}

