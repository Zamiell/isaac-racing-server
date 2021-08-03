package server

import (
	"math/rand"
	"strconv"
	"time"
	"unicode/utf8"

	melody "gopkg.in/olahol/melody.v1"
)

const (
	rateLimitRate       = float64(4)  // Number of races created
	rateLimitPer        = float64(60) // Per seconds
	automaticBanAdminID = 1
	automaticBanReason  = "race spamming"
)

func websocketRaceCreate(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	username := d.v.Username
	admin := d.v.Admin
	rateLimitAllowance := d.v.RateLimitAllowance
	rateLimitLastCheck := d.v.RateLimitLastCheck
	name := d.Name
	ruleset := d.Ruleset
	password := d.Password

	/*
		Validation
	*/

	// Validate that the server is not shutting down soon
	if shutdownMode > 0 && admin == 0 {
		websocketWarning(s, d.Command, "The server is restarting soon (when all ongoing races have finished). You cannot start any new races for the time being.")
		return
	}

	// Validate that the race name cannot be empty
	if name == "" {
		name = "-"
	}

	// Validate that the race name is not longer than 100 characters
	if utf8.RuneCountInString(name) > 100 {
		logger.Warning("User \"" + username + "\" sent a race name longer than 100 characters.")
		websocketError(s, d.Command, "Race names must not be longer than 100 characters.")
		return
	}

	// Validate that the ruleset options cannot be empty
	if ruleset.Format == "" {
		ruleset.Format = RaceFormatUnseeded
	}
	if ruleset.Character == "" {
		ruleset.Character = "Judas"
	}
	if ruleset.Goal == "" {
		ruleset.Goal = "The Chest"
	}

	// Validate the submitted ruleset
	if !raceValidateRuleset(s, d) {
		return
	}

	// Fix the ranking for multiplayer races
	if !ruleset.Solo {
		ruleset.Ranked = true
	} else {
		password = ""
	}

	// Pick a random character, if necessary
	if ruleset.Character == "random" {
		ruleset.CharacterRandom = true
		rand.Seed(time.Now().UnixNano())
		ruleset.Character = characters[rand.Intn(len(characters))] // nolint: gosec
	}

	// Pick a random starting build, if necessary
	// (StartingBuild will be -1 in non-seeded races)
	if ruleset.StartingBuild == 0 {
		ruleset.StartingBuildRandom = true
		rand.Seed(time.Now().UnixNano())
		// We don't want to select index 0
		ruleset.StartingBuild = rand.Intn(len(allBuilds)-1) + 1 // nolint: gosec
	}

	// Check if there are any ongoing races with this name
	for _, race := range races {
		if race.Name == name {
			websocketError(s, d.Command, "There is already a non-finished race with that name.")
			return
		}
	}

	// Validate that the user is not creating new races over and over, which will generate an annoying sound effect for everyone in the lobby
	// Algorithm from: http://stackoverflow.com/questions/667508/whats-a-good-rate-limiting-algorithm
	// (allow staff/admins to create unlimited races)
	if admin == 0 && !ruleset.Solo {
		now := time.Now()
		timePassed := now.Sub(rateLimitLastCheck).Seconds()
		s.Set("rateLimitLastCheck", now)
		logger.Info("User \"" + username + "\" has \"" + strconv.FormatFloat(timePassed, 'f', 2, 64) + "\" time passed since the last race creation.")

		newRateLimitAllowance := rateLimitAllowance + timePassed*(rateLimitRate/rateLimitPer)
		if newRateLimitAllowance > rateLimitRate {
			newRateLimitAllowance = rateLimitRate
		}

		if newRateLimitAllowance < 1 {
			// They are spamming new races, so automatically ban them as punishment
			logger.Warning("User \"" + username + "\" triggered new race rate-limiting; banning them.")
			ban(s, d)
			return
		}

		newRateLimitAllowance--
		s.Set("rateLimitAllowance", newRateLimitAllowance)
	}

	/*
		Create
	*/

	// Create and set a seed if necessary
	ruleset.Seed = "-"
	if ruleset.Format == RaceFormatSeeded {
		// Create a random Isaac seed
		// (using the current Epoch timestamp as a seed)
		ruleset.Seed = isaacGetRandomSeed()
	} else if ruleset.Format == RaceFormatDiversity {
		ruleset.Seed = diversityGetSeed(ruleset)
	}

	// Zero out the starting items field
	ruleset.StartingItems = make([]int, 0)

	/*
		Create the race in the database
		(it will have no data associated with it other than the automatically
		generated row ID; we want to use this ID as a unique map key)

		The benefit of doing this is that we won't reuse any race IDs after a
		server restart or crash.
		Furthermore, we want the ability for racers to be able to submit a race
		comment after the race has already ended. (Races are deleted from the
		internal map upon finishing.)
	*/
	var raceID int
	if v, err := db.Races.Insert(); err != nil {
		logger.Error("Database error while inserting the race:", err)
		websocketError(s, d.Command, "")
		return
	} else {
		raceID = v
	}

	// Create the race and keep track of it in the races map
	race := &Race{
		ID:              raceID,
		Name:            name,
		Status:          RaceStatusOpen,
		Ruleset:         ruleset,
		Captain:         username,
		Password:        password,
		SoundPlayed:     false,
		DatetimeCreated: getTimestamp(),
		DatetimeStarted: 0,
		Racers:          make(map[string]*Racer),
	}
	races[raceID] = race

	// Send everyone a notification that a new race has been started
	for _, s := range websocketSessions {
		websocketEmit(s, "raceCreated", &RaceCreatedMessage{
			ID:                  race.ID,
			Name:                race.Name,
			Status:              race.Status,
			Ruleset:             race.Ruleset,
			Captain:             race.Captain,
			IsPasswordProtected: len(race.Password) > 0,
			DatetimeCreated:     race.DatetimeCreated,
			DatetimeStarted:     race.DatetimeStarted,
			Racers:              make([]string, 0),
		})
	}

	d.ID = race.ID
	websocketRaceJoin(s, d)
}

func ban(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	userID := d.v.UserID

	/*
		This code is copied from the "websocketAdminBan()" function
	*/

	// Add this username to the ban list in the database
	if err := db.BannedUsers.Insert(userID, automaticBanAdminID, automaticBanReason); err != nil {
		logger.Error("Database error while userting the banned user:", err)
		websocketError(s, d.Command, "")
		return
	}

	// Add their IP to the banned IP list
	if err := db.BannedIPs.InsertUserIP(userID, automaticBanAdminID, automaticBanReason); err != nil {
		logger.Error("Database error while inserting the banned IP:", err)
		websocketError(s, d.Command, "")
		return
	}

	websocketError(
		s,
		"Banned",
		"New race spamming detected. You have been banned. If you think this was a mistake, please contact the administration to appeal.",
	)
	websocketClose(s)
}
