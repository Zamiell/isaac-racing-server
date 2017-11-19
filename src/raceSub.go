package main

import (
	"math/rand"
	"time"

	melody "gopkg.in/olahol/melody.v1"
)

/*
	Constants
*/

const numBuilds = 33

var characters = []string{
	"Isaac",
	"Magdalene",
	"Cain",
	"Judas",
	"Blue Baby",
	"Eve",
	"Samson",
	"Azazel",
	"Lazarus",
	"Eden",
	"The Lost",
	"Lilith",
	"Keeper",
	"Apollyon",
	"Samael",
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
		ruleset.Format != "unseeded-lite" &&
		ruleset.Format != "seeded-hard" &&
		ruleset.Format != "custom" {

		websocketError(s, d.Command, "That is not a valid ruleset.")
		return false
	}

	// Validate the character
	validCharacter := false
	for _, character := range characters {
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
	if ruleset.Character == "random" {
		ruleset.CharacterRandom = true
		rand.Seed(time.Now().UnixNano())
		ruleset.Character = characters[rand.Intn(len(characters))]
	}

	// Validate the goal
	if ruleset.Goal != "Blue Baby" &&
		ruleset.Goal != "The Lamb" &&
		ruleset.Goal != "Mega Satan" &&
		ruleset.Goal != "Everything" &&
		ruleset.Goal != "custom" {

		websocketError(s, d.Command, "That is not a valid goal.")
		return false
	}

	// Validate the starting build
	if ruleset.Format != "seeded" &&
		ruleset.Format != "seeded-hard" &&
		ruleset.StartingBuild != -1 {

		websocketError(s, d.Command, "You cannot set a starting build for a non-seeded race.")
		return false
	} else if (ruleset.Format == "seeded" || ruleset.Format == "seeded-hard") &&
		(ruleset.StartingBuild < 0 || ruleset.StartingBuild > numBuilds) { // There are 33 builds (0 is random)

		websocketError(s, d.Command, "That is not a valid starting build.")
		return false
	}
	if ruleset.StartingBuild == 0 {
		ruleset.StartingBuildRandom = true
		rand.Seed(time.Now().UnixNano())
		ruleset.StartingBuild = rand.Intn(numBuilds) + 1 // 1 to numBuilds
	}

	// Validate ranked games
	if ruleset.Ranked &&
		ruleset.Format != "seeded" &&
		ruleset.Format != "unseeded" {

		websocketError(s, d.Command, "Ranked races must be either seeded or unseeded.")
		return false
	}
	if ruleset.Ranked &&
		ruleset.Format == "seeded" &&
		ruleset.StartingBuild != 0 {

		websocketError(s, d.Command, "Ranked seeded races must have a random starting build.")
		return false
	}

	// Validate unseeded ranked games
	if ruleset.Ranked && ruleset.Format == "unseeded" && ruleset.Character != "Judas" {
		websocketError(s, d.Command, "Ranked unseeded races must have a character of Judas.")
		return false
	} else if ruleset.Ranked && ruleset.Format == "unseeded" && ruleset.Goal != "Blue Baby" {
		websocketError(s, d.Command, "Ranked unseeded races must have a goal of Blue Baby.")
		return false
	}

	return true
}
