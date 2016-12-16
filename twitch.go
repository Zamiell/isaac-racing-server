package main

/*
	Imports
*/

import (
	"github.com/Zamiell/isaac-racing-server/models"

	"bufio"
	"net"
	"net/textproto"
	"os"
	"regexp"
	"strings"
	"time"
)

/*
	Global variables
*/

var (
	twitchUsername string = "IsaacRacingPlus"
	conn           net.Conn
	twitchModMap   = make(map[string]bool)
)

func twitchInit() {
	// Read the OAuth secret from the environment variable (it was loaded from the .env file in main.go)
	oauthToken := os.Getenv("TWITCH_OAUTH")

	// Connect to the Twitch IRC server
	var err error
	conn, err = net.Dial("tcp", "irc.chat.twitch.tv:6667")
	if err != nil {
		log.Error("Failed to connect to the Twitch IRC:", err)
		time.Sleep(60 * time.Second)
		go twitchInit()
		return
	}
	defer conn.Close()

	// Send our Twitch credentials (the pass has to come before the nick)
	ircSend("PASS " + oauthToken)
	ircSend("NICK " + twitchUsername)

	// Request Twitch specific capabilities
	// This is required to see who is a moderator in the channel: https://github.com/justintv/Twitch-API/blob/master/IRC.md
	ircSend("CAP REQ :twitch.tv/membership")
	ircSend("CAP REQ :twitch.tv/commands")
	ircSend("CAP REQ :twitch.tv/tags")

	// Figure out which channels to join
	streams, err := db.Users.GetAllStreams()
	if err != nil {
		log.Error("Database error:", err)
		return
	}

	// Join all the channels
	for _, stream := range streams {
		// Just in case, skip over streams that are not Twitch URLs
		if strings.HasPrefix(stream, "https://www.twitch.tv/") {
			// Parse for the username
			re, err := regexp.Compile(`https://www.twitch.tv/(.+)`)
			if err != nil {
				log.Error("Failed to compile the Twitch username regular expression:", err)
				return
			}
			user := re.FindStringSubmatch(stream)[1]
			user = strings.ToLower(user)

			// Join it
			ircSend("JOIN #" + user)
		} else {
			log.Error("A user had a stream set to \"" + stream + "\" but their \"twitch_bot_enabled\" was set to 1, which should never happen.")
		}
	}

	// Listen for IRC commands
	tp := textproto.NewReader(bufio.NewReader(conn))
	for {
		// Block until we get a message
		msg, err := tp.ReadLine()
		if err != nil {
			log.Error("Failed to read from the Twitch IRC connection:", err)
			go twitchInit()
			return
		}

		// Log all messages
		log.Info("< " + msg)

		// Split the message by spaces
		msgParts := strings.Split(msg, " ")

		// We have to respond to all PINGS or else we will get kicked
		if msgParts[0] == "PING" {
			ircSend("PONG " + msgParts[1])
			continue
		}

		// Avoid potential index errors
		if len(msgParts) < 4 {
			continue
		}
		metadata := msgParts[0]
		command := msgParts[2]
		channel := msgParts[3][1:] // Remove the # at the beginning

		if command == "USERSTATE" {
			// Listen to see if we are a mod in the channel
			if strings.Contains(metadata, ";display-name=IsaacRacingPlus;") {
				if strings.Contains(metadata, "mod=0") {
					twitchModMap[channel] = false
				} else if strings.Contains(metadata, "mod=1") {
					twitchModMap[channel] = true
				}
			}
			continue
		}

		if command == "PRIVMSG" {
			// Remove the colon at the beginning of the message
			message := msgParts[4][1:]

			// Information commands
			if message == "!racingplus" || message == "!racing+" {
				twitchSend(channel, "Racing+ is a mod for The Binding of Isaac: Afterbirth+: https://isaacracing.net", 0)
			}
		}
	}
}

func ircSend(command string) {
	log.Info("> " + command)
	conn.Write([]byte(command + "\r\n"))
}

func twitchRacerSend(racer models.Racer, message string) {
	if racer.TwitchBotEnabled == 1 {
		if strings.HasPrefix(racer.Stream, "https://www.twitch.tv/") {
			// Parse for the username
			re, err := regexp.Compile(`https://www.twitch.tv/(.+)`)
			if err != nil {
				commandMutex.Unlock()
				log.Error("Failed to compile the Twitch username regular expression:", err)
				return
			}
			channel := re.FindStringSubmatch(racer.Stream)[1]
			channel = strings.ToLower(channel)

			twitchSend(channel, message, racer.TwitchBotDelay)
		} else {
			log.Error("User \"" + racer.Name + "\" had a stream set to \"" + racer.Stream + "\" but their \"twitch_bot_enabled\" was set to 1, which should never happen.")
		}
	}
}

func twitchSend(channel string, message string, delay int) {
	// Don't talk in any channels that we are not a moderator in
	// (we will get banned for spamming otherwise)
	if twitchModMap[channel] == false {
		return
	}

	go twitchSend2(channel, message, delay)
}

func twitchSend2(channel string, message string, delay int) {
	time.Sleep(time.Second * time.Duration(delay))
	ircSend(":" + twitchUsername + "!" + twitchUsername + "@" + twitchUsername + ".tmi.twitch.tv PRIVMSG #" + channel + " :" + message)
}
