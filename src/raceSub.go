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

		websocketWarning(s, d.Command, "That is not a valid ruleset.")
		return false
	}

	if ruleset.Format == RaceFormatSeeded && ruleset.Solo {
		websocketWarning(
			s,
			d.Command,
			"Racing+ online ranked solo season 2 has concluded. You cannot play ranked solo until season 3 starts.",
		)
		return false
	}

	// Validate the character
	// (valid characters are defined in "characters.go")
	validCharacter := stringInSlice(ruleset.Character, characters)
	if ruleset.Character == "random" {
		validCharacter = true
	}
	if !validCharacter {
		websocketWarning(s, d.Command, "That is not a valid character.")
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

		websocketWarning(s, d.Command, "That is not a valid goal.")
		return false
	}

	// Validate the starting build
	if ruleset.Format == RaceFormatSeeded {
		if ruleset.StartingBuild < 0 || ruleset.StartingBuild >= len(allBuilds) { // 0 is random
			msg := "The build of \"" + strconv.Itoa(ruleset.StartingBuild) + "\" is not a valid starting build."
			websocketWarning(s, d.Command, msg)
			return false
		}
	} else {
		if ruleset.StartingBuild != -1 {
			websocketWarning(s, d.Command, "You cannot set a starting build for a non-seeded race.")
			return false
		}
	}

	// Validate specific things for seeded races
	if ruleset.Format == RaceFormatSeeded {
		// Check for character + build anti-synergies
		illegalCharacters := buildExceptions[ruleset.StartingBuild]
		if stringInSlice(ruleset.Character, illegalCharacters) {
			msg := "The character of " + ruleset.Character + " is illegal in combination with the starting build of: " + getBuildName(ruleset.StartingBuild)
			websocketWarning(s, d.Command, msg)
			return false
		}

		if ruleset.Character == "Tainted Lazarus" {
			msg := "Tainted Lazarus is illegal for seeded races since his mechanics are difficult to seed properly."
			websocketWarning(s, d.Command, msg)
			return false
		}
	}

	if ruleset.Solo {
		return raceValidateRulesetSolo(s, d)
	}

	return raceValidateRulesetMultiplayer(s, d)
}

func raceValidateRulesetSolo(s *melody.Session, d *IncomingWebsocketData) bool {
	ruleset := d.Ruleset

	// Validate ranked solo games
	if ruleset.Ranked {
		return raceValidateRulesetRankedSolo(s, d)
	}

	return true
}

func raceValidateRulesetRankedSolo(s *melody.Session, d *IncomingWebsocketData) bool {
	ruleset := d.Ruleset

	if ruleset.Format != RaceFormatSeeded {
		websocketWarning(s, d.Command, "Ranked solo races must be seeded.")
		return false
	}

	if ruleset.Character != "Judas" {
		websocketWarning(s, d.Command, "Ranked solo races must have a character of Judas.")
		return false
	}

	if ruleset.Goal != "Blue Baby" {
		websocketWarning(s, d.Command, "Ranked solo races must have a goal of Blue Baby.")
		return false
	}

	// Validate the difficulty
	if ruleset.Difficulty != "normal" {
		websocketWarning(s, d.Command, "That is not a valid difficulty.")
		return false
	}

	return true
}

func raceValidateRulesetMultiplayer(s *melody.Session, d *IncomingWebsocketData) bool {
	ruleset := d.Ruleset

	if !ruleset.Ranked {
		websocketWarning(s, d.Command, "Multiplayer races must be ranked.")
		return false
	}

	return true
}
