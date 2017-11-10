package main

import (
	melody "gopkg.in/olahol/melody.v1"
)

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
	if ruleset.Character != "Isaac" &&
		ruleset.Character != "Magdalene" &&
		ruleset.Character != "Cain" &&
		ruleset.Character != "Judas" &&
		ruleset.Character != "Blue Baby" &&
		ruleset.Character != "Eve" &&
		ruleset.Character != "Samson" &&
		ruleset.Character != "Azazel" &&
		ruleset.Character != "Lazarus" &&
		ruleset.Character != "Eden" &&
		ruleset.Character != "The Lost" &&
		ruleset.Character != "Lilith" &&
		ruleset.Character != "Keeper" &&
		ruleset.Character != "Apollyon" &&
		ruleset.Character != "Samael" {

		websocketError(s, d.Command, "That is not a valid character.")
		return false
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
		(ruleset.StartingBuild < 1 || ruleset.StartingBuild > 33) { // There are 33 builds

		websocketError(s, d.Command, "That is not a valid starting build.")
		return false
	}

	// Validate ranked games
	if ruleset.Ranked &&
		ruleset.Format != "seeded" &&
		ruleset.Format != "unseeded" {

		websocketError(s, d.Command, "Ranked races must be either seeded or unseeded.")
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
