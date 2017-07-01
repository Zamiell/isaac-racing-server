package main

/*
	Imports
*/

import (
	"github.com/Zamiell/isaac-racing-server/models"

	"math/rand"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

/*
	Global variables
*/

var validDiversityActiveItems = [...]int{
	// Rebirth items
	33, 34, 35, 36, 37, 38, 39, 40, 41, 42,
	44, 45, 47, 49, 56, 58, 65, 66, 77, 78,
	83, 84, 85, 86, 93, 97, 102, 105, 107, 111,
	123, 124, 126, 127, 130, 133, 135, 136, 137, 145,
	146, 147, 158, 160, 164, 166, 171, 175, 177, 181,
	186, 192, 282, 285, 286, 287, 288, 289, 290, 291, // D100 (283) and D4 (284) are banned
	292, 293, 294, 295, 296, 297, 298, 323, 324, 325,
	326, 338,

	// Afterbirth items
	347, 348, 349, 351, 352, 357, 382, 383, 386, 396,
	406, 419, 421, 422, 427, 434, 437, 439, 441,

	// Afterbirth+ items
	475, 476, 477, 478, 479, 480, 481, 482, 483, 484,
	485, 486, 487, 488, 490, 504, 507, 510, // D Infinity (489) is banned

	// Booster Pack items
	512, 515, 516, 521, 522, 523, 527,
}
var validDiversityPassiveItems = [...]int{
	// Rebirth items
	1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
	11, 12, 13, 14, 17, 18, 19, 20, 21, 27, // <3 (15), Raw Liver (16), Lunch (22), Dinner (23), Dessert (24), Breakfast (25), and Rotten Meat (26) are banned
	28, 32, 48, 50, 51, 52, 53, 54, 55, 57, // Mom's Underwear (29), Moms Heels (30), Moms Lipstick (31), and Lucky Foot (46) are banned
	60, 62, 63, 64, 67, 68, 69, 70, 71, 72,
	73, 74, 75, 76, 79, 80, 81, 82, 87, 88,
	89, 90, 91, 94, 95, 96, 98, 99, 100, 101, // Super Bandage (92) is banned
	103, 104, 106, 108, 109, 110, 112, 113, 114, 115,
	116, 117, 118, 119, 120, 121, 122, 125, 128, 129,
	131, 132, 134, 138, 139, 140, 141, 142, 143, 144,
	148, 149, 150, 151, 152, 153, 154, 155, 156, 157,
	159, 161, 162, 163, 165, 167, 168, 169, 170, 172,
	173, 174, 178, 179, 180, 182, 183, 184, 185, 187, // Stem Cells (176) is banned
	188, 189, 190, 191, 193, 195, 196, 197, 198, 199, // Magic 8 Ball (194) is banned
	200, 201, 202, 203, 204, 205, 206, 207, 208, 209,
	210, 211, 212, 213, 214, 215, 216, 217, 218, 219,
	220, 221, 222, 223, 224, 225, 227, 228, 229, 230, // Black Lotus (226) is banned
	231, 232, 233, 234, 236, 237, 240, 241, 242, 243, // Key Piece #1 (238) and Key Piece #2 (239) are banned
	244, 245, 246, 247, 248, 249, 250, 251, 252, 254, // Magic Scab (253) is banned
	255, 256, 257, 259, 260, 261, 262, 264, 265, 266, // Missing No. (258) is banned
	267, 268, 269, 270, 271, 272, 273, 274, 275, 276,
	277, 278, 279, 280, 281, 299, 300, 301, 302, 303,
	304, 305, 306, 307, 308, 309, 310, 311, 312, 313,
	314, 315, 316, 317, 318, 319, 320, 321, 322, 327,
	328, 329, 330, 331, 332, 333, 335, 336, 337, 340, // The Body (334) and Safety Pin (339) are banned
	341, 342, 343, 345, // Match Book (344) and A Snack (346) are banned

	// Afterbirth items
	350, 353, 354, 356, 358, 359, 360, 361, 362, 363, // Mom's Pearls (355) is banned
	364, 365, 366, 367, 368, 369, 370, 371, 372, 373,
	374, 375, 376, 377, 378, 379, 380, 381, 384, 385,
	387, 388, 389, 390, 391, 392, 393, 394, 395, 397,
	398, 399, 400, 401, 402, 403, 404, 405, 407, 408,
	409, 410, 411, 412, 413, 414, 415, 416, 417, 418,
	420, 423, 424, 425, 426, 429, 430, 431, 432, 433, // PJs (428) is banned
	435, 436, 438, 440,

	// Afterbirth+ items
	442, 443, 444, 445, 446, 447, 448, 449, 450, 451,
	452, 453, 454, 457, 458, 459, 460, 461, 462, 463, // Dad's Lost Coin (455) and Moldy Bread (456) are banned
	464, 465, 466, 467, 468, 469, 470, 471, 472, 473,
	474, 491, 492, 493, 494, 495, 496, 497, 498, 499,
	500, 501, 502, 503, 505, 506, 508, 509,

	// Booster pack items
	511, 513, 514, 517, 518, 519, 520, 524, 525, 526, 528, 529,
}
var validDiversityTrinkets = [...]int{
	// Rebirth trinkets
	1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
	11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
	21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
	31, 32, 33, 34, 35, 36, 37, 38, 39, 40,
	41, 42, 43, 44, 45, 46, 48, 49, 50, 51,
	52, 53, 54, 55, 56, 57, 58, 59, 60, 61,

	// Afterbirth trinkets
	62, 63, 64, 65, 66, 67, 68, 69, 70, 71,
	72, 73, 74, 75, 76, 77, 78, 79, 80, 81,
	82, 83, 84, 86, 87, 88, 89, 90, // Karma (85) is banned

	// Afterbirth+ trinkets
	91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
	101, 102, 103, 104, 105, 106, 107, 108, 109, 110,
	111, 112, 113, 114, 115, 116, 117, 118, 119,

	// Booster pack trinkets
	120, 121, 122, 123,
}

/*
	WebSocket race command functions
*/

func raceCreate(conn *ExtendedConnection, data *IncomingCommandMessage) {
	// Local variables
	functionName := "raceCreate"
	userID := conn.UserID
	username := conn.Username
	name := data.Name
	ruleset := data.Ruleset

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()

	// Log the received command
	log.Debug("User \"" + username + "\" sent a " + functionName + " command.") // There is no race ID yet because they are creating a race

	// Rate limit all commands
	if commandRateLimit(conn) == true {
		return
	}

	// Validate that the race name cannot be empty
	if name == "" {
		name = "-"
	}

	// Validate that the race name is not longer than 100 characters
	if utf8.RuneCountInString(name) > 100 {
		log.Warning("User \"" + username + "\" sent a race name longer than 100 characters.")
		commandMutex.Unlock()
		connError(conn, functionName, "Race names must not be longer than 100 characters.")
		return
	}

	// Validate that the ruleset options cannot be empty
	if ruleset.Type == "" {
		ruleset.Type = "unranked"
	}
	if ruleset.Format == "" {
		ruleset.Format = "unseeded"
	}
	if ruleset.Character == "" {
		ruleset.Character = "Judas"
	}
	if ruleset.Goal == "" {
		ruleset.Goal = "The Chest"
	}

	// Validate the submitted ruleset
	if raceValidateRuleset(conn, data, functionName) == false {
		return
	}

	// Check if there are any non-finished races with the same name
	if raceWithSameName, err := db.Races.CheckName(name); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	} else if raceWithSameName == true {
		commandMutex.Unlock()
		connError(conn, functionName, "There is already a non-finished race with that name.")
		return
	}

	// Create the race
	raceID, err := db.Races.Insert(name, ruleset, userID)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Create and set a seed if necessary
	var seed string
	if ruleset.Format == "seeded" || ruleset.Format == "diversity" {
		if ruleset.Format == "seeded" {
			// Create a random Isaac seed
			outputByteArray, err := exec.Command(projectPath + "/scripts/get_seed").Output()
			if err != nil {
				commandMutex.Unlock()
				log.Error("Failed to execute the \"get_seed\" binary:", err)
				connError(conn, functionName, "Something went wrong. Please contact an administrator.")
				return
			}
			seed = strings.TrimSpace(string(outputByteArray)) // There is a newline at the end, so we have to remove it
		} else if ruleset.Format == "diversity" {
			// Get 1 random unique active item
			var items []int
			rand.Seed(time.Now().UnixNano())
			item := validDiversityActiveItems[rand.Intn(len(validDiversityActiveItems))]
			items = append(items, item)

			// Get 3 random unique passive items
			for i := 1; i <= 3; i++ {
				for {
					// Initialize the PRNG and get a random element from the slice
					// (if we don't do this, it will use a seed of 1)
					rand.Seed(time.Now().UnixNano())
					item := validDiversityPassiveItems[rand.Intn(len(validDiversityPassiveItems))]

					// Ensure this item is unique
					if intInSlice(item, items) == false {
						items = append(items, item)
						break
					}
				}
			}

			// Get 1 random trinket
			rand.Seed(time.Now().UnixNano())
			trinket := validDiversityTrinkets[rand.Intn(len(validDiversityTrinkets))]
			items = append(items, trinket)

			// The "seed" value is used to communicate the 5 random diversity items to the client
			for _, item := range items {
				seed += strconv.Itoa(item) + ","
			}
			seed = strings.TrimSuffix(seed, ",") // Remove the trailing comma
		}

		// Set the new seed
		if err = db.Races.SetSeed(raceID, seed); err != nil {
			commandMutex.Unlock()
			log.Error("Database error:", err)
			connError(conn, functionName, "Something went wrong. Please contact an administrator.")
			return
		}
	} else {
		seed = "-"
	}

	// Add this user to the participants list for that race
	if err = db.RaceParticipants.Insert(userID, raceID); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Send everyone a notification that a new race has been started
	for _, conn := range connectionMap.m {
		conn.Connection.Emit("raceCreated", &models.Race{
			ID:              raceID,
			Name:            name,
			Status:          "open",
			Ruleset:         ruleset,
			Seed:            seed,
			DatetimeCreated: makeTimestamp(),
			Captain:         username,
			Racers:          []string{username},
		})
	}

	// Get all the information about the racers in this race
	racerList, err := db.RaceParticipants.GetRacerList(raceID)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		return
	}

	// Send it to the user
	conn.Connection.Emit("racerList", &RacerList{raceID, racerList})

	// Join the user to the channel for that race
	roomJoinSub(conn, "_race_"+strconv.Itoa(raceID))

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()
}

