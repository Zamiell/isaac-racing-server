package main // In Go, executable commands must always use package main

/*
 *  Imports
 */

import (
	"github.com/Zamiell/isaac-racing-server/models"

	"net/http" // For establishing an HTTP server
	"os"       // For logging and reading environment variables
	"sync"     // For locking and unlocking the connection map
	"time"     // For dealing with timestamps

	"github.com/didip/tollbooth"  // For rate-limiting login requests
	"github.com/gorilla/context"  // For cookie sessions (1/2)
	"github.com/gorilla/sessions" // For cookie sessions (2/2)
	"github.com/joho/godotenv"    // For reading environment variables that contain secrets
	"github.com/op/go-logging"    // For logging
	"github.com/trevex/golem"     // The Golem WebSocket framework
)

/*
 *  Constants
 */

const (
	port        = "443"
	sessionName = "isaac.sid"
	domain      = "isaacitemtracker.com"
	auth0Domain = "isaacserver.auth0.com"
	useSSL      = true
	sslCertFile = "/etc/letsencrypt/live/isaacitemtracker.com/fullchain.pem"
	sslKeyFile  = "/etc/letsencrypt/live/isaacitemtracker.com/privkey.pem"
)

/*
 *  Global variables
 */

var (
	log           = logging.MustGetLogger("isaac")
	sessionStore  *sessions.CookieStore
	roomManager   = golem.NewRoomManager()
	pmManager     = golem.NewRoomManager()
	connectionMap = struct {
		// Maps are not safe for concurrent use: https://blog.golang.org/go-maps-in-action
		sync.RWMutex
		m map[string]*ExtendedConnection
	}{m: make(map[string]*ExtendedConnection)}
	chatRoomMap = struct {
		// Maps are not safe for concurrent use: https://blog.golang.org/go-maps-in-action
		sync.RWMutex
		m map[string][]User
	}{m: make(map[string][]User)}
	db *model.Model
)

/*
 *  Program entry point
 */

func main() {
	// Configure logging: http://godoc.org/github.com/op/go-logging#Formatter
	loggingBackend := logging.NewLogBackend(os.Stdout, "", 0)
	logFormat := logging.MustStringFormatter( // https://golang.org/pkg/time/#Time.Format
		`%{time:Mon Jan 2 15:04:05 MST 2006} - %{level:.4s} - %{shortfile} - %{message}`,
	)
	loggingBackendFormatted := logging.NewBackendFormatter(loggingBackend, logFormat)
	logging.SetBackend(loggingBackendFormatted)

	// Load the .env file which contains environment variables with secret values
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load .env file:", err)
	}

	// Create a session store
	sessionSecret := os.Getenv("SESSION_SECRET")
	sessionStore = sessions.NewCookieStore([]byte(sessionSecret))
	sessionStore.Options = &sessions.Options{
		Domain:   domain,
		Path:     "/",
		MaxAge:   5,    // 5 seconds
		Secure:   true, // Only send the cookie over HTTPS: https://www.owasp.org/index.php/Testing_for_cookies_attributes_(OTG-SESS-002)
		HttpOnly: true, // Mitigate XSS attacks: https://www.owasp.org/index.php/HttpOnly
	}

	// Initialize the database model
	db = model.GetModel(log)

	// Clean up any non-started races before we start
	db.Races.Cleanup()

	// Initialize the achievements
	achievementsInit()

	// Create a WebSocket router using the Golem framework
	router := golem.NewRouter()
	router.SetConnectionExtension(NewExtendedConnection)
	router.OnHandshake(validateSession)
	router.OnConnect(connOpen)
	router.OnClose(connClose)

	/*
	 *  The websocket commands
	 */

	// Chat commands
	router.On("roomJoin", roomJoin)
	router.On("roomLeave", roomLeave)
	router.On("roomMessage", roomMessage)
	router.On("privateMessage", privateMessage)
	router.On("roomGetAll", roomGetAll)

	// Race commands
	router.On("raceCreate", raceCreate)
	router.On("raceJoin", raceJoin)
	router.On("raceLeave", raceLeave)
	router.On("raceReady", raceReady)
	router.On("raceUnready", raceUnready)
	router.On("raceRuleset", raceRuleset)
	router.On("raceDone", raceDone)
	router.On("raceQuit", raceQuit)
	router.On("raceComment", raceComment)
	router.On("raceItem", raceItem)
	router.On("raceFloor", raceFloor)

	// Profile commands
	router.On("profileGet", profileGet)
	router.On("profileSetUsername", profileSetUsername)

	// Admin commands
	router.On("adminBan", adminBan)
	router.On("adminUnban", adminUnban)
	router.On("adminBanIP", adminBanIP)
	router.On("adminUnbanIP", adminUnbanIP)
	router.On("adminSquelch", adminSquelch)
	router.On("adminUnsquelch", adminUnsquelch)
	router.On("adminPromote", adminPromote)
	router.On("adminDemote", adminDemote)

	// Miscellaneous
	router.On("logout", logout)

	/*
	 *  HTTP stuff
	 */

	// Assign functions to URIs
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))            // Serve static files
	http.HandleFunc("/", serveTemplate)                                                                   // Anything that is not a static file will match this
	http.Handle("/login", tollbooth.LimitFuncHandler(tollbooth.NewLimiter(1, time.Second), loginHandler)) // Rate limit the login handler
	http.HandleFunc("/ws", router.Handler())                                                              // The golem router handles websockets

	// Welcome message
	log.Info("Starting isaac-racing-server on port " + port + ".")

	// Listen and serve
	if useSSL == true {
		if err := http.ListenAndServeTLS(
			":"+port, // Nothing before the colon implies 0.0.0.0
			sslCertFile,
			sslKeyFile,
			context.ClearHandler(http.DefaultServeMux), // We wrap with context.ClearHandler or else we will leak memory: http://www.gorillatoolkit.org/pkg/sessions
		); err != nil {
			log.Fatal("ListenAndServeTLS failed:", err)
		}
	} else {
		// Listen and serve (HTTP)
		if err := http.ListenAndServe(
			":"+port, // Nothing before the colon implies 0.0.0.0
			context.ClearHandler(http.DefaultServeMux), // We wrap with context.ClearHandler or else we will leak memory: http://www.gorillatoolkit.org/pkg/sessions
		); err != nil {
			log.Fatal("ListenAndServeTLS failed:", err)
		}
	}
}
