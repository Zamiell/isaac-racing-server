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
	"sync"
	"time"
)

/*
	Global variables
*/

var (
	twitchUsername string = "IsaacRacingPlus"
	conn           net.Conn
	twitchModMap   = struct {
		// Maps are not safe for concurrent use: https://blog.golang.org/go-maps-in-action
		sync.RWMutex
		m map[string]bool
	}{m: make(map[string]bool)}
	twitchModCheckingMap = struct {
		// Maps are not safe for concurrent use: https://blog.golang.org/go-maps-in-action
		sync.RWMutex
		m map[string]bool
	}{m: make(map[string]bool)}
)

func twitchInit() {
	// Read the OAuth secret from the environment variable (it was loaded from the .env file in main.go)
	oauthToken := os.Getenv("TWITCH_OAUTH")

	// Connect to the Twitch IRC server
	var err error
	conn, err = net.Dial("tcp", "irc.chat.twitch.tv:6667")
	if err != nil {
		log.Error("Failed to connect to the Twitch IRC:", err)
		time.Sleep(30 * time.Second)
		go twitchInit() // Reconnect after 30 seconds
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
	streamURLs, err := db.Users.GetAllStreamURLs()
	if err != nil {
		log.Error("Database error:", err)
		return
	}

	// Join all the channels
	for _, streamURL := range streamURLs {
		// Just in case, skip over streams that are not Twitch URLs
		if strings.HasPrefix(streamURL, "https://www.twitch.tv/") {
			// Parse for the username
			re, err := regexp.Compile(`https://www.twitch.tv/(.+)`)
			if err != nil {
				log.Error("Failed to compile the Twitch username regular expression:", err)
				return
			}
			user := re.FindStringSubmatch(streamURL)[1]
			user = strings.ToLower(user)

			// Join it
			ircSend("JOIN #" + user)
		} else {
			log.Error("A user had a stream URL set to \"" + streamURL + "\" but their \"twitch_bot_enabled\" was set to 1, which should never happen.")
		}
	}

	// Listen for IRC commands
	tp := textproto.NewReader(bufio.NewReader(conn))
	for {
		// Block until we get a message
		msg, err := tp.ReadLine()
		if err != nil {
			log.Error("Failed to read from the Twitch IRC connection:", err)
			go twitchInit() // Reconnect
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

		// Listen to see if we are a mod in the channel (1/2)
		// (this is emitted after joining a channel and after someone is modded or demodded)
		// MODE commands don't follow the same format as other commands
		// e.g. :jtv MODE #zamiell +o isaacracingplus
		if msgParts[0] == ":jtv" && msgParts[1] == "MODE" {
			channel := msgParts[2][1:] // Remove the # at the beginning
			twitchModMap.Lock()
			if msgParts[2] == "+o" {
				twitchModMap.m[channel] = true
			} else if msgParts[2] == "-o" {
				twitchModMap.m[channel] = false

				// Check again for mod after a delay
				go twitchCheckForMod(channel)
			}
			twitchModMap.Unlock()
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
				twitchModMap.Lock()
				if strings.Contains(metadata, ";mod=0;") {
					twitchModMap.m[channel] = false

					// Check again for mod after a delay
					go twitchCheckForMod(channel)
				} else if strings.Contains(metadata, ";mod=1;") {
					twitchModMap.m[channel] = true
				}
				twitchModMap.Unlock()
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

func ircSend(command string) {
	log.Info("> " + command)
	conn.Write([]byte(command + "\r\n"))
}

func twitchRacerSend(racer models.Racer, message string) {
	// We don't have to worry about the command mutex in this function because the parent is looping through multiple racers
	if racer.TwitchBotEnabled == 1 {
		if strings.HasPrefix(racer.StreamURL, "https://www.twitch.tv/") {
			// Parse for the username
			re, err := regexp.Compile(`https://www.twitch.tv/(.+)`)
			if err != nil {
				log.Error("Failed to compile the Twitch username regular expression:", err)
				return
			}
			channel := re.FindStringSubmatch(racer.StreamURL)[1]
			channel = strings.ToLower(channel)

			twitchSend(channel, message, racer.TwitchBotDelay)
		} else {
			log.Error("User \"" + racer.Name + "\" had a stream URL set to \"" + racer.StreamURL + "\" but their \"twitch_bot_enabled\" was set to 1, which should never happen.")
		}
	}
}

func twitchSend(channel string, message string, delay int) {
	// Don't talk in any channels that we are not a moderator in
	// (we will get banned for spamming otherwise)
	twitchModMap.RLock()
	isMod := twitchModMap.m[channel]
	twitchModMap.RUnlock()

	if isMod == false {
		return
	}

	go twitchSend2(channel, message, delay)
}

func twitchSend2(channel string, message string, delay int) {
	time.Sleep(time.Second * time.Duration(delay))
	ircSend(":" + twitchUsername + "!" + twitchUsername + "@" + twitchUsername + ".tmi.twitch.tv PRIVMSG #" + channel + " :" + message)
}

func twitchCheckForMod(channel string) {
	// If we are already checking for mod in this channel, don't do anything
	twitchModCheckingMap.RLock()
	isChecking := twitchModCheckingMap.m[channel]
	twitchModCheckingMap.RUnlock()

	if isChecking {
		return
	} else {
		twitchModCheckingMap.Lock()
		twitchModCheckingMap.m[channel] = true
		twitchModCheckingMap.Unlock()
	}

	// Wait 30 seconds before testing again to see if we are a moderator
	time.Sleep(time.Second * 30)

	// Double check to make sure that at least one user has their stream set to this Twitch channel and that they still have the Twitch Bot enabled
	streamURLs, err := db.Users.GetAllStreamURLs()
	if err != nil {
		twitchModCheckingMap.Lock()
		twitchModCheckingMap.m[channel] = false
		twitchModCheckingMap.Unlock()

		log.Error("Database error:", err)
		return
	}
	foundStream := false
	var twitchUser string
	for _, streamURL := range streamURLs {
		// Just in case, skip over streams that are not Twitch URLs
		if strings.HasPrefix(streamURL, "https://www.twitch.tv/") {
			// Parse for the username
			re, err := regexp.Compile(`https://www.twitch.tv/(.+)`)
			if err != nil {
				twitchModCheckingMap.Lock()
				twitchModCheckingMap.m[channel] = false
				twitchModCheckingMap.Unlock()

				log.Error("Failed to compile the Twitch username regular expression:", err)
				return
			}
			twitchUser = re.FindStringSubmatch(streamURL)[1]
			twitchUser = strings.ToLower(twitchUser)

			if channel == twitchUser {
				// We know that this stream has the Twitch Bot enabled because the GetAllStreams() function only returns streams with "twitch_bot_enabled = 1"
				foundStream = true
				break
			}
		}
	}

	if foundStream == false {
		// They either changed their stream URL or disabled the Twitch chat bot
		twitchModCheckingMap.Lock()
		twitchModCheckingMap.m[channel] = false
		twitchModCheckingMap.Unlock()
		return
	}

	// Force the server to give us a status update about whether we are a mod or not
	// (even though we are already joined to the channel, we can send another JOIN command without errors)
	twitchModCheckingMap.Lock()
	twitchModCheckingMap.m[channel] = false
	twitchModCheckingMap.Unlock()
	ircSend("JOIN #" + twitchUser)
}
