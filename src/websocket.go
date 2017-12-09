package main

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/Zamiell/isaac-racing-server/src/log"
	"github.com/Zamiell/isaac-racing-server/src/models"
	melody "gopkg.in/olahol/melody.v1" // A WebSocket framework
)

var (
	// This is the Melody WebSocket router
	m *melody.Melody

	// We keep track of all WebSocket sessions
	websocketSessions = make(map[string]*melody.Session)

	// We keep track of which chat rooms exist and which users are in each room
	chatRooms = make(map[string][]User)

	// Used to store all of the functions that handle each command
	commandHandlerMap = make(map[string]func(*melody.Session, *IncomingWebsocketData))

	// The WebSocket server needs to processes one action at a time;
	// otherwise, there would be chaos
	commandMutex = new(sync.Mutex)
)

func websocketInit() {
	/*
		Define all of the WebSocket commands
	*/

	// Room (chat) commands
	commandHandlerMap["roomJoin"] = websocketRoomJoin
	commandHandlerMap["roomLeave"] = websocketRoomLeave
	commandHandlerMap["roomMessage"] = websocketRoomMessage
	commandHandlerMap["privateMessage"] = websocketPrivateMessage

	// Race commands
	commandHandlerMap["raceCreate"] = websocketRaceCreate
	commandHandlerMap["raceJoin"] = websocketRaceJoin
	commandHandlerMap["raceLeave"] = websocketRaceLeave
	commandHandlerMap["raceReady"] = websocketRaceReady
	commandHandlerMap["raceUnready"] = websocketRaceUnready
	commandHandlerMap["raceFinish"] = websocketRaceFinish
	commandHandlerMap["raceQuit"] = websocketRaceQuit
	commandHandlerMap["raceSeed"] = websocketRaceSeed
	commandHandlerMap["raceFloor"] = websocketRaceFloor
	commandHandlerMap["raceItem"] = websocketRaceItem
	commandHandlerMap["raceRoom"] = websocketRaceRoom

	// Profile commands
	commandHandlerMap["profileSetStream"] = websocketProfileSetStream

	// Admin commands
	commandHandlerMap["adminMessage"] = websocketAdminMessage
	commandHandlerMap["adminShutdown"] = websocketAdminShutdown
	commandHandlerMap["adminUnshutdown"] = websocketAdminUnshutdown
	commandHandlerMap["adminBan"] = websocketAdminBan
	commandHandlerMap["adminUnban"] = websocketAdminUnban
	/*
		commandHandlerMap["adminBanIP"] = websocketAdminBanIP
		commandHandlerMap["adminUnbanIP"] = websocketAdminUnbanIP
		commandHandlerMap["adminMute"] = websocketAdminMute
		commandHandlerMap["adminUnmute"] = websocketAdminUnmute
		commandHandlerMap["adminPromote"] = websocketAdminPromote
		commandHandlerMap["adminDemote"] = websocketAdminDemote

	*/

	// Miscellaneous commands
	commandHandlerMap["debug"] = websocketDebug

	// Define a new Melody router and attach a message handler
	m = melody.New()
	m.HandleConnect(websocketHandleConnect)
	m.HandleDisconnect(websocketHandleDisconnect)
	m.HandleMessage(websocketHandleMessage)
	// We could also attach a function to HandleError, but this fires on routine
	// things like disconnects, so it is undesirable
}

/*
	WebSocket miscellaneous subroutines
*/

