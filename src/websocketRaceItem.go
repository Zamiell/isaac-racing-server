package main

/*
	Imports
*/

import (
	"strconv"

	"github.com/Zamiell/isaac-racing-server/src/log"
	"github.com/Zamiell/isaac-racing-server/src/models"
	melody "gopkg.in/olahol/melody.v1"
)

func websocketRaceItem(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	raceID := d.ID
	itemID := d.ItemID
	userID := d.v.UserID
	username := d.v.Username

	// Validate that the item number is sane
	// The highest item ID is 510 (Delirious), and the Racing+ mod has a bunch of custom items
	// So just check for over 600 to be safe
	if itemID < 1 || itemID > 600 {
		log.Warning("User \"" + username + "\" attempted to add item " + strconv.Itoa(itemID) + " to their build, but that is a bogus number.")
		websocketError(s, d.Command, "That is not a valid item ID.")
		return
	}

	// Validate basic things about the race ID
	if !raceValidate(s, d) {
		return
	}

	// Validate that the race has started
	if !raceValidateStatus(s, d, "in progress") {
		return
	}

	// Validate that they are in the race
	if !raceValidateIn2(s, d) {
		return
	}

	// Validate that their status is set to "racing" status
	if !racerValidateStatus(s, d, "racing") {
		return
	}

	// Get their current floor
	floorNum, stageType, err := db.RaceParticipants.GetFloor(userID, raceID)
	if err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	}

	// Add this item to their build
	if err = db.RaceParticipantItems.Insert(userID, raceID, itemID, floorNum, stageType); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	}

	// Get the list of racers for this race
	racerNames, err := db.RaceParticipants.GetRacerNames(raceID)
	if err != nil {
		log.Error("Database error:", err)
		return
	}

	// Send a notification to all the people in this particular race that the user got an item
	for _, racer := range racerNames {
		// Not all racers may be online during a race
		if s, ok := websocketSessions[racer]; ok {
			item := models.Item{
				ID:        itemID,
				FloorNum:  floorNum,
				StageType: stageType,
			}
			websocketEmit(s, "racerAddItem", &RacerAddItemMessage{raceID, username, item})
		}
	}
}
