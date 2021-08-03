package server

import (
	"strconv"
	"strings"

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
	if race.Status != RaceStatusInProgress {
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
	if racer.Status != RacerStatusRacing {
		return
	}

	// Validate that the item number is sane
	// The the base game there are over 500 items and the Racing+ mod has a bunch of custom items
	// Furthermore, we hardcode some custom items in the 3000-3999 range
	if itemID < 1 || itemID > 4000 {
		logger.Warning("User \"" + username + "\" attempted to add item " + strconv.Itoa(itemID) + " to their build, but that is a bogus number.")
		websocketError(s, d.Command, "That is not a valid item ID.")
		return
	}

	// Custom items are handled manually
	// The final vanilla item is Mom's Shovel (552)
	// The first custom item is 3001
	if itemID > 552 && itemID <= 3000 && itemID != 560 {
		return
	}

	// Validate that this is not a More Options (414) that is given for Basement 1 only
	if itemID == 414 && len(racer.Rooms) == 1 && race.Ruleset.Character != "Eden" {
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

	if itemID == 560 { // Checkpoint
		racer.CharacterNum++
		race.SetAllPlaceMid()

		for racerName := range race.Racers {
			// Not all racers may be online during a race
			if s, ok := websocketSessions[racerName]; ok {
				// Send the message about the new character
				type RacerCharacterMessage struct {
					ID           int    `json:"id"`
					Name         string `json:"name"`
					CharacterNum int    `json:"characterNum"`
				}
				websocketEmit(s, "racerCharacter", &RacerCharacterMessage{
					ID:           raceID,
					Name:         racer.Name,
					CharacterNum: racer.CharacterNum,
				})
			}
		}
	}

	// Check to see if this is their starting item
	startingItem := false
	if race.Ruleset.Format != RaceFormatSeeded &&
		racer.StartingItem == 0 &&
		len(racer.Rooms) > 1 {

		// Every character starts with the D6
		startedWithThis := false
		if itemID == 105 {
			startedWithThis = true
		}

		// For Diversity races, check to see if this item was already given to them at the start
		if race.Ruleset.Format == RaceFormatDiversity {
			for i, startingItem := range strings.Split(race.Ruleset.Seed, ",") {
				if i == 4 {
					// We don't want to compare to the trinket
					continue
				}
				if startingItemInt, err := strconv.Atoi(startingItem); err != nil {
					logger.Error("Failed to parse the Diversity seed when checking for the starting item:", err)
					continue
				} else {
					if itemID == startingItemInt {
						startedWithThis = true
						break
					}
				}
			}
		}

		if !startedWithThis {
			racer.StartingItem = itemID
			startingItem = true
		}
	}

	for racerName := range race.Racers {
		// Not all racers may be online during a race
		if s, ok := websocketSessions[racerName]; ok {
			// Send the message about the item
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