func raceJoin(conn *ExtendedConnection, data *IncomingCommandMessage) {
	// Local variables
	functionName := "raceJoin"
	userID := conn.UserID
	username := conn.Username
	raceID := data.ID

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()

	// Log the received command
	log.Debug("User \"" + username + "\" sent a " + functionName + " command for race ID: " + strconv.Itoa(raceID))

	// Rate limit all commands
	if commandRateLimit(conn) == true {
		return
	}

	// Validate basic things about the race ID
	if raceValidate(conn, data, functionName) == false {
		return
	}

	// Validate that the race is open
	if raceValidateStatus(conn, data, "open", functionName) == false {
		return
	}

	// Validate that they are not in the race
	if raceValidateOut2(conn, data, functionName) == false {
		return
	}

	// Validate that this is not a solo race
	if solo, err := db.Races.CheckSolo(raceID); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	} else if solo {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" attempted to call " + functionName + " on race ID " + strconv.Itoa(raceID) + ", but it is a solo race.")
		connError(conn, functionName, "Race ID "+strconv.Itoa(raceID)+" is a solo race, so you cannot join it.")
		return
	}

	// Add this user to the participants list for that race
	if err := db.RaceParticipants.Insert(userID, raceID); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Send everyone a notification that the user joined
	for _, conn := range connectionMap.m {
		conn.Connection.Emit("raceJoined", RaceMessage{raceID, username})
	}

	// Get all the information about the racers in this race
	racerList, err := db.RaceParticipants.GetRacerList(raceID)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		return
	}

	// Send it to the user
	conn.Connection.Emit("racerList", &RacerList{raceID, racerList})

	// Join the user to the channel for that race
	roomJoinSub(conn, "_race_"+strconv.Itoa(raceID))

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()
}

func raceJoinSpectate(conn *ExtendedConnection, data *IncomingCommandMessage) {
	// Local variables
	functionName := "raceJoinSpectate"
	userID := conn.UserID
	username := conn.Username
	raceID := data.ID

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()

	// Log the received command
	log.Debug("User \"" + username + "\" sent a " + functionName + " command for race ID: " + strconv.Itoa(raceID))

	// Rate limit all commands
	if commandRateLimit(conn) == true {
		return
	}

	// Validate basic things about the race ID
	if raceValidate(conn, data, functionName) == false {
		return
	}

	// Validate that they are not in the race
	if raceValidateOut2(conn, data, functionName) == false {
		return
	}

	// Add this user to the participants list for that race
	if err := db.RaceParticipants.Insert(userID, raceID); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Send everyone a notification that the user joined
	for _, conn := range connectionMap.m {
		conn.Connection.Emit("raceJoined", RaceMessage{raceID, username})
	}

	// Get all the information about the racers in this race
	racerList, err := db.RaceParticipants.GetRacerList(raceID)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		return
	}

	// Send it to the user
	conn.Connection.Emit("racerList", &RacerList{raceID, racerList})

	// Join the user to the channel for that race
	roomJoinSub(conn, "_race_"+strconv.Itoa(raceID))

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()
}

