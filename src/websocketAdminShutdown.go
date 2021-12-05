package server

import (
	"os/exec"
	"path"
	"time"

	melody "gopkg.in/olahol/melody.v1"
)

func websocketAdminShutdown(s *melody.Session, d *IncomingWebsocketData) {
	username := d.v.Username
	admin := d.v.Admin
	comment := d.Comment

	// Validate that the user is an admin
	if admin != 2 {
		logger.Warning("User \"" + username + "\" tried to send turn on the shutdown mode, but they are not an administrator.")
		websocketError(s, d.Command, "Only administrators can do that.")
		return
	}

	if comment == "restart" {
		// We want to automatically restart the server once all races are finished
		shutdownMode = 1
	} else {
		// We don't want to automatically restart the server and have to probably perform some manual administrative activity
		shutdownMode = 2
	}

	if len(races) > 0 {
		d.Message = "The server will restart when all ongoing races have finished. New race creation has been disabled."
		websocketAdminMessage(s, d)
		go websocketAdminShutdownSub(s, d)
	} else {
		restartServer(s, d)
	}
}

func websocketAdminShutdownSub(s *melody.Session, d *IncomingWebsocketData) {
	for {
		time.Sleep(time.Second)

		if shutdownMode == 0 {
			logger.Info("shutdownMode changed to false. Shutdown aborted.")
			break
		}

		// Check to see if all races are finished
		if len(races) == 0 {
			// Wait 30 seconds so that the last people finishing a race are not immediately booted upon finishing
			time.Sleep(time.Second * 30)

			restartServer(s, d)
			break
		}
	}
}

func restartServer(s *melody.Session, d *IncomingWebsocketData) {
	d.Message = "The server is restarting; please stand by."
	websocketAdminMessage(s, d)

	if shutdownMode == 2 {
		// We don't want to perform an automatic restart
		return
	}

	cmd := exec.Command(path.Join(projectPath, "restart.sh")) // nolint: gosec
	if output, err := cmd.Output(); err != nil {
		logger.Error("Failed to execute \"restart.sh\":", err)
	} else {
		logger.Info("\"restart.sh\" completed:", string(output))
	}
}
