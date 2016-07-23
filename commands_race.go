package main

/*
 *  Imports
 */

import (
	"github.com/Zamiell/isaac-racing-server/models"

	"strconv"
	"time"
)

/*
 *  WebSocket race command functions
 */

func raceCreate(conn *ExtendedConnection, data *RaceCreateMessage) {
	// Local variables
	functionName := "raceCreate"
	userID := conn.UserID
	username := conn.Username
	name := data.Name
	ruleset := data.Ruleset

	// Log the received command
	log.Debug("User \""+username+"\" sent a", functionName, "command.")

	// Rate limit all commands
	if commandRateLimit(conn) == true {
		return
	}

	// Validate that the race name cannot be empty
	if name == "" {
		name = "-"
	}

	// Validate that the ruleset cannot be empty
	if ruleset == "" {
		ruleset = "unseeded"
	}

	// Validate the ruleset
	if ruleset != "unseeded" && ruleset != "seeded" && ruleset != "diversity" {
		connError(conn, functionName, "That is not a valid ruleset.")
		return
	}

	// Check if this user has started 2 races
	if count, err := db.Races.CaptainCount(userID); err != nil {
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	} else if count >= 2 {
		log.Info("New race request denied; user is captain of ", count, "races.")
		connError(conn, functionName, "To prevent abuse, you are only allowed to create 2 new races at a time.")
		return
	}

	// Check if there are non-finished races with the same name
	if raceWithSameName, err := db.Races.CheckName(name); err != nil {
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	} else if raceWithSameName == true {
		connError(conn, functionName, "There is already a non-finished race with that name.")
		return
	}

	// Create the race (and add this user to the participants list for that race)
	raceID, err := db.Races.Insert(userID, name)
	if err != nil {
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Send success confirmation
	data.ID = raceID // Requested by Chronometrics as extra information for the client
	connSuccess(conn, functionName, data)

	// Join the user to the channel for that race
	roomJoinSub(conn, "_race_"+strconv.Itoa(raceID))

	// Send everyone the new list of races
	raceUpdateAll()

	// Send the people in this race an update
	racerUpdate(raceID)
}

func raceJoin(conn *ExtendedConnection, data *RaceMessage) {
	// Local variables
	functionName := "raceJoin"
	userID := conn.UserID
	username := conn.Username
	raceID := data.ID

	// Log the received command
	log.Debug("User \""+username+"\" sent a", functionName, "command.")

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
	if raceValidateOut(conn, data, functionName) == false {
		return
	}

	// Add this user to the participants list for that race
	if err := db.RaceParticipants.Insert(userID, raceID); err != nil {
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Get the name of this race
	name, err := db.Races.GetName(raceID)
	if err != nil {
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Join the user to the channel for that race
	roomJoinSub(conn, "_race_"+strconv.Itoa(raceID))

	// Send success confirmation
	data.Name = name // Requested by Chronometrics as extra information for the client
	connSuccess(conn, functionName, data)

	// Send everyone the new list of races
	raceUpdateAll()

	// Send the people in this race an update
	racerUpdate(raceID)
}

func raceLeave(conn *ExtendedConnection, data *RaceMessage) {
	// Local variables
	functionName := "raceLeave"
	userID := conn.UserID
	username := conn.Username
	raceID := data.ID

	// Log the received command
	log.Debug("User \""+username+"\" sent a", functionName, "command.")

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
	if raceValidateIn(conn, data, functionName) == false {
		return
	}

	// Remove this user from the participants list for that race
	if err := db.RaceParticipants.Delete(userID, raceID); err != nil {
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Disconnect the user from the channel for that race
	roomLeaveSub(conn, "_race_"+strconv.Itoa(raceID))

	// Send success confirmation
	connSuccess(conn, functionName, data)

	// Send everyone the new list of races
	raceUpdateAll()

	// Send the people in this race an update
	racerUpdate(raceID)

	// Check to see if the race is ready to start
	raceStart(raceID)
}

func raceReady(conn *ExtendedConnection, data *RaceMessage) {
	// Local variables
	functionName := "raceReady"
	userID := conn.UserID
	username := conn.Username
	raceID := data.ID

	// Log the received command
	log.Debug("User \""+username+"\" sent a", functionName, "command.")

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
	if raceValidateIn(conn, data, functionName) == false {
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

	// Send success confirmation
	connSuccess(conn, functionName, data)

	// Send the people in this race an update
	racerUpdate(raceID)

	// Check to see if the race is ready to start
	raceStart(raceID)
}

func raceUnready(conn *ExtendedConnection, data *RaceMessage) {
	// Local variables
	functionName := "raceUnready"
	userID := conn.UserID
	username := conn.Username
	raceID := data.ID

	// Log the received command
	log.Debug("User \""+username+"\" sent a", functionName, "command.")

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
	if raceValidateIn(conn, data, functionName) == false {
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

	// Send success confirmation
	connSuccess(conn, functionName, data)

	// Send the people in this race an update
	racerUpdate(raceID)
}

func raceRuleset(conn *ExtendedConnection, data *RaceRulesetMessage) {
	// Local variables
	functionName := "raceUnready"
	userID := conn.UserID
	username := conn.Username
	raceID := data.ID
	ruleset := data.Ruleset

	// Log the received command
	log.Debug("User \""+username+"\" sent a", functionName, "command.")

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
	if raceValidateIn(conn, data, functionName) == false {
		return
	}

	// Validate that they are the race captain
	if isCaptain, err := db.Races.CheckCaptain(raceID, userID); err != nil {
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	} else if isCaptain == false {
		connError(conn, functionName, "Only the captain of the race can change the ruleset.")
		return
	}

	// Validate the ruleset
	if ruleset != "unseeded" && ruleset != "seeded" && ruleset != "diversity" {
		connError(conn, functionName, "That is not a valid ruleset.")
		return
	}

	// Change the ruleset
	if err := db.Races.SetRuleset(raceID, ruleset); err != nil {
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Set everyone's status to "not ready"
	if err := db.RaceParticipants.SetAllStatus(raceID, "not ready"); err != nil {
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Send success confirmation
	connSuccess(conn, functionName, data)

	// Send everyone the new list of races
	raceUpdateAll()

	// Send the people in this race an update
	racerUpdate(raceID)
}

func raceDone(conn *ExtendedConnection, data *RaceMessage) {
	// Local variables
	functionName := "raceDone"
	userID := conn.UserID
	username := conn.Username
	raceID := data.ID

	// Log the received command
	log.Debug("User \""+username+"\" sent a", functionName, "command.")

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
	if raceValidateIn(conn, data, functionName) == false {
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

	// Send success confirmation
	connSuccess(conn, functionName, data)

	// Send the people in this race an update
	racerUpdate(raceID)

	// Check to see if the race is ready to finish
	raceFinish(raceID)
}

func raceQuit(conn *ExtendedConnection, data *RaceMessage) {
	// Local variables
	functionName := "raceQuit"
	userID := conn.UserID
	username := conn.Username
	raceID := data.ID

	// Log the received command
	log.Debug("User \""+username+"\" sent a", functionName, "command.")

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
	if raceValidateIn(conn, data, functionName) == false {
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

	// Send success confirmation
	connSuccess(conn, functionName, data)

	// Send the people in this race an update
	racerUpdate(raceID)

	// Check to see if the race is ready to finish
	raceFinish(raceID)
}

func raceComment(conn *ExtendedConnection, data *RaceCommentMessage) {
	// Local variables
	functionName := "raceQuit"
	userID := conn.UserID
	username := conn.Username
	raceID := data.ID
	comment := data.Comment

	// Log the received command
	log.Debug("User \""+username+"\" sent a", functionName, "command.")

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
	if raceValidateIn(conn, data, functionName) == false {
		return
	}

	// Set their comment in the database
	if err := db.RaceParticipants.SetComment(userID, raceID, comment); err != nil {
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Send success confirmation
	connSuccess(conn, functionName, data)

	// Send the people in this race an update
	racerUpdate(raceID)
}

func raceItem(conn *ExtendedConnection, data *RaceItemMessage) {
	// Local variables
	functionName := "raceItem"
	userID := conn.UserID
	username := conn.Username
	raceID := data.ID
	itemID := data.ItemID

	// Log the received command
	log.Debug("User \""+username+"\" sent a", functionName, "command.")

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
	if raceValidateIn(conn, data, functionName) == false {
		return
	}

	// Validate that their status is set to "racing" status
	if racerValidateStatus(conn, userID, raceID, "racing", functionName) == false {
		return
	}

	// Validate that the item number is sane
	if itemID < 1 || itemID > 441 { // This will need to be updated once we know the highest item ID in Afterbirth+
		log.Warning("User \""+username+"\" attempted to add an item", itemID, "to their build, but that is a bogus number.")
		connError(conn, functionName, "That is not a valid item ID.")
		return
	}

	// Add this item to their build
	if err := db.RaceParticipantItems.Insert(userID, raceID, itemID); err != nil {
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Send success confirmation
	connSuccess(conn, functionName, data)

	// Send the people in this race an update
	racerUpdate(raceID)
}

func raceFloor(conn *ExtendedConnection, data *RaceFloorMessage) {
	// Local variables
	functionName := "raceFloor"
	userID := conn.UserID
	username := conn.Username
	raceID := data.ID
	floor := data.Floor

	// Log the received command
	log.Debug("User \""+username+"\" sent a", functionName, "command.")

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
	if raceValidateIn(conn, data, functionName) == false {
		return
	}

	// Validate that their status is set to "racing" status
	if racerValidateStatus(conn, userID, raceID, "racing", functionName) == false {
		return
	}

	// Validate that the floor is sane
	if floor < 1 || floor > 10 {
		log.Warning("User \""+username+"\" attempted to update their floor, but", floor, "is a bogus number.")
		connError(conn, functionName, "That is not a valid floor.")
		return
	}

	// Set their floor in the database
	if err := db.RaceParticipants.SetFloor(userID, raceID, floor); err != nil {
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Send success confirmation
	connSuccess(conn, functionName, data)

	// Send the people in this race an update
	racerUpdate(raceID)
}

/*
 *  Race subroutines
 */

func raceValidate(conn *ExtendedConnection, data interface{}, functionName string) bool {
	// Local variables
	username := conn.Username
	raceID := data.(*RaceMessage).ID

	// Validate that the requested race is sane
	if raceID <= 0 {
		log.Warning("User \""+username+"\" attempted to call", functionName, "with a bogus ID of "+strconv.Itoa(raceID)+".")
		connError(conn, functionName, "You must provide a valid race number.")
		return false
	}

	// Validate that the requested race exists
	if exists, err := db.Races.Exists(raceID); err != nil {
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return false
	} else if exists == false {
		log.Warning("User \""+username+"\" attempted to call", functionName, "on race ID "+strconv.Itoa(raceID)+", but it doesn't exist.")
		connError(conn, functionName, "Race ID "+strconv.Itoa(raceID)+" does not exist.")
		return false
	}

	// The user's request seems to be valid
	return true
}

func raceValidateStatus(conn *ExtendedConnection, data interface{}, status string, functionName string) bool {
	// Local variables
	username := conn.Username
	raceID := data.(*RaceMessage).ID

	// Validate that the race is set to the correct status
	if correctStatus, err := db.Races.CheckStatus(raceID, status); err != nil {
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return false
	} else if correctStatus == false {
		log.Warning("User \""+username+"\" attempted to call", functionName, "on race ID "+strconv.Itoa(raceID)+", but race is not set to status \""+status+"\".")
		connError(conn, functionName, "Race ID "+strconv.Itoa(raceID)+" is not set to status \""+status+"\".")
		return false
	}

	// The race is the correct status
	return true
}

func raceValidateIn(conn *ExtendedConnection, data interface{}, functionName string) bool {
	// Local variables
	userID := conn.UserID
	username := conn.Username
	raceID := data.(*RaceMessage).ID

	// Validate that they are in the race
	if userInRace, err := db.RaceParticipants.CheckInRace(userID, raceID); err != nil {
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return false
	} else if userInRace == false {
		log.Warning("User \""+username+"\" attempted to call", functionName, "on race ID "+strconv.Itoa(raceID)+", but they are not in that race.")
		connError(conn, functionName, "You are not in race ID "+strconv.Itoa(raceID)+".")
		return false
	}

	// The user is in the race
	return true
}

func raceValidateOut(conn *ExtendedConnection, data interface{}, functionName string) bool {
	// Local variables
	userID := conn.UserID
	username := conn.Username
	raceID := data.(*RaceMessage).ID

	// Validate that they are not already in the race
	if userInRace, err := db.RaceParticipants.CheckInRace(userID, raceID); err != nil {
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return false
	} else if userInRace == true {
		log.Warning("User \""+username+"\" attempted to call", functionName, "on race ID "+strconv.Itoa(raceID)+", but they are already in that race.")
		connError(conn, functionName, "You are already in race ID "+strconv.Itoa(raceID)+".")
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
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return false
	} else if correctStatus == false {
		log.Warning("User \""+username+"\" attempted to call", functionName, "on race ID "+strconv.Itoa(raceID)+", but they are not set to status \""+status+"\".")
		connError(conn, functionName, "You can only do that if your status is set to \""+status+"\".")
		return false
	}

	// The user has the correct status
	return true
}

func racerSetStatus(conn *ExtendedConnection, username string, raceID int, status string, functionName string) bool {
	// Change the status in the database
	if err := db.RaceParticipants.SetStatus(username, raceID, status); err != nil {
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return false
	}

	// The change was successful
	return true
}

// Called when a client initially connects
func raceUpdate(conn *ExtendedConnection) {
	// Local variables
	functionName := "raceUpdate"

	// Get the current races
	var raceList []model.Race
	raceList, err := db.Races.GetCurrentRaces()
	if err != nil {
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Send it to the user
	conn.Connection.Emit("raceList", raceList)
}

// Called whenever someone joins or leaves a race, a race changes status, or a race changes ruleset
func raceUpdateAll() {
	// Get the current races
	var raceList []model.Race
	raceList, err := db.Races.GetCurrentRaces()
	if err != nil {
		return
	}

	// Send it to all users
	connectionMap.RLock()
	for _, conn := range connectionMap.m {
		conn.Connection.Emit("raceList", raceList)
	}
	connectionMap.RUnlock()
}

// Called whenever someone does something inside of a race
func racerUpdate(raceID int) {
	// Get the list of racers for this race
	var racerList []model.Racer
	racerList, err := db.RaceParticipants.GetRacerList(raceID)
	if err != nil {
		return
	}

	// Send it to all the people in this particular race
	connectionMap.RLock()
	for _, racer := range racerList {
		username := racer.Name
		conn, ok := connectionMap.m[username]
		if ok == true { // Not all racers may be online during a running race
			conn.Connection.Emit("racerList", &RacerList{
				raceID,
				racerList,
			})
		}
	}
	connectionMap.RUnlock()
}

// Called after a user leaves a race, someone readies, or someone quits
func raceCheckStartFinish(raceID int) {
	// Get the status of the race
	if status, err := db.Races.GetStatus(raceID); err != nil {
		return
	} else if status == "open" {
		go raceStart(raceID) // Need to use a goroutine since this function uses sleeps
	} else if status == "in progress" {
		raceFinish(raceID)
	}
}

// Check to see if a race is ready to start, and if so, start it
func raceStart(raceID int) {
	// Check if everyone is ready
	if sameStatus, err := db.RaceParticipants.CheckAllStatus(raceID, "ready"); err != nil {
		return
	} else if sameStatus == false {
		return
	}

	// Log the race starting
	log.Info("Race " + strconv.Itoa(raceID) + " started.")

	// Change the status for this race to "starting"
	if err := db.Races.SetStatus(raceID, "starting"); err != nil {
		return
	}

	// Send everyone the new list of races
	raceUpdateAll()

	// Get the list of people in this race
	racers, err := db.RaceParticipants.GetRacerNames(raceID)
	if err != nil {
		return
	}

	// Send everyone in the race a message describing exactly when it will start
	connectionMap.RLock()
	for _, username := range racers {
		conn, ok := connectionMap.m[username]
		if ok == true {
			conn.Connection.Emit("raceStart", &RaceStartMessage{
				raceID,
				time.Now().Add(10 * time.Second).UnixNano(), // 10 seconds in the future
			})
		} else {
			log.Warning("Failed to send a raceStart message to user \"" + username + "\". This should never happen.")
		}
	}
	connectionMap.RUnlock()

	// Sleep 10 seconds
	time.Sleep(10 * time.Second)

	// Start the race (which will set everyone's status to "racing")
	if err := db.Races.Start(raceID); err != nil {
		return
	}

	// Send everyone the new list of races
	raceUpdateAll()

	// Send the people in this race an update
	racerUpdate(raceID)

	// Sleep 30 minutes
	time.Sleep(30 * time.Minute)

	// Find out if the race is finished
	if status, err := db.Races.GetStatus(raceID); err != nil {
		return
	} else if status == "finished" {
		return
	}

	// The race is still going, so get the list of people still in this race
	racerList, err := db.RaceParticipants.GetRacerList(raceID)
	if err != nil {
		return
	}

	// If any are still racing, force them to quit
	for _, racer := range racerList {
		if racer.Status == "racing" {
			if err := db.RaceParticipants.SetStatus(racer.Name, raceID, "quit"); err != nil {
				return
			}
		}
	}

	// Close down the race
	raceFinish(raceID)
}

// Check to see if a rate is ready to finish, and if so, finish it
func raceFinish(raceID int) {
	// Check if anyone is still racing
	if stillRacing, err := db.RaceParticipants.CheckStillRacing(raceID); err != nil {
		return
	} else if stillRacing == true {
		return
	}

	// Log the race finishing
	log.Info("Race " + strconv.Itoa(raceID) + " finished.")

	// Finish the race
	if err := db.Races.Finish(raceID); err != nil {
		return
	}

	// Send everyone the new list of races
	raceUpdateAll()
}
