package main

import (
	"os"
	"strings"

	"github.com/Zamiell/isaac-racing-server/src/log"
	"github.com/bwmarrin/discordgo"
)

const (
	// This is the ID of the "Isaac Speedrunning & Racing" server (guild)
	discordServerID = "83214009964171264"

	// This is the ID of the "racing-plus-lobby" text channel
	discordLobbyChannelID = "286115994621968384"
)

var (
	discord        *discordgo.Session
	discordEnabled = false
	discordBotID   string
)

/*
	Initialization functions
*/

func discordInit() {
	// Read the OAuth secret from the environment variable
	// (it was loaded from the .env file in main.go)
	token := os.Getenv("DISCORD_TOKEN")
	if len(token) == 0 {
		log.Info("The \"DISCORD_TOKEN\" environment variable is blank; aborting Discord bot initialization.")
		return
	}

	discordEnabled = true
	go discordConnect(token)
}

func discordConnect(token string) {
	// Bot accounts must be prefixed with "Bot"
	if d, err := discordgo.New("Bot " + token); err != nil {
		log.Error("Error creating Discord session: ", err)
		return
	} else {
		discord = d
	}

	// Register function handlers for various events
	discord.AddHandler(discordReady)
	discord.AddHandler(discordMessageCreate)

	// Open the websocket and begin listening
	if err := discord.Open(); err != nil {
		log.Error("Error opening Discord session: ", err)
	}
}

/*
	Event handlers
*/

func discordReady(s *discordgo.Session, event *discordgo.Ready) {
	log.Info("Discord bot connected with username: " + event.User.Username)
	discordBotID = event.User.ID
}

func discordMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == discordBotID {
		return
	}

	// Log the message
	log.Info("[D" + m.ChannelID + "] <" + m.Author.Username + "#" + m.Author.Discriminator + "> " + m.Content)

	// Info commands
	message := strings.ToLower(m.Content)
	if message == "!r+" ||
		message == "!racing+" ||
		message == "!racingplus" {

		discordSend(m.ChannelID, "Racing+ is a mod for The Binding of Isaac: Afterbirth+: https://isaacracing.net")
		return
	} else if message == "!corrupt" ||
		message == "!corrupted" ||
		message == "!corruptmod" ||
		message == "!corruptedmod" {

		discordSend(m.ChannelID, "`What do I do if the mod doesn't seem to be working correctly in-game?`\n<https://isaacracing.net/info#faq-corrupt>")
	} else if message == "!doc" ||
		message == "!documentation" {

		discordSend(m.ChannelID, "Everything in the mod has detailed documentation if you just bother to look on the website! Here's a handy link for you: https://github.com/Zamiell/isaac-racing-client/blob/master/mod/CHANGES.md")
	}

	// Copy messages from "racing-plus-lobby"
	if m.ChannelID == discordLobbyChannelID {
		// Send everyone the notification
		commandMutex.Lock()
		type discordMessageMessage struct {
			Name    string `json:"name"`
			Message string `json:"message"`
		}
		for _, s := range websocketSessions {
			websocketEmit(s, "discordMessage", &discordMessageMessage{
				Name:    m.Author.Username + "#" + m.Author.Discriminator,
				Message: message,
			})
		}
		commandMutex.Unlock()
	}
}

/*
	Miscellaneous functions
*/

func discordSend(channelID string, message string) {
	if !discordEnabled {
		return
	}

	_, err := discord.ChannelMessageSend(channelID, message)
	if err != nil {
		errorMessage := "Failed to send message to \""
		if channelID == discordLobbyChannelID {
			errorMessage += "racing-plus-lobby"
		} else {
			errorMessage += channelID
		}
		errorMessage += "\": " + message
		log.Warning(errorMessage)
	}
}
