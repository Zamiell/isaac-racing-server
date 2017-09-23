package main

import (
	"path"

	"github.com/Zamiell/isaac-racing-server/src/log"
	"github.com/joho/godotenv"
	melody "gopkg.in/olahol/melody.v1"
)

func websocketAdminLoadEnv(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	username := d.v.Username
	admin := d.v.Admin

	// Validate that the user is an admin
	if admin != 2 {
		log.Warning("User \"" + username + "\" tried to send turn off the shutdown mode, but they are not an administrator.")
		websocketError(s, d.Command, "Only administrators can do that.")
		return
	}

	// Reload the ".env" file that we first loaded upon server initialization
	if err := godotenv.Load(path.Join(projectPath, ".env")); err != nil {
		log.Fatal("Failed to load .env file:", err)
	}
}
