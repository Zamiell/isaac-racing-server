package main

import (
	melody "gopkg.in/olahol/melody.v1"
)

/*
	Constants
*/

var characters = []string{
	"Isaac",         // 0
	"Magdalene",     // 1
	"Cain",          // 2
	"Judas",         // 3
	"Blue Baby",     // 4
	"Eve",           // 5
	"Samson",        // 6
	"Azazel",        // 7
	"Lazarus",       // 8
	"Eden",          // 9
	"The Lost",      // 10
	"Lilith",        // 11
	"Keeper",        // 12
	"Apollyon",      // 13
	"The Forgotten", // 14
	"Samael",        // 15
	"Random Baby",   // 16
}

/*
	Race validation subroutines
*/

func raceValidateRuleset(s *melody.Session, d *IncomingWebsocketData) bool {
	// Local variables
	ruleset := d.Ruleset

	// Validate the ruleset format
	if ruleset.Format != "unseeded" &&
		ruleset.Format != "seeded" &&
		ruleset.Format != "diversity" &&
		ruleset.Format != "custom" {

		websocketError(s, d.Command, "That is not a valid ruleset.")
		return false
	}

	// Validate the character
	validCharacter := false
	for _, character := range characters { // Valid characters are defined above
		if ruleset.Character == character {
			validCharacter = true
			break
		}
	}
	if ruleset.Character == "random" {
		validCharacter = true
	}
	if !validCharacter {
		websocketError(s, d.Command, "That is not a valid character.")
		return false
	}

	// Validate the goal
	if ruleset.Goal != "Blue Baby" &&
		ruleset.Goal != "The Lamb" &&
		ruleset.Goal != "Mega Satan" &&
		ruleset.Goal != "Hush" &&
		ruleset.Goal != "Delirium" &&
		ruleset.Goal != "Boss Rush" &&
		ruleset.Goal != "Everything" &&
		ruleset.Goal != "custom" {

		websocketError(s, d.Command, "That is not a valid goal.")
		return false
	}

	// Validate the starting build
	if ruleset.Format != "seeded" &&
		ruleset.StartingBuild != -1 {

		websocketError(s, d.Command, "You cannot set a starting build for a non-seeded race.")
		return false
	} else if ruleset.Format == "seeded" &&
		(ruleset.StartingBuild < 0 || ruleset.StartingBuild > len(allBuilds)) { // 0 is random

		websocketError(s, d.Command, "That is not a valid starting build.")
		return false
	}

	// Validate multiplayer ranked games
	if !ruleset.Solo {
		if ruleset.Ranked {
			websocketError(s, d.Command, "Multiplayer races must be unranked.")
			return false
		} else {
			// Set the ruleset to ranked since it is a multiplayer game
			// (in the past, there was multiplayer unranked and ranked,
			// so this is a monkey fix to avoid changing the client)
			return true
		}
	}

	// Validate solo ranked games
	if !ruleset.Ranked {
		return true
	}
	if ruleset.Format != "seeded" &&
		ruleset.Format != "unseeded" {

		websocketError(s, d.Command, "Solo ranked races must be either seeded or unseeded.")
		return false
	}
	if ruleset.Character != "Judas" {
		websocketError(s, d.Command, "Solo ranked races must have a character of Judas.")
		return false
	}
	if ruleset.Goal != "Blue Baby" {
		websocketError(s, d.Command, "Solo ranked races must have a goal of Blue Baby.")
		return false
	}
	if ruleset.Format == "seeded" &&
		ruleset.StartingBuild != 0 {

		websocketError(s, d.Command, "Solo ranked seeded races must have a random starting build.")
		return false
	}

	// Validate the difficulty
	if ruleset.Difficulty != "normal" && ruleset.Difficulty != "hard" {
		websocketError(s, d.Command, "That is not a valid difficulty.")
		return false
	}

	return true
}
