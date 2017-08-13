package main

import (
	"os/exec"
	"path"
	"time"

	"github.com/Zamiell/isaac-racing-server/src/log"
	melody "gopkg.in/olahol/melody.v1"
)

func websocketAdminShutdown(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	username := d.v.Username
	admin := d.v.Admin

	// Validate that the user is an admin
	if admin != 2 {
		log.Warning("User \"" + username + "\" tried to send turn on the shutdown mode, but they are not an administrator.")
		websocketError(s, d.Command, "Only administrators can do that.")
		return
	}

	shutdownMode = true
	d.Message = "The server is restarting soon. New race creation has been disabled for the time being."
	websocketAdminMessage(s, d)

	go websocketAdminShutdownSub(s, d)
}

func websocketAdminShutdownSub(s *melody.Session, d *IncomingWebsocketData) {
	for {
		time.Sleep(time.Second)

		if !shutdownMode {
			log.Info("shutdownMode changed to false. Shutdown aborted.")
			break
		}

		// Check to see if all races are finished
		if len(races) == 0 {
			// Wait 30 seconds so that the last people finishing a race are not immediately booted upon finishing
			time.Sleep(time.Second * 30)

			d.Message = "All races have completed. Initiating shutdown and restart."
			websocketAdminMessage(s, d)
			restart()
			break
		}
	}
}

func restart() {
	exec.Command(path.Join(projectPath, "restart.sh"))
}