func raceLeave(conn *ExtendedConnection, data *IncomingCommandMessage) {
	// Local variables
	functionName := "raceLeave"
	username := conn.Username
	raceID := data.ID

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()

	// Log the received command
	log.Debug("User \"" + username + "\" sent a " + functionName + " command for race ID: " + strconv.Itoa(raceID))

	// Rate limit all commands
	if commandRateLimit(conn) == true {
		return
	}

	// Validate basic things about the race ID
	if raceValidate(conn, data, functionName) == false {
		return
	}

	// Validate that the race is open
	if raceValidateStatus(conn, data, "open", functionName) == false {
		return
	}

	// Validate that they are in the race
	if raceValidateIn2(conn, data, functionName) == false {
		return
	}

	// Remove this user from the participants list for that race
	if err := db.RaceParticipants.Delete(username, raceID); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Disconnect the user from the channel for that race
	roomLeaveSub(conn, "_race_"+strconv.Itoa(raceID))

	// Send everyone a notification that the user left the race
	for _, conn := range connectionMap.m {
		conn.Connection.Emit("raceLeft", RaceMessage{raceID, username})
	}

	// If the race went from 2 people to 1, automatically unready the last person so that they don't start the race by themsevles
	if racerNames, err := db.RaceParticipants.GetRacerNames(raceID); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	} else if len(racerNames) == 1 {
		if currentStatus, err := db.RaceParticipants.GetStatus(racerNames[0], raceID); err != nil {
			commandMutex.Unlock()
			log.Error("Database error:", err)
			connError(conn, functionName, "Something went wrong. Please contact an administrator.")
			return
		} else if currentStatus == "ready" {
			// Set them from "ready" to "not ready"
			if err := db.RaceParticipants.SetStatus(racerNames[0], raceID, "not ready"); err != nil {
				commandMutex.Unlock()
				log.Error("Database error:", err)
				connError(conn, functionName, "Something went wrong. Please contact an administrator.")
				return
			}

			// Tell them
			conn, ok := connectionMap.m[racerNames[0]]
			if ok == true { // They should definately be online, but check anyway just in case
				conn.Connection.Emit("racerSetStatus", &RacerSetStatusMessage{raceID, username, "not ready", 0})
			}
		}
	}

	// Check to see if the race is ready to start
	raceCheckStart(raceID)

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()
}

func raceReady(conn *ExtendedConnection, data *IncomingCommandMessage) {
	// Local variables
	functionName := "raceReady"
	userID := conn.UserID
	username := conn.Username
	raceID := data.ID

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()

	// Log the received command
	log.Debug("User \"" + username + "\" sent a " + functionName + " command for race ID: " + strconv.Itoa(raceID))

	// Rate limit all commands
	if commandRateLimit(conn) == true {
		return
	}

	// Validate basic things about the race ID
	if raceValidate(conn, data, functionName) == false {
		return
	}

	// Validate that the race is open
	if raceValidateStatus(conn, data, "open", functionName) == false {
		return
	}

	// Validate that they are in the race
	if raceValidateIn2(conn, data, functionName) == false {
		return
	}

	// Validate that their status is set to "not ready"
	if racerValidateStatus(conn, userID, raceID, "not ready", functionName) == false {
		return
	}

	// Change their status to "ready"
	if racerSetStatus(conn, username, raceID, "ready", functionName) == false {
		return
	}

	// Get the list of racers for this race
	racerNames, err := db.RaceParticipants.GetRacerNames(raceID)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		return
	}

	// Send a notification to all the people in this particular race that the user is ready
	for _, racer := range racerNames {
		conn, ok := connectionMap.m[racer]
		if ok == true { // Not all racers may be online during a race
			conn.Connection.Emit("racerSetStatus", &RacerSetStatusMessage{raceID, username, "ready", 0})
		}
	}

	// Check to see if the race is ready to start
	raceCheckStart(raceID)

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()
}

func raceUnready(conn *ExtendedConnection, data *IncomingCommandMessage) {
	// Local variables
	functionName := "raceUnready"
	userID := conn.UserID
	username := conn.Username
	raceID := data.ID

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()

	// Log the received command
	log.Debug("User \"" + username + "\" sent a " + functionName + " command for race ID: " + strconv.Itoa(raceID))

	// Rate limit all commands
	if commandRateLimit(conn) == true {
		return
	}

	// Validate basic things about the race ID
	if raceValidate(conn, data, functionName) == false {
		return
	}

	// Validate that the race is open
	if raceValidateStatus(conn, data, "open", functionName) == false {
		return
	}

	// Validate that they are in the race
	if raceValidateIn2(conn, data, functionName) == false {
		return
	}

	// Validate that their status is set to "ready"
	if racerValidateStatus(conn, userID, raceID, "ready", functionName) == false {
		return
	}

	// Change their status to "not ready"
	if racerSetStatus(conn, username, raceID, "not ready", functionName) == false {
		return
	}

	// Get the list of racers for this race
	racerNames, err := db.RaceParticipants.GetRacerNames(raceID)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		return
	}

	// Send a notification to all the people in this particular race that the user is not ready
	for _, racer := range racerNames {
		conn, ok := connectionMap.m[racer]
		if ok == true { // Not all racers may be online during a race
			conn.Connection.Emit("racerSetStatus", &RacerSetStatusMessage{raceID, username, "not ready", 0})
		}
	}

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()
}

