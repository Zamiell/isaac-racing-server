package server

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/Zamiell/isaac-racing-server/models"
	"github.com/joho/godotenv"
)

var (
	projectPath string

	logger           *Logger
	gitCommitOnStart string
	isDev            bool
	usingSentry      bool
	db               *models.Models
	races            = make(map[int]*Race)
	shutdownMode     = 0
)

func Init() {
	// Initialize logging (in "logger.go")
	logger = NewLogger()

	// Welcome message
	logger.Info("+-------------------------------+")
	logger.Info("| Starting isaac-racing-server. |")
	logger.Info("+-------------------------------+")

	// Get the project path
	// https://stackoverflow.com/questions/18537257/
	if v, err := os.Executable(); err != nil {
		logger.Fatal("Failed to get the path of the currently running executable:", err)
	} else {
		projectPath = filepath.Dir(v)
	}

	// Record the commit that corresponds with when the Golang code was compiled
	cmd := exec.Command("git", "rev-parse", "HEAD")
	if stdout, err := cmd.Output(); err != nil {
		logger.Fatal("Failed to perform a \"git rev-parse HEAD\":", err)
		return
	} else {
		gitCommitOnStart = strings.TrimSpace(string(stdout))
	}

	// Check to see if the ".env" file exists
	envPath := path.Join(projectPath, ".env")
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		logger.Fatal("The \"" + envPath + "\" file does not exist.")
		return
	} else if err != nil {
		logger.Fatal("Failed to check if the \""+envPath+"\" file exists:", err)
		return
	}

	// Load the ".env" file which contains environment variables with secret values
	if err := godotenv.Load(path.Join(projectPath, ".env")); err != nil {
		logger.Fatal("Failed to load .env file:", err)
	}

	if os.Getenv("DOMAIN") == "" ||
		os.Getenv("DOMAIN") == "localhost" ||
		strings.HasPrefix(os.Getenv("DOMAIN"), "192.168.") ||
		strings.HasPrefix(os.Getenv("DOMAIN"), "10.") {

		isDev = true
	}

	// Initialize Sentry (in "sentry.go")
	usingSentry = sentryInit()

	// Initialize the database model
	if v, err := models.Init(); err != nil {
		logger.Fatal("Failed to open the database:", err)
	} else {
		db = v
	}
	defer db.Close()

	// Clean up any unfinished races from the database
	if nonStartedRaces, err := db.Races.Cleanup(); err != nil {
		logger.Fatal("Failed to cleanup the leftover races:", err)
	} else {
		for _, raceID := range nonStartedRaces {
			logger.Info("Deleted race", raceID, "during starting cleanup.")
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

	// Initialize the needed static maps for items (in constants.go)
	loadAllItems()
	loadAllBuilds()

	// Initialize the needed static maps for tournaments (in constants.go)
	loadAllTournaments()

	shadowInit()

	// Initialize an HTTP router using the Gin framework (in http.go)
	// (the "ListenAndServe" functions located inside here are blocking)
	httpInit()
}
