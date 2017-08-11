package main

/*
	Data types that are shared between different elements of the program
*/

// Used when passing the cookie values from "httpValidateSession" to "httpWS"
// Also used in "websocketHandleMessage" to stuff WebSocket session values into
// the WebsocketData object as a convience for command handler functions
type SessionValues struct {
	UserID    int
	Username  string
	Admin     int
	Muted     bool
	StreamURL string
}