// Currently not implemented client-side
/*
func raceRuleset(conn *ExtendedConnection, data *IncomingCommandMessage) {
	// Local variables
	functionName := "raceUnready"
	userID := conn.UserID
	username := conn.Username
	raceID := data.ID
	ruleset := data.Ruleset

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()

	// Log the received command
	log.Debug("User \"" + username + "\" sent a " + functionName + " command for race ID: " + strconv.Itoa(raceID))

	// Rate limit all commands
	if commandRateLimit(conn) == true {
		return
	}

	// Get the current ruleset
	currentRuleset, err := db.Races.GetRuleset(raceID)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Check to see if anything has changed
	if currentRuleset.Format == ruleset.Format &&
		currentRuleset.Character == ruleset.Character &&
		currentRuleset.Goal == ruleset.Goal &&
		currentRuleset.StartingBuild == ruleset.StartingBuild {

		commandMutex.Unlock()
		connError(conn, functionName, "The race ruleset is already set to those values.")
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
	if raceValidateRuleset(conn, data, functionName) == false {
		return
	}

	// Validate basic things about the race ID
	if raceValidate(conn, data, functionName) == false {
		return
	}

	// Validate that the race is open
	if raceValidateStatus(conn, data, "open", functionName) == false {
		return
	}

	// Validate that they are in the race
	if raceValidateIn2(conn, data, functionName) == false {
		return
	}

	// Validate that they are the race captain
	if isCaptain, err := db.Races.CheckCaptain(raceID, userID); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	} else if isCaptain == false {
		commandMutex.Unlock()
		connError(conn, functionName, "Only the captain of the race can change the ruleset.")
		return
	}

	// Get and set a seed if necessary
	if (ruleset.Format == "seeded" || ruleset.Format == "diversity") && ruleset.Format != currentRuleset.Format {
		var seed string
		if ruleset.Format == "seeded" {
			// Create a random Isaac seed
			outputByteArray, err := exec.Command(projectPath + "/scripts/get_seed").Output()
			if err != nil {
				commandMutex.Unlock()
				log.Error("Failed to execute the \"get_seed\" binary:", err)
				connError(conn, functionName, "Something went wrong. Please contact an administrator.")
				return
			}
			seed = string(outputByteArray)
		} else if ruleset.Format == "diversity" {
			// Get 3 random unique passive items
			var items []int
			for i := 1; i <= 3; i++ {
				for {
					// Initialize the PRNG and get a random element from the slice
					// (if we don't do this, it will use a seed of 1)
					rand.Seed(time.Now().UnixNano())
					item := validDiversityItems[rand.Intn(len(validDiversityItems))]

					// Ensure this item is unique
					if intInSlice(item, items) == false {
						items = append(items, item)
						break
					}
				}
			}

			// The "seed" value is used to communicate the 3 random diversity items to the client
			for _, item := range items {
				seed += strconv.Itoa(item) + ","
			}
			seed = strings.TrimSuffix(seed, ",") // Remove the trailing comma
		}

		// Set the new seed
		if err := db.Races.SetSeed(raceID, seed); err != nil {
			commandMutex.Unlock()
			log.Error("Database error:", err)
			connError(conn, functionName, "Something went wrong. Please contact an administrator.")
			return
		}

	}

	// Change the ruleset
	if err := db.Races.SetRuleset(raceID, ruleset); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Set everyone's status to "not ready"
	if err := db.RaceParticipants.SetAllStatus(raceID, "not ready"); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Send everyone a notification that the ruleset has changed for this race
	for _, conn := range connectionMap.m {
		conn.Connection.Emit("raceSetRuleset", RaceSetRulesetMessage{raceID, ruleset})
	}

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()
}
*/

func raceFinish(conn *ExtendedConnection, data *IncomingCommandMessage) {
	// Local variables
	functionName := "raceFinish"
	userID := conn.UserID
	username := conn.Username
	raceID := data.ID

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()

	// Log the received command
	log.Debug("User \"" + username + "\" sent a " + functionName + " command for race ID: " + strconv.Itoa(raceID))

	// Rate limit all commands
	if commandRateLimit(conn) == true {
		return
	}

	// Validate basic things about the race ID
	if raceValidate(conn, data, functionName) == false {
		return
	}

	// Validate that the race has started
	if raceValidateStatus(conn, data, "in progress", functionName) == false {
		return
	}

	// Validate that they are in the race
	if raceValidateIn2(conn, data, functionName) == false {
		return
	}

	// Validate that their status is set to "racing" status
	if racerValidateStatus(conn, userID, raceID, "racing", functionName) == false {
		return
	}

	// Change their status to "finished"
	if racerSetStatus(conn, username, raceID, "finished", functionName) == false {
		return
	}

	// Set their finish time
	if err := db.RaceParticipants.SetDatetimeFinished(username, raceID, int(makeTimestamp())); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Get the place of the last person that finished so far
	currentPlace, err := db.RaceParticipants.GetCurrentPlace(raceID)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Set their (final) place
	if err = db.RaceParticipants.SetPlace(username, raceID, currentPlace+1); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Get the list of racers for this race
	racerList, err := db.RaceParticipants.GetRacerList(raceID)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Recalculate everyones mid-race places
	if racerSetAllPlaceMid(conn, raceID, racerList, functionName) == false {
		return
	}

	// Send a notification to all the people in this particular race that the user finished
	for _, racer := range racerList {
		racerConn, ok := connectionMap.m[racer.Name]
		if ok == true { // Not all racers may be online during a race
			racerConn.Connection.Emit("racerSetStatus", &RacerSetStatusMessage{raceID, username, "finished", currentPlace + 1})
		}
	}

	// Calculate their run time
	started, err := db.Races.GetDatetimeStarted(raceID)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		return
	}
	var runTime int
	var place int
	for _, racer := range racerList {
		if racer.Name == username {
			runTime = racer.DatetimeFinished - started
			place = racer.Place
			break
		}
	}
	minutes := strconv.Itoa(runTime / 1000 / 60)
	seconds := strconv.Itoa(runTime / 1000 % 60)
	if len(seconds) == 1 {
		seconds = "0" + seconds
	}
	timeString := "(" + minutes + ":" + seconds + ")"
	placeString := getOrdinal(place)

	// Get the number of people left in the race
	peopleLeft, err := db.RaceParticipants.GetPeopleLeft(raceID)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Make the Twitch bot announce that the person finished
	twitchString := "/me - " + placeString + " - " + username + " " + timeString + " - "
	if peopleLeft == 0 {
		twitchString += "Race completed."
	} else {
		twitchString += strconv.Itoa(peopleLeft) + " left"
	}
	for _, racer := range racerList {
		twitchRacerSend(racer, twitchString)
	}

	// Check to see if the race is ready to finish
	raceCheckFinish(raceID)

	// Update fields in the users table (e.g. average, ELO)
	// (we calculate stats for seeded races only when the race is completed)
	raceUpdateUnseededStats(raceID, username)

	// Check to see if the user got any achievements
	achievementsCheck(username)

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()
}

