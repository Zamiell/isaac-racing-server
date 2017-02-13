package main

/*
	Imports
*/

import (
	"regexp"
	"strings"
)

/*
	WebSocket profile command functions
*/

func profileSetStream(conn *ExtendedConnection, data *IncomingCommandMessage) {
	/*
		This function gets the stream URL, the Twitch bot setting, and the Twitch bot delay
	*/

	// Local variables
	functionName := "profileSetStream"
	userID := conn.UserID
	username := conn.Username
	newStreamURL := data.Name
	newTwitchBotEnabled := data.Enabled
	newTwitchBotDelay := data.Value

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()

	// Log the received command
	log.Debug("User \"" + username + "\" sent a " + functionName + " command.")

	// Rate limit all commands
	if commandRateLimit(conn) == true {
		return
	}

	// Get the user's current stream URL
	oldStreamURL, err := db.Users.GetStreamURL(userID)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Get the user's current Twitch bot setting
	oldTwitchBotEnabled, err := db.Users.GetTwitchBotEnabled(userID)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Get the user's current Twitch bot delay
	oldTwitchBotDelay, err := db.Users.GetTwitchBotDelay(userID)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Prepare some regular expressions for later
	twitchStreamRegExp, err := regexp.Compile(`https://www.twitch.tv/(.+)`)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Failed to compile the Twitch stream regular expression.")
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}
	twitchUserValidRegExp, err := regexp.Compile(`^[a-zA-Z0-9_]{4,25}$`)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Failed to compile the Twitch username validity regular expression.")
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Validate that the submitted stream URL is different than the current one
	if oldStreamURL != newStreamURL {
		// Validate that the submitted stream URL is not malicious
		if newStreamURL == "-" {
			// Do nothing
		} else if strings.HasPrefix(newStreamURL, "https://www.twitch.tv/") {
			// Do nothing
		} else {
			commandMutex.Unlock()
			connError(conn, functionName, "Stream URLs must either be \"-\" or begin with \"https://www.twitch.tv/\".")
			return
		}

		// Check to see if anyone else has claimed this stream URL
		streamURLs, err := db.Users.GetAllStreamURLs()
		if err != nil {
			commandMutex.Unlock()
			log.Error("Database error:", err)
			connError(conn, functionName, "Something went wrong. Please contact an administrator.")
			return
		}
		for _, streamURL := range streamURLs {
			if strings.ToLower(newStreamURL) == strings.ToLower(streamURL) {
				commandMutex.Unlock()
				connWarning(conn, functionName, "Someone else has already claimed that stream URL. If you are the real owner of this stream, please contact an administrator.")
				return
			}
		}

		// If this is a Twitch stream, validate that the Twitch username is valid
		if strings.HasPrefix(newStreamURL, "https://www.twitch.tv/") {
			// Parse the for username
			newTwitchUser := twitchStreamRegExp.FindStringSubmatch(newStreamURL)[1]

			// Validate the username (from https://www.reddit.com/r/Twitch/comments/32w5b2/username_requirements/)
			if twitchUserValidRegExp.FindString(newTwitchUser) == "" {
				commandMutex.Unlock()
				connError(conn, functionName, "The stream URL submitted does not have a valid Twitch username.")
				return
			}
		}

		// Since they are changing streams, we no longer need the Twitch bot to be in their old channel
		if oldTwitchBotEnabled {
			// Parse for the old Twitch username
			oldTwitchUser := twitchStreamRegExp.FindStringSubmatch(oldStreamURL)[1]
			oldTwitchUser = strings.ToLower(oldTwitchUser)

			// Leave the channel
			ircSend("PART #" + oldTwitchUser)
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
				commandMutex.Unlock()
				log.Warning("User \"" + username + "\" tried to enable the Twitch bot without having a Twitch stream URL set.")
				connError(conn, functionName, "You must have a Twitch stream URL set in order to use the Twitch chat bot.")
				return
			}

			// Join the channel
			ircSend("JOIN #" + newTwitchUser)
		}

		// Set the new stream URL in the database
		if err := db.Users.SetStreamURL(userID, newStreamURL); err != nil {
			commandMutex.Unlock()
			log.Error("Database error:", err)
			connError(conn, functionName, "Something went wrong. Please contact an administrator.")
			return
		}

		// Update the connections for everyone
		// TODO
	}

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
				commandMutex.Unlock()
				log.Warning("User \"" + username + "\" tried to enable the Twitch bot without having a Twitch stream URL set.")
				connError(conn, functionName, "You must have a Twitch stream URL set in order to use the Twitch chat bot.")
				return
			}

			// If it is a Twitch stream, make the Twitch IRC bot join their channel
			ircSend("JOIN #" + newTwitchUser)
		} else {
			// Parse for the username
			oldTwitchUser := twitchStreamRegExp.FindStringSubmatch(oldStreamURL)[1]
			oldTwitchUser = strings.ToLower(oldTwitchUser)

			// Leave the channel
			ircSend("PART #" + oldTwitchUser)
		}

		// Set the new Twitch bot setting in the database
		if err := db.Users.SetTwitchBotEnabled(userID, newTwitchBotEnabled); err != nil {
			commandMutex.Unlock()
			log.Error("Database error:", err)
			connError(conn, functionName, "Something went wrong. Please contact an administrator.")
			return
		}
	}

	// Validate that the delay is different than the current one
	if oldTwitchBotDelay != newTwitchBotDelay {
		// Validate that it is a sane delay
		if newTwitchBotDelay < 0 || newTwitchBotDelay > 60 {
			commandMutex.Unlock()
			connError(conn, functionName, "Your Twitch bot delay must be between 0 and 60.")
			return
		}

		// Set the new Twitch bot delay in the database
		if err := db.Users.SetTwitchBotDelay(userID, newTwitchBotDelay); err != nil {
			commandMutex.Unlock()
			log.Error("Database error:", err)
			connError(conn, functionName, "Something went wrong. Please contact an administrator.")
			return
		}
	}

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()
}
