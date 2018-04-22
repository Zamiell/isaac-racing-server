package main

import (
	"fmt"
	"strconv"

	"github.com/Zamiell/isaac-racing-server/src/log"
	melody "gopkg.in/olahol/melody.v1"
)

func websocketDebug(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	username := d.v.Username
	admin := d.v.Admin

	// Validate that the user is an admin
	if admin == 0 {
		log.Info("User \"" + username + "\" tried to do a debug command, but they are not staff/admin.")
		websocketError(s, d.Command, "Only staff members or administrators can do that.")
		return
	}

	var i int

	// Print out all of the connected users
	log.Debug("----------------")
	log.Debug("Connected users:")
	i = 0
	for name := range websocketSessions {
		i++
		log.Debug(strconv.Itoa(i)+")", name)
	}

	// Print out all of the races
	log.Debug("--------------")
	log.Debug("Ongoing races:")
	i = 0
	for _, race := range races {
		i++
		raceString := fmt.Sprintf("%+v", race)
		log.Debug(strconv.Itoa(i)+")", raceString)
		j := 0
		for _, racer := range race.Racers {
			j++
			racerString := fmt.Sprintf("%+v", racer)
			log.Debug("\t"+strconv.Itoa(j)+")", racerString)
		}
	}
	log.Debug("--------------")

	// Do extra stuff
	leaderboardRecalculateSoloUnseeded()
}