func raceQuit(conn *ExtendedConnection, data *IncomingCommandMessage) {
	// Local variables
	functionName := "raceQuit"
	userID := conn.UserID
	username := conn.Username
	raceID := data.ID

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()

	// Log the received command
	log.Debug("User \"" + username + "\" sent a " + functionName + " command for race ID: " + strconv.Itoa(raceID))

	// Rate limit all commands
	if commandRateLimit(conn) == true {
		return
	}

	// Validate basic things about the race ID
	if raceValidate(conn, data, functionName) == false {
		return
	}

	// Validate that the race has started
	if raceValidateStatus(conn, data, "in progress", functionName) == false {
		return
	}

	// Validate that they are in the race
	if raceValidateIn2(conn, data, functionName) == false {
		return
	}

	// Validate that their status is set to "racing" status
	if racerValidateStatus(conn, userID, raceID, "racing", functionName) == false {
		return
	}

	// Change their status to "quit"
	if racerSetStatus(conn, username, raceID, "quit", functionName) == false {
		return
	}

	// Set their finish time
	if err := db.RaceParticipants.SetDatetimeFinished(username, raceID, int(makeTimestamp())); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Set their (final) place to -1 (which indicates a quit status)
	if err := db.RaceParticipants.SetPlace(username, raceID, -1); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Get the list of racers for this race
	racerList, err := db.RaceParticipants.GetRacerList(raceID)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		return
	}

	// Recalculate everyones mid-race places
	if racerSetAllPlaceMid(conn, raceID, racerList, functionName) == false {
		return
	}

	// Send a notification to all the people in this particular race that the user quit
	for _, racer := range racerList {
		racerConn, ok := connectionMap.m[racer.Name]
		if ok == true { // Not all racers may be online during a race
			racerConn.Connection.Emit("racerSetStatus", &RacerSetStatusMessage{raceID, username, "quit", -1})
		}
	}

	// Get the number of people left in the race
	peopleLeft, err := db.RaceParticipants.GetPeopleLeft(raceID)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Make the Twitch bot announce that the person quit
	twitchString := "/me - " + username + " quit - "
	if peopleLeft == 0 {
		twitchString += "Race completed."
	} else {
		twitchString += strconv.Itoa(peopleLeft) + " left"
	}
	for _, racer := range racerList {
		twitchRacerSend(racer, twitchString)
	}

	// Check to see if the race is ready to finish
	raceCheckFinish(raceID)

	// Update fields in the users table (e.g. average, ELO)
	// (we calculate stats for seeded races only when the race is completed)
	raceUpdateUnseededStats(raceID, username)

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()
}

func raceComment(conn *ExtendedConnection, data *IncomingCommandMessage) {
	// Local variables
	functionName := "raceQuit"
	userID := conn.UserID
	username := conn.Username
	raceID := data.ID
	comment := data.Comment

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()

	// Log the received command
	log.Debug("User \"" + username + "\" sent a " + functionName + " command for race ID: " + strconv.Itoa(raceID))

	// Rate limit all commands
	if commandRateLimit(conn) == true {
		return
	}

	// Strip leading and trailing whitespace from the comment
	comment = strings.TrimSpace(comment)

	// Validate that the comment is not empty
	if comment == "" {
		commandMutex.Unlock()
		connError(conn, functionName, "That is an invalid comment.")
		return
	}

	// Validate that the comment is not excessively long
	if len(comment) < 150 {
		commandMutex.Unlock()
		connError(conn, functionName, "Comments must not be longer than 150 characters.")
		return
	}

	// Validate that the comment does not contain special characters
	if isAlphaNumericUnderscore(comment) == false {
		commandMutex.Unlock()
		connError(conn, functionName, "Your comment must contain only letters, numbers, and underscores.")
		return
	}

	// Validate that the user is not muted
	if conn.Muted == 1 {
		commandMutex.Unlock()
		connError(conn, functionName, "You have been muted by an administrator, so you cannot submit comments.")
		return
	}

	// Validate basic things about the race ID
	if raceValidate(conn, data, functionName) == false {
		return
	}

	// Validate that the race has started
	if raceValidateStatus(conn, data, "in progress", functionName) == false {
		return
	}

	// Validate that they are in the race
	if raceValidateIn2(conn, data, functionName) == false {
		return
	}

	// Set their comment in the database
	if err := db.RaceParticipants.SetComment(userID, raceID, comment); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Get the list of racers for this race
	racerNames, err := db.RaceParticipants.GetRacerNames(raceID)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		return
	}

	// Send a notification to all the people in this particular race that the user added or changed their comment
	for _, racer := range racerNames {
		conn, ok := connectionMap.m[racer]
		if ok == true { // Not all racers may be online during a race
			conn.Connection.Emit("racerSetComment", &RacerSetCommentMessage{raceID, username, comment})
		}
	}

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()
}

func raceItem(conn *ExtendedConnection, data *IncomingCommandMessage) {
	// Local variables
	functionName := "raceItem"
	userID := conn.UserID
	username := conn.Username
	raceID := data.ID
	itemID := data.ItemID

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()

	// Log the received command
	log.Debug("User \"" + username + "\" sent a " + functionName + " command for race ID: " + strconv.Itoa(raceID))

	// Rate limit all commands
	if commandRateLimit(conn) == true {
		return
	}

	// Validate that the item number is sane
	// The highest item ID is 510 (Delirious), and the Racing+ mod has a bunch of custom items
	// So just check for over 600 to be safe
	if itemID < 1 || itemID > 600 {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" attempted to add item " + strconv.Itoa(itemID) + " to their build, but that is a bogus number.")
		connError(conn, functionName, "That is not a valid item ID.")
		return
	}

	// Validate basic things about the race ID
	if raceValidate(conn, data, functionName) == false {
		return
	}

	// Validate that the race has started
	if raceValidateStatus(conn, data, "in progress", functionName) == false {
		return
	}

	// Validate that they are in the race
	if raceValidateIn2(conn, data, functionName) == false {
		return
	}

	// Validate that their status is set to "racing" status
	if racerValidateStatus(conn, userID, raceID, "racing", functionName) == false {
		return
	}

	// Get their current floor
	floorNum, stageType, err := db.RaceParticipants.GetFloor(userID, raceID)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Add this item to their build
	if err = db.RaceParticipantItems.Insert(userID, raceID, itemID, floorNum, stageType); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Get the list of racers for this race
	racerNames, err := db.RaceParticipants.GetRacerNames(raceID)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		return
	}

	log.Debug("Getting here 1.")

	// Send a notification to all the people in this particular race that the user got an item
	for _, racer := range racerNames {
		conn, ok := connectionMap.m[racer]
		if ok == true { // Not all racers may be online during a race
			item := models.Item{itemID, floorNum, stageType}
			conn.Connection.Emit("racerAddItem", &RacerAddItemMessage{raceID, username, item})
		}
	}

	log.Debug("Getting here 2.")

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()
}

