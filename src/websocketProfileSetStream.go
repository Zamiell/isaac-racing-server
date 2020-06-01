package main

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/Zamiell/isaac-racing-server/src/log"
	melody "gopkg.in/olahol/melody.v1"
)

/*
	This function sets the stream URL, the Twitch bot setting, and the Twitch bot delay
*/

func websocketProfileSetStream(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	userID := d.v.UserID
	username := d.v.Username
	oldStreamURL := d.v.StreamURL
	oldTwitchBotEnabled := d.v.TwitchBotEnabled
	oldTwitchBotDelay := d.v.TwitchBotDelay
	newStreamURL := d.Name
	newTwitchBotEnabled := d.Enabled
	newTwitchBotDelay := d.Value

	/*
		Validation
	*/

	// Validate that the stream URL is not empty
	if newStreamURL == "" {
		newStreamURL = "-"
	}

	// Prepare some regular expressions for later
	twitchStreamRegExp, err := regexp.Compile(`https://www.twitch.tv/(.+)`)
	if err != nil {
		log.Error("Failed to compile the Twitch stream regular expression.")
		websocketError(s, d.Command, "")
		return
	}
	twitchUserValidRegExp, err := regexp.Compile(`^[a-zA-Z0-9_]{4,25}$`)
	if err != nil {
		log.Error("Failed to compile the Twitch username validity regular expression.")
		websocketError(s, d.Command, "")
		return
	}

	/*
		Stream URL
	*/

	// Validate that the submitted stream URL is different than the current one
	if oldStreamURL != newStreamURL {
		// Validate that the submitted stream URL is not malicious
		if newStreamURL == "-" {
			// Do nothing
		} else if strings.HasPrefix(newStreamURL, "https://www.twitch.tv/") {
			// Do nothing
		} else {
			websocketError(s, d.Command, "Stream URLs must either be \"-\" or begin with \"https://www.twitch.tv/\".")
			return
		}

		// Check to see if anyone else has claimed this stream URL
		streamURLs, err := db.Users.GetAllStreamURLs()
		if err != nil {
			log.Error("Database error while getting all of the stream URLs:", err)
			websocketError(s, d.Command, "")
			return
		}
		for _, streamURL := range streamURLs {
			if strings.EqualFold(newStreamURL, streamURL) {
				websocketWarning(s, d.Command, "Someone else has already claimed that stream URL. If you are the real owner of this stream, please contact an administrator.")
				return
			}
		}

		// If this is a Twitch stream, validate that the Twitch username is valid
		if strings.HasPrefix(newStreamURL, "https://www.twitch.tv/") {
			// Parse the for username
			newTwitchUser := twitchStreamRegExp.FindStringSubmatch(newStreamURL)[1]

			// Validate the username
			// https://www.reddit.com/r/Twitch/comments/32w5b2/username_requirements/
			if twitchUserValidRegExp.FindString(newTwitchUser) == "" {
				websocketError(s, d.Command, "The stream URL submitted does not have a valid Twitch username.")
				return
			}
		}

		// Since they are changing streams, we no longer need the Twitch bot to be in their old channel
		if oldTwitchBotEnabled {
			// Parse for the old Twitch username
			oldTwitchUser := twitchStreamRegExp.FindStringSubmatch(oldStreamURL)[1]
			oldTwitchUser = strings.ToLower(oldTwitchUser)

			// Leave the channel
			twitchLeaveChannel(oldTwitchUser)
		}

		// Since they are changing streams, we need to join the bot to their new channel
		if newTwitchBotEnabled {
			// Validate that they have a Twitch stream URL set
			var newTwitchUser string
			if strings.HasPrefix(newStreamURL, "https://www.twitch.tv/") {
				// Parse for the username
				newTwitchUser = twitchStreamRegExp.FindStringSubmatch(newStreamURL)[1]
				newTwitchUser = strings.ToLower(newTwitchUser)
			}
			if newTwitchUser == "" {
				log.Warning("User \"" + username + "\" tried to enable the Twitch bot without having a Twitch stream URL set.")
				websocketError(s, d.Command, "You must have a Twitch stream URL set in order to use the Twitch chat bot.")
				return
			}

			// Join the channel
			twitchJoinChannel(newTwitchUser)
		}

		// Set the new stream URL in the database
		if err := db.Users.SetStreamURL(userID, newStreamURL); err != nil {
			log.Error("Database error while setting the stream URL for user "+strconv.Itoa(userID)+":", err)
			websocketError(s, d.Command, "")
			return
		}

		// Set the new stream URL in the WebSocket session
		s.Set("streamURL", newStreamURL)

		// It has to also be updated in all chat rooms
		chatRoomsUpdate(username, "StreamURL", newStreamURL)
	}

	/*
		Twitch bot enabled
	*/

	// Validate that the submitted Twitch bot setting is different than the current one
	if oldTwitchBotEnabled != newTwitchBotEnabled {
		if newTwitchBotEnabled {
			// Validate that they have a Twitch stream URL set
			var newTwitchUser string
			if strings.HasPrefix(newStreamURL, "https://www.twitch.tv/") {
				// Parse for the username
				newTwitchUser = twitchStreamRegExp.FindStringSubmatch(newStreamURL)[1]
				newTwitchUser = strings.ToLower(newTwitchUser)
			}
			if newTwitchUser == "" {
				log.Warning("User \"" + username + "\" tried to enable the Twitch bot without having a Twitch stream URL set.")
				websocketError(s, d.Command, "You must have a Twitch stream URL set in order to use the Twitch chat bot.")
				return
			}

			// If it is a Twitch stream, make the Twitch IRC bot join their channel
			twitchJoinChannel(newTwitchUser)
		} else {
			// Parse for the username
			oldTwitchUser := twitchStreamRegExp.FindStringSubmatch(oldStreamURL)[1]
			oldTwitchUser = strings.ToLower(oldTwitchUser)

			// Leave the channel
			twitchLeaveChannel(oldTwitchUser)
		}

		// Set the new Twitch bot setting in the database
		if err := db.Users.SetTwitchBotEnabled(userID, newTwitchBotEnabled); err != nil {
			log.Error("Database error while setting the twitch bot setting for user "+strconv.Itoa(userID)+":", err)
			websocketError(s, d.Command, "")
			return
		}

		// Set the new stream URL in the WebSocket session
		s.Set("twitchBotEnabled", newTwitchBotEnabled)
	}

	/*
		Twitch bot delay
	*/

	// Validate that the delay is different than the current one
	if oldTwitchBotDelay != newTwitchBotDelay {
		// Validate that it is a sane delay
		if newTwitchBotDelay < 0 || newTwitchBotDelay > 60 {
			websocketError(s, d.Command, "Your Twitch bot delay must be between 0 and 60.")
			return
		}

		// Set the new Twitch bot delay in the database
		if err := db.Users.SetTwitchBotDelay(userID, newTwitchBotDelay); err != nil {
			log.Error("Database error while setting the twitch bot delay for user "+strconv.Itoa(userID)+":", err)
			websocketError(s, d.Command, "")
			return
		}

		// Set the new stream URL in the WebSocket session
		s.Set("twitchBotDelay", newTwitchBotDelay)
	}
}
