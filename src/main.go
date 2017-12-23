package main // In Go, executable commands must always use package main

import (
	"os"
	"path"

	"github.com/Zamiell/isaac-racing-server/src/log"
	"github.com/Zamiell/isaac-racing-server/src/models"
	raven "github.com/getsentry/raven-go"
	"github.com/joho/godotenv"
)

var (
	projectPath  = path.Join(os.Getenv("GOPATH"), "src", "github.com", "Zamiell", "isaac-racing-server")
	db           *models.Models
	races        = make(map[int]*Race)
	shutdownMode = 0
)

func main() {
	// Initialize logging (in log/log.go)
	log.Init()

	// Welcome message
	log.Info("+-------------------------------+")
	log.Info("| Starting isaac-racing-server. |")
	log.Info("+-------------------------------+")

	// Load the ".env" file which contains environment variables with secret values
	if err := godotenv.Load(path.Join(projectPath, ".env")); err != nil {
		log.Fatal("Failed to load .env file:", err)
	}

	// Configure error reporting to Sentry
	sentrySecret := os.Getenv("SENTRY_SECRET")
	if len(sentrySecret) == 0 {
		log.Info("The \"SENTRY_SECRET\" environment variable is blank; aborting Sentry initialization.")
	} else {
		raven.SetDSN("https://0d0a2118a3354f07ae98d485571e60be:" + sentrySecret + "@sentry.io/124813")
	}

	// Initialize the database model
	if v, err := models.Init(); err != nil {
		log.Fatal("Failed to open the database:", err)
	} else {
		db = v
	}
	defer db.Close()

	// Clean up any unfinished races from the database
	if nonStartedRaces, err := db.Races.Cleanup(); err != nil {
		log.Fatal("Failed to cleanup the leftover races:", err)
	} else {
		for _, raceID := range nonStartedRaces {
			log.Info("Deleted race", raceID, "during starting cleanup.")
		}
	}

	// Populate the achievements map (in achievements.go)
	achievementsInit()

	// Start the Twitch bot (in twitch.go)
	twitchInit()

	// Start the Discord bot (in discord.go)
	discordInit()

	// Initialize a WebSocket router using the Melody framework
	// (in websocket.go)
	websocketInit()

	loadAllItems()
	// Initialize an HTTP router using the Gin framework (in http.go)
	// (the "ListenAndServe" functions located inside here are blocking)

	httpInit()
}