func raceFloor(conn *ExtendedConnection, data *IncomingCommandMessage) {
	// Local variables
	functionName := "raceFloor"
	userID := conn.UserID
	username := conn.Username
	raceID := data.ID
	floorNum := data.FloorNum
	stageType := data.StageType

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()

	// Log the received command
	log.Debug("User \"" + username + "\" sent a " + functionName + " command for race ID: " + strconv.Itoa(raceID))

	// Rate limit all commands
	if commandRateLimit(conn) == true {
		return
	}

	// Validate basic things about the race ID
	if raceValidate(conn, data, functionName) == false {
		return
	}

	// Validate that the race has started
	if raceValidateStatus(conn, data, "in progress", functionName) == false {
		return
	}

	// Validate that they are in the race
	if raceValidateIn2(conn, data, functionName) == false {
		return
	}

	// Validate that their status is set to "racing" status
	if racerValidateStatus(conn, userID, raceID, "racing", functionName) == false {
		return
	}

	// Validate that the floor is sane
	if floorNum < 1 || floorNum > 12 {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" attempted to update their floor, but \"" + strconv.Itoa(floorNum) + "\" is a bogus floor number.")
		connError(conn, functionName, "That is not a valid floor number.")
		return
	} else if stageType < 0 || stageType > 3 {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" attempted to update their floor, but \"" + strconv.Itoa(stageType) + "\" is a bogus stage type.")
		connError(conn, functionName, "That is not a valid stage type.")
		return
	}

	// Set their floor in the database
	floorArrived, err := db.RaceParticipants.SetFloor(userID, raceID, floorNum, stageType)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// The floor gets sent as 1 when a reset occurs
	if floorNum == 1 {
		// Reset all of their accumulated items
		if err = db.RaceParticipantItems.Reset(userID, raceID); err != nil {
			commandMutex.Unlock()
			log.Error("Database error:", err)
			connError(conn, functionName, "Something went wrong. Please contact an administrator.")
			return
		}

		// Reset all of their visited rooms
		if err = db.RaceParticipantItems.Reset(userID, raceID); err != nil {
			commandMutex.Unlock()
			log.Error("Database error:", err)
			connError(conn, functionName, "Something went wrong. Please contact an administrator.")
			return
		}
	}

	// Get the list of racers for this race
	racerList, err := db.RaceParticipants.GetRacerList(raceID)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		return
	}

	// Recalculate everyones mid-race places
	if racerSetAllPlaceMid(conn, raceID, racerList, functionName) == false {
		return
	}

	log.Debug("Getting here 1.")

	// Send a notification to all the people in this particular race that the user got to a new floor
	for _, racer := range racerList {
		conn, ok := connectionMap.m[racer.Name]
		if ok == true { // Not all racers may be online during a race
			conn.Connection.Emit("racerSetFloor", &RacerSetFloorMessage{raceID, username, floorNum, stageType, floorArrived})
		}
	}

	log.Debug("Getting here 2.")

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()
}

func raceRoom(conn *ExtendedConnection, data *IncomingCommandMessage) {
	// Local variables
	functionName := "raceRoom"
	userID := conn.UserID
	username := conn.Username
	raceID := data.ID
	roomID := data.RoomID

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()

	// Log the received command
	log.Debug("User \"" + username + "\" sent a " + functionName + " command for race ID: " + strconv.Itoa(raceID))

	// Rate limit all commands
	if commandRateLimit(conn) == true {
		return
	}

	// Validate basic things about the race ID
	if raceValidate(conn, data, functionName) == false {
		return
	}

	// Validate that the race has started
	if raceValidateStatus(conn, data, "in progress", functionName) == false {
		return
	}

	// Validate that they are in the race
	if raceValidateIn2(conn, data, functionName) == false {
		return
	}

	// Validate that their status is set to "racing" status
	if racerValidateStatus(conn, userID, raceID, "racing", functionName) == false {
		return
	}

	// Add the room to their list of visited rooms
	if err := db.RaceParticipantRooms.Insert(userID, raceID, roomID); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()
}

/*
	Race subroutines
*/

func raceValidate(conn *ExtendedConnection, data *IncomingCommandMessage, functionName string) bool {
	// Local variables
	username := conn.Username
	raceID := data.ID

	// Validate that the requested race is sane
	if raceID <= 0 {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" attempted to call " + functionName + " with a bogus ID of " + strconv.Itoa(raceID) + ".")
		connError(conn, functionName, "You must provide a valid race number.")
		return false
	}

	// Validate that the requested race exists
	if exists, err := db.Races.Exists(raceID); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return false
	} else if exists == false {
		commandMutex.Unlock()
		log.Info("User \"" + username + "\" attempted to call " + functionName + " on race ID " + strconv.Itoa(raceID) + ", but it doesn't exist.")
		// Don't send an error to the user; this kind of thing can happen if their internet is lagging
		return false
	}

	// The user's request seems to be valid
	return true
}

func raceValidateStatus(conn *ExtendedConnection, data *IncomingCommandMessage, status string, functionName string) bool {
	// Local variables
	username := conn.Username
	raceID := data.ID

	// Validate that the race is set to the correct status
	if correctStatus, err := db.Races.CheckStatus(raceID, status); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return false
	} else if correctStatus == false {
		commandMutex.Unlock()
		log.Info("User \"" + username + "\" attempted to call " + functionName + " on race ID " + strconv.Itoa(raceID) + ", but race is not set to status \"" + status + "\".")
		// Don't send an error to the user; this kind of thing can happen if their internet is lagging
		return false
	}

	// The race is the correct status
	return true
}

