package bot

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/brewinski/home-lab-server/src/metro-volleyball-bot/flags"
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

		// send a ready message if the feature flag is enabled
		if flags.BotReadyMessage {
			s.ChannelMessageSend(channel.ID, fmt.Sprintf("metro bot ready, monitoring page: %s", b.config.MonitorUrl)) // add some emojis
		}
	}
}

func (b *Bot) ChangeHandler(s *discordgo.Session, message string) ([]*discordgo.Message, []error) {
	// Get a list of all the guilds that are available for messages
	guilds, err := s.UserGuilds(100, "", "")
	if err != nil {
		return nil, []error{fmt.Errorf("unable to list guilds: %w", err)}
	}

	var errors []error
	var messages []*discordgo.Message

	// Send a message to each guild
	for _, guild := range guilds {
		channel, err := createChannelIfNotExists(s, guild.ID, b.config.UpdatesChannel)
		if err != nil {
			errors = append(errors, fmt.Errorf("unable to create channel: %w, guild[%s]", err, guild.Name))
			continue
		}

		message, err := s.ChannelMessageSend(channel.ID, message)
		if err != nil {
			// if the message fails to send, add the error to the list of errors and continue to the next guild
			errors = append(errors, fmt.Errorf("unable to send message: %w, guild[%s], channel[%s]", err, guild.Name, channel.Name))
			continue
		}

		// add the successful message to the list of messages to be returned
		messages = append(messages, message)
	}

	return messages, errors
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
