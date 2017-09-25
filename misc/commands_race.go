package main

/*
	WebSocket race command functions
*/

// Currently not implemented client-side
/*
func websocketRaceJoinSpectate(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	userID := d.v.UserID
	username := d.v.Username
	raceID := d.ID

	// Validate basic things about the race ID
	if !raceValidate(s, d) {
		return
	}

	// Validate that they are not in the race
	if !raceValidateOut2(s, d) {
		return
	}

	// Add this user to the participants list for that race
	if err := db.RaceParticipants.Insert(userID, raceID); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	}

	// Send everyone a notification that the user joined
	for _, conn := range websocketSessions {
		websocketEmit(s, "raceJoined", &RaceMessage{
			raceID,
			username,
		})
	}

	// Send it to the user
	websocketEmit(s, "racerList", &RacerList{
		raceID,
		racerList,
	})

	// Join the user to the channel for that race
	d.Room = "_race_"+strconv.Itoa(raceID)
	websocketRoomJoinSub(s, d)
}
*/

// Currently not implemented client-side
/*
func websocketRaceRuleset(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	userID := d.v.UserID
	username := d.v.Username
	raceID := d.ID
	ruleset := d.Ruleset

	// Get the current ruleset
	// TODO

	// Check to see if anything has changed
	if currentRuleset.Format == ruleset.Format &&
		currentRuleset.Character == ruleset.Character &&
		currentRuleset.Goal == ruleset.Goal &&
		currentRuleset.StartingBuild == ruleset.StartingBuild {

		websocketError(s, d.Command, "The race ruleset is already set to those values.")
		return
	}

	// If they didn't specify something, set it to the existing value
	if ruleset.Format == "" {
		ruleset.Format = currentRuleset.Format
	}
	if ruleset.Character == "" {
		ruleset.Character = currentRuleset.Character
	}
	if ruleset.Goal == "" {
		ruleset.Goal = currentRuleset.Goal
	}
	if ruleset.StartingBuild == 0 {
		ruleset.StartingBuild = currentRuleset.StartingBuild
	}

	// Validate the submitted ruleset
	if !raceValidateRuleset(s, d) {
		return
	}

	// Validate basic things about the race ID
	if !raceValidate(s, d) {
		return
	}

	// Validate that the race is open
	if !raceValidateStatus(s, d, "open") {
		return
	}

	// Validate that they are in the race
	if !raceValidateIn2(s, d) {
		return
	}

	// Validate that they are the race captain
	// TODO

	// Get and set a seed if necessary
	if (ruleset.Format == "seeded" || ruleset.Format == "diversity") && ruleset.Format != currentRuleset.Format {
		var seed string
		if ruleset.Format == "seeded" {
			// TODO
		} else if ruleset.Format == "diversity" {
			// TODO
		}

		// Set the new seed
		// TODO
	}

	// Change the ruleset
	// TODO

	// Set everyone's status to "not ready"
	if err := db.RaceParticipants.SetAllStatus(raceID, "not ready"); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	}

	// Send everyone a notification that the ruleset has changed for this race
	for _, conn := range websocketSessions {
		type RaceSetRulesetMessage struct {
			ID      int            `json:"id"`
			Ruleset models.Ruleset `json:"ruleset"`
		}
		websocketEmit(s, "raceSetRuleset", &RaceSetRulesetMessage{
			raceID,
			ruleset,
		})
	}
}
*/

// Currently not implemented client-side
/*
func websocketRaceComment(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	userID := d.v.UserID
	username := d.v.Username
	muted := d.v.Muted
	raceID := d.ID
	comment := d.Comment

	// Strip leading and trailing whitespace from the comment
	comment = strings.TrimSpace(comment)

	// Validate that the comment is not empty
	if comment == "" {
		websocketError(s, d.Command, "That is an invalid comment.")
		return
	}

	// Validate that the comment is not excessively long
	if len(comment) < 150 {
		websocketError(s, d.Command, "Comments must not be longer than 150 characters.")
		return
	}

	// Validate that the comment does not contain special characters
	if !isAlphaNumericUnderscore(comment) {
		websocketError(s, d.Command, "Your comment must contain only letters, numbers, and underscores.")
		return
	}

	// Validate that the user is not muted
	if muted {
		websocketError(s, d.Command, "You have been muted by an administrator, so you cannot submit comments.")
		return
	}

	// Validate basic things about the race ID
	if !raceValidate(s, d) {
		return
	}

	// Validate that the race has started
	// TODO needs custom logic to verify that it is either "in progress" or "finished"

	// Validate that they are in the race
	if !raceValidateIn2(s, d) {
		return
	}

	// Set their comment in the database
	if err := db.RaceParticipants.SetComment(userID, raceID, comment); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	}

	// Send a notification to all the people in this particular race that the user added or changed their comment
	for _, racer := range racerNames {
		// Not all racers may be online during a race
		if s, ok := websocketSessions[racer]; ok {
			type RacerSetCommentMessage struct {
				ID      int    `json:"id"`
				Name    string `json:"name"`
				Comment string `json:"comment"`
			}
			websocketEmit(s, "racerSetComment", &RacerSetCommentMessage{
				raceID,
				username,
				comment,
			})
		}
	}
}
*/