func raceValidateRuleset(conn *ExtendedConnection, data *IncomingCommandMessage, functionName string) bool {
	// Local variables
	ruleset := data.Ruleset

	// Validate the ruleset type
	if ruleset.Type != "ranked" &&
		ruleset.Type != "unranked" {

		commandMutex.Unlock()
		connError(conn, functionName, "That is not a valid type.")
		return false
	}

	// Validate the ruleset format
	if ruleset.Format != "unseeded" &&
		ruleset.Format != "seeded" &&
		ruleset.Format != "diversity" &&
		ruleset.Format != "custom" {

		commandMutex.Unlock()
		connError(conn, functionName, "That is not a valid ruleset.")
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
		ruleset.Character != "custom" {

		commandMutex.Unlock()
		connError(conn, functionName, "That is not a valid character.")
		return false
	}

	// Validate the goal
	if ruleset.Goal != "Blue Baby" &&
		ruleset.Goal != "The Lamb" &&
		ruleset.Goal != "Mega Satan" &&
		ruleset.Goal != "custom" {

		commandMutex.Unlock()
		connError(conn, functionName, "That is not a valid goal.")
		return false
	}

	// Validate the starting build
	if ruleset.Format != "seeded" && ruleset.StartingBuild != -1 {
		commandMutex.Unlock()
		connError(conn, functionName, "You cannot set a starting build for a non-seeded race.")
		return false
	} else if ruleset.Format == "seeded" && (ruleset.StartingBuild < 1 || ruleset.StartingBuild > 32) { // There are 32 builds currently
		commandMutex.Unlock()
		connError(conn, functionName, "That is not a valid starting build.")
		return false
	}

	// Validate unseeded ranked games
	if ruleset.Type == "ranked" && ruleset.Format == "unseeded" && ruleset.Character != "Judas" {
		commandMutex.Unlock()
		connError(conn, functionName, "Ranked unseeded races must have a character of Judas.")
		return false
	}
	if ruleset.Type == "ranked" && ruleset.Format == "unseeded" && ruleset.Goal != "Blue Baby" {
		commandMutex.Unlock()
		connError(conn, functionName, "Ranked unseeded races must have a goal of Blue Baby.")
		return false
	}

	return true
}

// Playing or observing
func raceValidateIn1(conn *ExtendedConnection, data *IncomingCommandMessage, functionName string) bool {
	// Local variables
	userID := conn.UserID
	username := conn.Username
	raceID := data.ID

	// Validate that they are in the race
	if userInRace, err := db.RaceParticipants.CheckInRace1(userID, raceID); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return false
	} else if userInRace == false {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" attempted to call " + functionName + " on race ID " + strconv.Itoa(raceID) + ", but they are not in that race.")
		connError(conn, functionName, "You are not playing in or observing race ID "+strconv.Itoa(raceID)+".")
		return false
	}

	// The user is in the race
	return true
}

// ONLY playing (not observing)
func raceValidateIn2(conn *ExtendedConnection, data *IncomingCommandMessage, functionName string) bool {
	// Local variables
	userID := conn.UserID
	username := conn.Username
	raceID := data.ID

	// Validate that they are in the race
	if userInRace, err := db.RaceParticipants.CheckInRace2(userID, raceID); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return false
	} else if userInRace == false {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" attempted to call " + functionName + " on race ID " + strconv.Itoa(raceID) + ", but they are not in that race.")
		connError(conn, functionName, "You are not playing in race ID "+strconv.Itoa(raceID)+".")
		return false
	}

	// The user is in the race
	return true
}

func raceValidateOut1(conn *ExtendedConnection, data *IncomingCommandMessage, functionName string) bool {
	// Local variables
	userID := conn.UserID
	username := conn.Username
	raceID := data.ID

	// Validate that they are not already in the race
	if userInRace, err := db.RaceParticipants.CheckInRace1(userID, raceID); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return false
	} else if userInRace == true {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" attempted to call " + functionName + " on race ID " + strconv.Itoa(raceID) + ", but they are already in that race.")
		connError(conn, functionName, "You are already playing or observing race ID "+strconv.Itoa(raceID)+".")
		return false
	}

	// The user is not in the race
	return true
}

func raceValidateOut2(conn *ExtendedConnection, data *IncomingCommandMessage, functionName string) bool {
	// Local variables
	userID := conn.UserID
	username := conn.Username
	raceID := data.ID

	// Validate that they are not already in the race
	if userInRace, err := db.RaceParticipants.CheckInRace2(userID, raceID); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return false
	} else if userInRace == true {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" attempted to call " + functionName + " on race ID " + strconv.Itoa(raceID) + ", but they are already in that race.")
		connError(conn, functionName, "You are already playing in race ID "+strconv.Itoa(raceID)+".")
		return false
	}

	// The user is not in the race
	return true
}

func racerValidateStatus(conn *ExtendedConnection, userID int, raceID int, status string, functionName string) bool {
	// Local variables
	username := conn.Username

	// Validate that the user is set to the correct status
	if correctStatus, err := db.RaceParticipants.CheckStatus(userID, raceID, status); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return false
	} else if correctStatus == false {
		commandMutex.Unlock()
		log.Info("User \"" + username + "\" attempted to call " + functionName + " on race ID " + strconv.Itoa(raceID) + ", but they are not set to status \"" + status + "\".")
		// Don't send an error to the user; just silently fail
		// This type of thing can occur if, for example, they try to unready immediately before the race begins
		return false
	}

	// The user has the correct status
	return true
}

func racerSetStatus(conn *ExtendedConnection, username string, raceID int, status string, functionName string) bool {
	// Change the status in the database
	if err := db.RaceParticipants.SetStatus(username, raceID, status); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return false
	}

	// The change was successful
	return true
}

// Recalculate everyones mid-race places
func racerSetAllPlaceMid(conn *ExtendedConnection, raceID int, racerList []models.Racer, functionName string) bool {
	// Get the current (final) place
	var currentPlace int
	for _, racer := range racerList {
		if racer.Place > currentPlace {
			currentPlace = racer.Place
		}
	}

	// Recalculate everyones mid-race places
	for _, racer := range racerList {
		if racer.Status != "racing" {
			continue // We don't need to calculate the mid-race place of someone who already finished or quit
		}
		racer.PlaceMid = currentPlace + 1
		for _, racer2 := range racerList {
			if racer.Status != "racing" {
				continue // We don't count people who finished or quit since our starting point was on the currentPlace
			}
			if racer2.FloorNum > racer.FloorNum {
				racer.PlaceMid++
			} else if racer2.FloorNum == racer.FloorNum && racer2.FloorArrived < racer.FloorArrived {
				racer.PlaceMid++
			}
		}
		if err := db.RaceParticipants.SetPlaceMid(racer.Name, raceID, racer.PlaceMid); err != nil {
			commandMutex.Unlock()
			log.Error("Database error:", err)
			connError(conn, functionName, "Something went wrong. Please contact an administrator.")
			return false
		}
	}

	// Everything was set successfully
	return true
}

// Called after someone disconnects or someone is banned
// (the commandMutex should be locked when getting here)
func raceCheckStartFinish(raceID int) {
	// Get the status of the race
	if status, err := db.Races.GetStatus(raceID); err != nil {
		log.Error("Database error:", err)
		return
	} else if status == "open" {
		raceCheckStart(raceID)
	} else if status == "in progress" {
		raceCheckFinish(raceID)
	}
}

