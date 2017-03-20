package main

/*
	Imports
*/

import (
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

/*
	Constants
*/

const (
	discordServerID       = "83214009964171264"  // This is the ID of the "Isaac Speedrunning & Racing" server
	discordLobbyChannelID = "286115994621968384" // This is the ID of the "racing-plus-lobby" text channel
)

/*
	Global variables
*/

var (
	discord      *discordgo.Session
	discordBotID string
)

func discordInit() {
	if useDiscord == false {
		return
	}

	// Local variables
	var err error

	// Read the OAuth secret from the environment variable (it was loaded from the .env file in main.go)
	oauthToken := os.Getenv("DISCORD_OAUTH")

	// Connect
	discord, err = discordgo.New("Bot " + oauthToken) // Bot accounts must be prefixed with "Bot"
	if err != nil {
		log.Error("Error creating Discord session: ", err)
		return
	}

	// Register ready as a callback for the ready events
	discord.AddHandler(discordReady)

	// Register messageCreate as a callback for the messageCreate events
	discord.AddHandler(discordMessageCreate)

	// Open the websocket and begin listening.
	err = discord.Open()
	if err != nil {
		log.Error("Error opening Discord session: ", err)
	}
}

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

		discordSend(m.ChannelID, "Go to https://isaacracing.net/info and Ctrl+F for: `What do I do if the mod doesn't seem to be working correctly in-game?`")
	}

	// Copy messages from "racing-plus-lobby"
	if m.ChannelID == discordLobbyChannelID {
		// Send everyone the notification
		connectionMap.RLock()
		for _, conn := range connectionMap.m {
			conn.Connection.Emit("discordMessage", &RoomMessageMessage{
				Name:    m.Author.Username + "#" + m.Author.Discriminator,
				Message: message,
			})
		}
		connectionMap.RUnlock()
	}
}

func discordSend(channelID string, message string) {
	if useDiscord == false {
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
