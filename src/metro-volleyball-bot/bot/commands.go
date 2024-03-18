package bot

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name: "vb-help",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: "help documentation.",
		},
		{
			Name:        "vb-ladder",
			Description: "view the latest ladder results.",
			Version:     "1.0.0",
		},
	}
)

// RegisterCommands registers all commands for the bot.
func RegisterCommands(s *discordgo.Session) error {
	// Get a list of all the guilds that are available for messages
	guilds, err := s.UserGuilds(100, "", "")
	if err != nil {
		return err
	}

	// register commands for each guild
	for _, guild := range guilds {
		for _, command := range commands {
			_, err := s.ApplicationCommandCreate(s.State.User.ID, guild.ID, command)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// RemoveCommands deregisters all commands for the bot.
func RemoveCommands(s *discordgo.Session) error {
	// Get a list of all the guilds that are available for messages
	guilds, err := s.UserGuilds(100, "", "")
	if err != nil {
		return err
	}

	// register commands for each guild
	for _, guild := range guilds {
		guildCommands, err := s.ApplicationCommands(s.State.User.ID, guild.ID)
		if err != nil {
			return err
		}

		for _, command := range guildCommands {
			slog.Info("removing command", "command", command.Name)
			err := s.ApplicationCommandDelete(s.State.User.ID, guild.ID, command.ID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// OnCommandHandler handles all commands for the bot. and allows the user to register
// a callback that returns the string to send to the channel.
func OnCommandHandlerFactory(callback func(string) (string, error)) func(*discordgo.Session, *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		slog.Info("handling command", "command", i.ApplicationCommandData().Name)

		message, err := callback(i.ApplicationCommandData().Name)
		if err != nil {
			slog.Error("error handling command", "error", err)
			responseWithMessage(s, i, "something went wrong, please try again later")
			return
		}

		responseWithMessage(s, i, message)
	}
}

// responseWithMessage handles the help command.
func responseWithMessage(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	// Create the response object
	response := discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	}
	slog.Info("responding to interaction", "response", response)
	// Respond to the interaction with the response object
	s.InteractionRespond(i.Interaction, &response)
}