// Get all values from the session and fill in the IncomingWebsocketData object
func websocketGetSessionValues(s *melody.Session, d *IncomingWebsocketData) bool {
	/*
		Get the values from the session
	*/

	var userID int
	if v, exists := s.Get("userID"); !exists {
		log.Error("Failed to get \"userID\" from the session (in the \"" + d.Command + "\" function).")
		return false
	} else {
		userID = v.(int)
	}

	var username string
	if v, exists := s.Get("username"); !exists {
		log.Error("Failed to get \"username\" from the session (in the \"" + d.Command + "\" function).")
		return false
	} else {
		username = v.(string)
	}

	var admin int
	if v, exists := s.Get("admin"); !exists {
		log.Error("Failed to get \"admin\" from the session (in the \"" + d.Command + "\" function).")
		return false
	} else {
		admin = v.(int)
	}

	var muted bool
	if v, exists := s.Get("muted"); !exists {
		log.Error("Failed to get \"muted\" from the session (in the \"" + d.Command + "\" function).")
		return false
	} else {
		muted = v.(bool)
	}

	var streamURL string
	if v, exists := s.Get("streamURL"); !exists {
		log.Error("Failed to get \"streamURL\" from the session (in the \"" + d.Command + "\" function).")
		return false
	} else {
		streamURL = v.(string)
	}

	var twitchBotEnabled bool
	if v, exists := s.Get("twitchBotEnabled"); !exists {
		log.Error("Failed to get \"twitchBotEnabled\" from the session (in the \"" + d.Command + "\" function).")
		return false
	} else {
		twitchBotEnabled = v.(bool)
	}

	var twitchBotDelay int
	if v, exists := s.Get("twitchBotDelay"); !exists {
		log.Error("Failed to get \"twitchBotDelay\" from the session (in the \"" + d.Command + "\" function).")
		return false
	} else {
		twitchBotDelay = v.(int)
	}

	var rateLimitAllowance float64
	if v, exists := s.Get("rateLimitAllowance"); !exists {
		log.Error("Failed to get \"rateLimitAllowance\" from the session (in the \"" + d.Command + "\" function).")
		return false
	} else {
		rateLimitAllowance = v.(float64)
	}

	var rateLimitLastCheck time.Time
	if v, exists := s.Get("rateLimitLastCheck"); !exists {
		log.Error("Failed to get \"rateLimitLastCheck\" from the session (in the \"" + d.Command + "\" function).")
		return false
	} else {
		rateLimitLastCheck = v.(time.Time)
	}

	/*
		Stick them inside the data object
	*/

	// "SessionValues" is defined in the "users.go" file
	d.v = &models.SessionValues{
		UserID:             userID,
		Username:           username,
		Admin:              admin,
		Muted:              muted,
		StreamURL:          streamURL,
		TwitchBotEnabled:   twitchBotEnabled,
		TwitchBotDelay:     twitchBotDelay,
		Banned:             false,
		RateLimitAllowance: rateLimitAllowance,
		RateLimitLastCheck: rateLimitLastCheck,
	}
	return true
}

// Send a message to a client using the Golem-style protocol described above
func websocketEmit(s *melody.Session, command string, d interface{}) {
	// Convert the data to JSON
	var ds string
	if dj, err := json.Marshal(d); err != nil {
		log.Error("Failed to marshal data when writing to a Melody session:", err)
		return
	} else {
		ds = string(dj)
	}

	// Send the message as bytes
	msg := command + " " + ds
	bytes := []byte(msg)
	if err := s.Write(bytes); err != nil {
		// This can routinely fail if the session is closed, so just return
		return
	}
}

// Used in the "error" and "warning" functions
type ErrorMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// Sent to the client if either their command was unsuccessful or something else went wrong
// (client-side, this will cause a WebSocket disconnect and the program to completely restart)
func websocketError(s *melody.Session, functionName string, message string) {
	if message == "" {
		// Specify a default error message
		message = "Something went wrong. Please contact an administrator."
	}
	websocketEmit(s, "error", &ErrorMessage{
		functionName,
		message,
	})
}

// Sent to the client if something unexpected happened
// (client-side, this will make a popup appear but still allow them to continue what they were doing)
func websocketWarning(s *melody.Session, functionName string, message string) {
	websocketEmit(s, "warning", &ErrorMessage{
		functionName,
		message,
	})
}

func websocketClose(s *melody.Session) {
	if err := s.Close(); err != nil {
		log.Error("Attempted to manually close a WebSocket connection, but it failed.")
	} else {
		log.Info("Successfully terminated a WebSocket connection.")
	}
}