// Check to see if a race is ready to start, and if so, start it
// (the commandMutex should be locked when getting here)
func raceCheckStart(raceID int) {
	// Check to see if the race was deleted
	if exists, err := db.Races.Exists(raceID); err != nil {
		log.Error("Database error:", err)
		return
	} else if exists == false {
		return
	}

	// Check to see if this is a solo race
	solo, err2 := db.Races.CheckSolo(raceID)
	if err2 != nil {
		log.Error("Database error:", err2)
		return
	}

	// Check to see if there is only 1 person in the race
	if racerNames, err := db.RaceParticipants.GetRacerNames(raceID); err != nil {
		log.Error("Database error:", err)
		return
	} else if len(racerNames) == 1 && solo == false {
		return
	}

	// Check if everyone is ready
	if sameStatus, err := db.RaceParticipants.CheckAllStatus(raceID, "ready"); err != nil {
		log.Error("Database error:", err)
		return
	} else if sameStatus == false {
		return
	}

	// Log the race starting
	log.Debug("Race " + strconv.Itoa(raceID) + " starting in 10 seconds.")

	// Change the status for this race to "starting"
	if err := db.Races.SetStatus(raceID, "starting"); err != nil {
		log.Error("Database error:", err)
		return
	}

	// Send everyone a notification that the race is starting soon
	for _, conn := range connectionMap.m {
		conn.Connection.Emit("raceSetStatus", &RaceSetStatusMessage{raceID, "starting"})
	}

	// Get the list of people in this race
	racerList, err := db.RaceParticipants.GetRacerList(raceID)
	if err != nil {
		log.Error("Database error:", err)
		return
	}

	// Get the time X seconds in the future
	var secondsToWait time.Duration
	if solo {
		secondsToWait = 3
	} else {
		secondsToWait = 10
	}
	startTime := time.Now().Add(secondsToWait*time.Second).UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))

	// Send everyone in the race a message describing exactly when it will start
	for _, racer := range racerList {
		conn, ok := connectionMap.m[racer.Name]
		if ok == true {
			conn.Connection.Emit("raceStart", &RaceStartMessage{
				raceID,
				startTime,
			})
		} else {
			log.Warning("Failed to send a raceStart message to user \"" + racer.Name + "\". This should never happen.")
		}
	}

	// Make the Twitch bot announce that the race is starting in 10 seconds
	if solo == false {
		for _, racer := range racerList {
			twitchRacerSend(racer, "/me - The race is starting in 10 seconds!")
		}
	}

	// Return for now and do more things in 10 seconds
	go raceCheckStart2(raceID)
}

func raceCheckStart2(raceID int) {
	// Check to see if this is a solo race
	solo, err2 := db.Races.CheckSolo(raceID)
	if err2 != nil {
		log.Error("Database error:", err2)
		return
	}

	// Sleep 3 or 10 seconds
	var sleepTime time.Duration
	if solo {
		sleepTime = 3
	} else {
		sleepTime = 10
	}
	time.Sleep(sleepTime * time.Second)

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()

	// Check to see if the race was deleted
	if exists, err := db.Races.Exists(raceID); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		return
	} else if exists == false {
		commandMutex.Unlock()
		return
	}

	// Get the amount of people in this race
	racerNames, err := db.RaceParticipants.GetRacerNames(raceID)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		return
	}

	// Log the race starting
	log.Info("Race", raceID, "started with", len(racerNames), "participants:", racerNames)

	// Change the status for this race to "in progress" and set "datetime_started" equal to now
	if err := db.Races.Start(raceID); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		return
	}

	// Update the status for everyone in the race to "racing"
	if err := db.RaceParticipants.SetAllStatus(raceID, "racing"); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		return
	}

	// Send everyone a notification that the race is now in progress
	for _, conn := range connectionMap.m {
		conn.Connection.Emit("raceSetStatus", &RaceSetStatusMessage{raceID, "in progress"})
	}

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()

	// Return for now and do more things in 30 minutes
	go raceCheckStart3(raceID)
}

func raceCheckStart3(raceID int) {
	// Sleep 30 minutes
	time.Sleep(30 * time.Minute)

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()

	// Find out if the race is finished
	if status, err := db.Races.GetStatus(raceID); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		return
	} else if status == "finished" {
		commandMutex.Unlock()
		return
	}

	// The race is still going, so get the list of people still in this race
	racerList, err := db.RaceParticipants.GetRacerList(raceID)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		return
	}

	// If any are still racing, force them to quit
	for _, racer := range racerList {
		if racer.Status == "racing" {
			if err := db.RaceParticipants.SetStatus(racer.Name, raceID, "quit"); err != nil {
				commandMutex.Unlock()
				log.Error("Database error:", err)
				return
			}

			// Send a notification to all the people in this particular race that the user quit
			for _, racer2 := range racerList {
				conn, ok := connectionMap.m[racer2.Name]
				if ok == true { // Not all racers may be online during a race
					conn.Connection.Emit("racerSetStatus", &RacerSetStatusMessage{raceID, racer.Name, "quit", -1})
				}
			}
		}
	}

	// Close down the race
	raceCheckFinish(raceID)

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()
}

// Check to see if a rate is ready to finish, and if so, finish it
// (the commandMutex should be locked when getting here)
func raceCheckFinish(raceID int) {
	// Check to see if the race was deleted
	if exists, err := db.Races.Exists(raceID); err != nil {
		log.Error("Database error:", err)
		return
	} else if exists == false {
		return
	}

	// Check if anyone is still racing
	if stillRacing, err := db.RaceParticipants.CheckStillRacing(raceID); err != nil {
		log.Error("Database error:", err)
		return
	} else if stillRacing == true {
		return
	}

	// Log the race finishing
	log.Info("Race " + strconv.Itoa(raceID) + " finished.")

	// Change the status for this race to "finished" and set "datetime_finished" equal to now
	if err := db.Races.Finish(raceID); err != nil {
		log.Error("Database error:", err)
		return
	}

	// Send everyone a notification that the race is now finished
	for _, conn := range connectionMap.m {
		conn.Connection.Emit("raceSetStatus", &RaceSetStatusMessage{raceID, "finished"})
	}
}

// Now that a user has finished, quit, or been disqualified from a race, update fields in the users table for unseeded races
func raceUpdateUnseededStats(raceID int, username string) {
	// Don't do anything if this is not an unseeded race (or an unranked race)
	/*
		if unseededAndRanked, err := db.Races.CheckUnseededRanked(raceID); err != nil {
			commandMutex.Unlock()
			log.Error("Database error:", err)
			return
		} else if unseededAndRanked == false {
			return
		}
	*/

	// Get their unseeded stats
	/*
		if statsUnseeded, err := db.Users.GetStatsUnseeded(username); err != nil {
			commandMutex.Unlock()
			log.Error("Database error:", err)
			return
		}
	*/

	// Update all the stats

}

// Now that the race has finished, update fields in the users table for seeded races
func raceUpdateSeededStats(raceID int, username string) {

}
