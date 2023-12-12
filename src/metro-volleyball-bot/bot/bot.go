package bot

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	lastPageResponse string
	config           Config
}

type Config struct {
	MonitorUrl     string
	UpdatesChannel string
	TickSpeed      time.Duration
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
			continue
		}

		s.ChannelMessageSend(channel.ID, fmt.Sprintf("metro bot ready, monitoring page: %s", b.config.MonitorUrl)) // add some emojis
	}
}

// Monitor a specific page for changes
func (b *Bot) MonitorListenAndServe(s *discordgo.Session) {
	// if the last response is empty we need to get the initial response.
	initResp, err := MonitorPage(b.config.MonitorUrl)
	if err != nil {
		slog.Error("request to monitor page failed for seeding initial data", "error", err)
		return
	}
	// set the last response to the initial response
	b.lastPageResponse = initResp

	// loop for the duration of the program.
	for range time.NewTicker(b.config.TickSpeed).C {
		resp, err := MonitorPage(b.config.MonitorUrl)
		if err != nil {
			slog.Error("monitor request failed", "url", b.config.MonitorUrl, "error", err)
		}

		if resp == b.lastPageResponse {
			slog.Info("monitor request, nothing has changed", b.config.MonitorUrl, "response", resp)
			continue
		}

		// if the response is different than the last response, set the last response to the new response.
		b.lastPageResponse = resp

		// if the response is the same as the last response, send a message to the updates channel.
		guilds, err := s.UserGuilds(100, "", "")
		if err != nil {
			slog.Error("Could not get guilds", "error", err)
			continue
		}

		for _, guild := range guilds {
			channel, err := createChannelIfNotExists(s, guild.ID, b.config.UpdatesChannel)
			if err != nil {
				slog.Error("Could not create channel", "error", err, "guild_id", guild.ID)
				continue
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
			slog.Info("channel exists, skipping create", "channel_id", channel.ID, "channel_name", channel.Name, "guild_id", guildId)
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
