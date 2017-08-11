package main

/*
	Imports
*/

import (
	"bufio"
	"net"
	"net/textproto"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/Zamiell/isaac-racing-server/src/log"
	"github.com/Zamiell/isaac-racing-server/src/models"
)

/*
	Global variables
*/

var (
	twitchUsername string = "IsaacRacingPlus"
	twitchConn     net.Conn
)

/*
	Initialization function
*/

func twitchInit() {
	if useTwitch {
		go twitchConnect()
	}
}

func twitchConnect() {
	// Read the OAuth secret from the environment variable (it was loaded from the .env file in main.go)
	oauthToken := os.Getenv("TWITCH_OAUTH")

	// Connect to the Twitch IRC server
	var err error
	twitchConn, err = net.Dial("tcp", "irc.chat.twitch.tv:6667")
	if err != nil {
		log.Error("Failed to connect to the Twitch IRC:", err)
		time.Sleep(30 * time.Second)
		go twitchConnect() // Reconnect after 30 seconds
		return
	}
	defer twitchConn.Close()

	// Send our Twitch credentials (the pass has to come before the nick)
	twitchIRCSend("PASS " + oauthToken)
	twitchIRCSend("NICK " + twitchUsername)

	// Request Twitch specific capabilities
	// (this is required to see who is a moderator in the channel)
	// https://github.com/justintv/Twitch-API/blob/master/IRC.md
	twitchIRCSend("CAP REQ :twitch.tv/membership")
	twitchIRCSend("CAP REQ :twitch.tv/commands")
	twitchIRCSend("CAP REQ :twitch.tv/tags")

	// Figure out which channels to join
	streamURLs, err := db.Users.GetAllStreamURLs()
	if err != nil {
		log.Error("Database error:", err)
		return
	}

	// Join all the channels
	for _, streamURL := range streamURLs {
		// Just in case, ensure that this is a Twitch URL
		if !strings.HasPrefix(streamURL, "https://www.twitch.tv/") {
			log.Error("A user had a stream URL set to \"" + streamURL + "\" but their \"twitch_bot_enabled\" was set to 1, which should never happen.")
			continue
		}

		// Parse for the username
		re, err := regexp.Compile(`https://www.twitch.tv/(.+)`)
		if err != nil {
			log.Error("Failed to compile the Twitch username regular expression:", err)
			return
		}
		user := re.FindStringSubmatch(streamURL)[1]
		user = strings.ToLower(user)

		// Join it
		twitchJoinChannel(user)
	}

	// Listen for IRC commands
	tp := textproto.NewReader(bufio.NewReader(twitchConn))
	for {
		// Block until we get a message
		msg, err := tp.ReadLine()
		if err != nil {
			// Ocassionally the connection is reset, so don't log this as an error
			log.Info("Failed to read from the Twitch IRC connection:", err)
			go twitchInit() // Reconnect
			return
		}

		// Log all messages
		log.Info("< " + msg)

		// Split the message by spaces
		msgParts := strings.Split(msg, " ")

		// We have to respond to all PINGS or else we will get kicked
		if msgParts[0] == "PING" {
			twitchIRCSend("PONG " + msgParts[1])
			continue
		}

		// Listen to see if we are a mod in the channel (1/2)
		// (this is emitted after joining a channel and after someone is modded or demodded)
		// MODE commands don't follow the same format as other commands
		// e.g. :jtv MODE #zamiell +o isaacracingplus
		if msgParts[0] == ":jtv" && msgParts[1] == "MODE" {
			channel := msgParts[2][1:] // Remove the # at the beginning
			if msgParts[2] == "-o" {
				twitchNotMod(channel)
			}
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
			// Listen to see if we are a mod in the channel (2/2)
			// (this is emitted after joining a channel and after we talk in a channel)
			if strings.Contains(metadata, ";display-name=IsaacRacingPlus;") {
				if strings.Contains(metadata, ";mod=0;") {
					twitchNotMod(channel)
				}
			}
			continue

		} else if command == "PRIVMSG" {
			// Remove the colon at the beginning of the message and make it lowercase for easier parsing
			message := msgParts[4][1:]
			message = strings.ToLower(message)

			// User commands
			if message == "!racingplus" || message == "!racing+" || message == "!r+" {
				twitchSend(channel, "Racing+ is a mod for The Binding of Isaac: Afterbirth+: https://isaacracing.net", 0)
			} else if message == "!left" {
				// TODO
			} else if message == "!entrants" {
				// TODO
			}
		}
	}
}

/*
	Mod functions
*/

func twitchNotMod(channel string) {
	log.Info("Detected that we are not a mod in the channel of \"" + channel + "\".")

	twitchLeaveChannel(channel)

	// Get the user ID and username of the user that matches this stream
	streamURL := "https://www.twitch.tv/" + channel
	userID, username, err := db.Users.FindUserFromStreamURL(streamURL)
	if err != nil {
		log.Error("Database error:", err)
		return
	} else if userID == 0 {
		log.Error("Was not able to find the user ID that goes along with the stream URL of: " + streamURL)
		return
	}

	// Disable the Twitch bot in the database
	if err := db.Users.SetTwitchBotEnabled(userID, false); err != nil {
		log.Error("Database error:", err)
		return
	}

	// If the user is online, send them a warning to let them know what happened
	if s, ok := websocketSessions[username]; ok {
		websocketWarning(s, "twitch.notMod", "The Twitch bot has been disabled for your account because it is not a moderator in your channel. Please type <code>/mod IsaacRacingPlus</code> in your Twitch chat, wait a few minutes, and then check the box again in the Racing+ settings.")
	}
}

// Called from the "raceFinish", "raceQuit", and the "raceCheckStart" functions
func twitchRacerSend(racer models.Racer, message string) {
	if !useTwitch {
		return
	}

	if racer.TwitchBotEnabled == 0 {
		return
	}

	if !strings.HasPrefix(racer.StreamURL, "https://www.twitch.tv/") {
		log.Error("User \"" + racer.Name + "\" had a stream URL set to \"" + racer.StreamURL + "\" but their \"TwitchBotEnabled\" was set to 1, which should never happen.")
		return
	}

	// Parse for the username
	re, err := regexp.Compile(`https://www.twitch.tv/(.+)`)
	if err != nil {
		log.Error("Failed to compile the Twitch username regular expression:", err)
		return
	}
	channel := re.FindStringSubmatch(racer.StreamURL)[1]
	channel = strings.ToLower(channel)

	twitchSend(channel, message, racer.TwitchBotDelay)
}

func twitchSend(channel string, message string, delay int) {
	go func(channel string, message string, delay int) {
		time.Sleep(time.Second * time.Duration(delay))
		twitchIRCSend(":" + twitchUsername + "!" + twitchUsername + "@" + twitchUsername + ".tmi.twitch.tv PRIVMSG #" + channel + " :" + message)
	}(channel, message, delay)
}

/*
	Miscellaneous functions
*/

func twitchJoinChannel(channel string) {
	twitchIRCSend("JOIN #" + channel)
}

func twitchLeaveChannel(channel string) {
	twitchIRCSend("PART #" + channel)
}

func twitchIRCSend(command string) {
	log.Info("> " + command)
	twitchConn.Write([]byte(command + "\r\n"))
}
