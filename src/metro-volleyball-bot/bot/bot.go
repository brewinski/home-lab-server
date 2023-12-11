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
	lastPageResponse string
	config Config
}

type Config struct {
	MonitorUrl string
	UpdatesChannel string
	TickSpeed time.Duration
}

func New(cfg Config) *Bot {
	return &Bot{
		"",
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

        s.ChannelMessageSend(channel.ID, fmt.Sprintf("metro bot ready, monitoring page: %s", b.config.MonitorUrl) ) // add some emojis
    }
}

// Monitor a specific page for changes
func (b *Bot) MoitorListenAndServe(s *discordgo.Session) {
	// if the last response is empty we need to get the initial response.
	_, err := b.monitorPage(b.config.MonitorUrl)
	if err != nil {
		slog.Error("request to monitor page failed for seeding initial data", "error", err)
		return
	}

	// loop for the duration of the program.
    for range time.Tick(b.config.TickSpeed) {
        status, err := b.monitorPage(b.config.MonitorUrl)
		if err != nil {
			slog.Error("monitor request failed", "url", b.config.MonitorUrl, "error", err)
		}

        if !status {
            slog.Info("monitor request, nothing has changed", b.config.MonitorUrl, slog.Bool("status", status))
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

            slog.Info("message sending", "channel_id", channel.ID, "channel_name", channel.Name, "guild_id", guild.ID)

            s.ChannelMessageSend(
				channel.ID,
				fmt.Sprintf("monitored page has changed, go to %s and review the changes", b.config.MonitorUrl),
			)

			slog.Info("message sent", "channel_id", channel.ID, "channel_name", channel.Name, "guild_id", guild.ID)
        }
    }
}

// ping the page and check if anything has changed.
// it will return false if the page hasn't been checked in the past
func (b *Bot) monitorPage(pageUrl string) (bool, error) {
	response, err := http.Get(pageUrl)
	if err != nil {
		return false, err
	}

	defer response.Body.Close()
	
	respBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return false, err
	}

	slog.Info("page response", "status", response.Status, "content-length", response.ContentLength)
	
	if b.lastPageResponse != string(respBytes) {
		b.lastPageResponse = string(respBytes)
		return true, nil
	} else {
		return false, nil
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

