package main

import (
	"strconv"

	"github.com/Zamiell/isaac-racing-server/src/log"
	melody "gopkg.in/olahol/melody.v1"
)

func websocketRaceItem(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	username := d.v.Username
	raceID := d.ID
	itemID := d.ItemID

	/*
		Validation
	*/

	// Validate that the race exists
	var race *Race
	if v, ok := races[d.ID]; !ok {
		return
	} else {
		race = v
	}

	// Validate that the race has started
	if race.Status != "in progress" {
		return
	}

	// Validate that they are in the race
	var racer *Racer
	if v, ok := race.Racers[username]; !ok {
		return
	} else {
		racer = v
	}

	// Validate that they are still racing
	if racer.Status != "racing" {
		return
	}

	// Validate that the item number is sane
	// The the base game there are over 500 items and the Racing+ mod has a bunch of custom items
	// So just check for over 700 to be safe
	if itemID < 1 || itemID > 700 {
		log.Warning("User \"" + username + "\" attempted to add item " + strconv.Itoa(itemID) + " to their build, but that is a bogus number.")
		websocketError(s, d.Command, "That is not a valid item ID.")
		return
	}

	/*
		Add the item
	*/

	item := &Item{
		itemID,
		racer.FloorNum,
		racer.StageType,
		getTimestamp(),
	}
	racer.Items = append(racer.Items, item)

	// Check to see if this is their starting item
	startingItem := false
	if race.Ruleset.Format != "seeded" &&
		race.Ruleset.Format != "seeded-hard" &&
		racer.StartingItem == 0 &&
		len(racer.Rooms) > 1 {

		racer.StartingItem = itemID
		startingItem = true
	}

	for racerName := range race.Racers {
		// Not all racers may be online during a race
		if s, ok := websocketSessions[racerName]; ok {
			// Send the message about the item
			type RacerAddItemMessage struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
				Item *Item  `json:"item"`
			}
			websocketEmit(s, "racerAddItem", &RacerAddItemMessage{
				raceID,
				username,
				item,
			})

			if startingItem {
				websocketEmit(s, "racerSetStartingItem", &RacerAddItemMessage{
					raceID,
					username,
					item,
				})
			}
		}
	}
}
