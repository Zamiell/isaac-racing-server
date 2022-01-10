package server

import (
	"strconv"

	melody "gopkg.in/olahol/melody.v1"
)

/*
	Race validation subroutines
*/

func raceValidateRuleset(s *melody.Session, d *IncomingWebsocketData) bool {
	ruleset := d.Ruleset

	// Validate the ruleset format
	if ruleset.Format != RaceFormatUnseeded &&
		ruleset.Format != RaceFormatSeeded &&
		ruleset.Format != RaceFormatDiversity &&
		ruleset.Format != RaceFormatCustom {

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
	if ruleset.Goal != RaceGoalBlueBaby &&
		ruleset.Goal != RaceGoalTheLamb &&
		ruleset.Goal != RaceGoalMegaSatan &&
		ruleset.Goal != RaceGoalHush &&
		ruleset.Goal != RaceGoalDelirium &&
		ruleset.Goal != RaceGoalMother &&
		ruleset.Goal != RaceGoalBeast &&
		ruleset.Goal != RaceGoalBossRush &&
		ruleset.Goal != RaceGoalCustom {

		websocketError(s, d.Command, "That is not a valid goal.")
		return false
	}

	// Validate the starting build
	if ruleset.Format != RaceFormatSeeded &&
		ruleset.StartingBuild != -1 {

		websocketError(s, d.Command, "You cannot set a starting build for a non-seeded race.")
		return false
	}

	// Validate multiplayer ranked games
	if !ruleset.Solo {
		if ruleset.Ranked {
			websocketError(s, d.Command, "Multiplayer races must be unranked.")
			return false
		}

		// Set the ruleset to ranked since it is a multiplayer game
		// (in the past, there was multiplayer unranked and ranked,
		// so this is a monkey fix to avoid changing the client)
		return true
	}

	// Validate ranked solo games
	if ruleset.Ranked && ruleset.Solo {
		return raceValidateRulesetRankedSolo(s, d)
	}

	// Validate seeded games
	if ruleset.Format == RaceFormatSeeded {
		return raceValidateRulesetSeeded(s, d)
	}

	return true
}

func raceValidateRulesetRankedSolo(s *melody.Session, d *IncomingWebsocketData) bool {
	ruleset := d.Ruleset

	if ruleset.Format != RaceFormatSeeded {
		websocketError(s, d.Command, "Ranked solo races must be seeded.")
		return false
	}

	if ruleset.Character != "Judas" {
		websocketError(s, d.Command, "Ranked solo races must have a character of Judas.")
		return false
	}

	if ruleset.Goal != "Blue Baby" {
		websocketError(s, d.Command, "Ranked solo races must have a goal of Blue Baby.")
		return false
	}

	// Validate the difficulty
	if ruleset.Difficulty != "normal" {
		websocketError(s, d.Command, "That is not a valid difficulty.")
		return false
	}

	return true
}

func raceValidateRulesetSeeded(s *melody.Session, d *IncomingWebsocketData) bool {
	ruleset := d.Ruleset

	if ruleset.StartingBuild < 0 || ruleset.StartingBuild >= len(allBuilds) { // 0 is random
		msg := "The build of \"" + strconv.Itoa(ruleset.StartingBuild) + "\" is not a valid starting build."
		websocketError(s, d.Command, msg)
		return false
	}

	if ruleset.Character == "Tainted Lazarus" {
		msg := "Tainted Lazarus is illegal for seeded races since his mechanics are difficult to seed properly."
		websocketError(s, d.Command, msg)
		return false
	}

	illegalCharacters := buildExceptions[ruleset.StartingBuild]
	if stringInSlice(ruleset.Character, illegalCharacters) {
		msg := "The character of " + ruleset.Character + " is illegal in combination with the starting build of: " + getBuildName(ruleset.StartingBuild)
		websocketError(s, d.Command, msg)
		return false
	}

	return true
}
