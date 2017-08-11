package main // In Go, executable commands must always use package main

/*
	Imports
*/

import (
	"net/http"
	"os"
	"path"
	"time"

	"github.com/Zamiell/isaac-racing-server/src/log"
	"github.com/Zamiell/isaac-racing-server/src/models"
	raven "github.com/getsentry/raven-go"
	"github.com/joho/godotenv"
)

/*
	Constants
*/

const (
	domain       = "isaacracing.net"
	useSSL       = false
	sslCertFile  = "/etc/letsencrypt/live/" + domain + "/fullchain.pem"
	sslKeyFile   = "/etc/letsencrypt/live/" + domain + "/privkey.pem"
	useTwitch    = false           // If true, it will run a Twitch chat bot in a new Goroutine
	useDiscord   = false           // If true, it will run a Discord chat bot in a new Goroutine
	useLocalhost = true            // If true, it will use "localhost" instead of the above domain
	GATrackingID = "UA-91999156-1" // For Google Analytics
)

/*
	Global variables
*/

var (
	projectPath = path.Join(os.Getenv("GOPATH"), "src", "github.com", "Zamiell", "isaac-racing-server")
	db          *models.Models

	// We don't want to use the default http.Client structure because it has no default timeout set
	myHTTPClient = &http.Client{
		Timeout: 10 * time.Second,
	}
)

/*
	No directory listing stuff from: https://marc.ttias.be/golang-nuts/2016-03/msg00888.php
*/

/*
type justFilesFilesystem struct {
	fs http.FileSystem
}

func (fs justFilesFilesystem) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return neuteredReaddirFile{f}, nil
}

type neuteredReaddirFile struct {
	http.File
}

func (f neuteredReaddirFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}
*/

/*
	Program entry point
*/

func main() {
	// Initialize logging
	log.Init()

	// Welcome message
	log.Info("+-------------------------------+")
	log.Info("| Starting isaac-racing-server. |")
	log.Info("+-------------------------------+")

	// Load the .env file which contains environment variables with secret values
	err := godotenv.Load(path.Join(projectPath, ".env"))
	if err != nil {
		log.Fatal("Failed to load .env file:", err)
	}

	// Configure error reporting to Sentry
	sentrySecret := os.Getenv("SENTRY_SECRET")
	raven.SetDSN("https://0d0a2118a3354f07ae98d485571e60be:" + sentrySecret + "@sentry.io/124813")

	// Initialize the database model
	if db, err = models.Init(path.Join(projectPath, "database.sqlite")); err != nil {
		log.Fatal("Failed to open the database:", err)
	}

	// Clean up any non-started races before we start
	if nonStartedRaces, err := db.Races.Cleanup(); err != nil {
		log.Fatal("Failed to cleanup the leftover races:", err)
	} else {
		for _, raceID := range nonStartedRaces {
			log.Info("Deleted race", raceID, "during starting cleanup.")
		}
	}

	// Initiate the "end of race" function for each of the non-finished races in 30 minutes
	// (this is normally initiated on race start)
	/*
		if startedRaces, err := db.Races.GetCurrentRaces(); err != nil {
			log.Fatal("Failed to start  the leftover races:", err)
		} else {
			// TODO
				for _, race := range startedRaces {
					go raceCheckStart3(race.ID)
				}
		}
	*/

	// Add the achievements to the database (in achievements.go)
	//achievementsInit()

	// Start the Twitch bot (in twitch.go)
	twitchInit()

	// Start the Discord bot (in discord.go)
	discordInit()

	// Initialize a WebSocket router using the Melody framework (in websocket.go)
	websocketInit()

	// Initialize an HTTP router using the Gin framework (in http.go)
	httpInit()
}
