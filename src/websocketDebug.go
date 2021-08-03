package server

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	melody "gopkg.in/olahol/melody.v1"
)

func websocketDebug(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	username := d.v.Username
	admin := d.v.Admin

	// Validate that the user is an admin
	if admin == 0 {
		logger.Info("User \"" + username + "\" tried to do a debug command, but they are not staff/admin.")
		websocketError(s, d.Command, "Only staff members or administrators can do that.")
		return
	}

	logger.Debug("---------------------------------------------------------------")

	// Print out all of the current races
	if len(races) == 0 {
		logger.Debug("[no current races]")
	}
	for i, race := range races { // This is a map[int]*Race
		logger.Debug(strconv.Itoa(i) + " - " + race.Name)
		logger.Debug("\n")

		// Print out all of the fields
		// From: https://stackoverflow.com/questions/24512112/how-to-print-struct-variables-in-console
		logger.Debug("    All fields:")
		fieldsToIgnore := []string{
			"Racers",
			"Ruleset",
		}
		s := reflect.ValueOf(race).Elem()
		maxChars := 0
		for i := 0; i < s.NumField(); i++ {
			fieldName := s.Type().Field(i).Name
			if stringInSlice(fieldName, fieldsToIgnore) {
				continue
			}
			if len(fieldName) > maxChars {
				maxChars = len(fieldName)
			}
		}
		for i := 0; i < s.NumField(); i++ {
			fieldName := s.Type().Field(i).Name
			if stringInSlice(fieldName, fieldsToIgnore) {
				continue
			}
			f := s.Field(i)
			line := "  "
			for i := len(fieldName); i < maxChars; i++ {
				line += " "
			}
			line += "%s = %v"
			line = fmt.Sprintf(line, fieldName, f.Interface())
			if strings.HasSuffix(line, " = ") {
				line += "[empty string]"
			}
			line += "\n"
			logger.Debug(line)
		}
		logger.Debug("\n")

		// Manually enumerate the slices and maps
		logger.Debug("    Racers:")
		for name, racer := range race.Racers {
			logger.Debug("        " + name)
			s3 := reflect.ValueOf(racer).Elem()
			maxChars3 := 0
			for i := 0; i < s3.NumField(); i++ {
				fieldName := s3.Type().Field(i).Name
				if len(fieldName) > maxChars3 {
					maxChars3 = len(fieldName)
				}
			}
			for i := 0; i < s3.NumField(); i++ {
				fieldName := s3.Type().Field(i).Name
				f := s3.Field(i)
				line := "    "
				for i := len(fieldName); i < maxChars3; i++ {
					line += " "
				}
				line += "%s = %v"
				line = fmt.Sprintf(line, fieldName, f.Interface())
				if strings.HasSuffix(line, " = ") {
					line += "[empty string]"
				}
				line += "\n"
				logger.Debug(line)
			}
			logger.Debug("\n")
		}

		logger.Debug("---------------------------------------------------------------")
	}

	// Print out all of the current users
	logger.Debug("Current users:")
	if len(websocketSessions) == 0 {
		logger.Debug("    [no users]")
	}
	i := 1
	for name := range websocketSessions { // This is a map[string]*melody.Session
		logger.Debug("    " +
			strconv.Itoa(i) + ") " + name)
	}
	logger.Debug("---------------------------------------------------------------")

	// Do extra stuff
	/*
		logger.Debug("Calculating unseeded solo stats.")
		leaderboardRecalculateSoloUnseeded()
		logger.Debug("Finished.")
	*/
}
