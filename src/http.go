package server

import (
	"html/template"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/Zamiell/isaac-racing-server/models"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

const (
	HTTPSessionName  = "isaac.sid"
	HTTPWriteTimeout = 10 * time.Second
)

var (
	sessionStore cookie.Store
	GATrackingID string

	// HTTPClientWithTimeout specifies a timeout on the default http.Client structure
	HTTPClientWithTimeout = &http.Client{
		Timeout: HTTPWriteTimeout,
	}
)

/*
	Data structures
*/

type TemplateData struct {
	Title string

	// Races stuff
	RaceResults       []models.RaceHistory
	SingleRaceResults models.RaceHistory
	SingleRaceFormat  string
	ResultsRaces      []models.RaceHistory
	TotalRaceCount    int
	TotalPages        int
	PreviousPage      int
	NextPage          int
	RaceResultsRanked []models.RaceHistory
	RaceResultsAll    []models.RaceHistory

	// Profiles/profile stuff
	ResultsProfiles   []models.ProfilesRow
	ResultsProfile    models.ProfileData
	TotalProfileCount int
	UsersPerPage      int
	TotalTime         int
	MissingPlayer     string

	// Leaderboard stuff
	LeaderboardSeeded       []models.LeaderboardRowSeeded
	LeaderboardSeededSolo   []models.LeaderboardRowSeededSolo
	LeaderboardUnseeded     []models.LeaderboardRowUnseeded
	LeaderboardUnseededSolo []models.LeaderboardRowUnseededSolo
	LeaderboardDiversity    []models.LeaderboardRowDiversity
	//LeaderboardSeeded []models.LeaderboardRowSeeded

	// Hall of Fame stuff
	Season1r9  []models.SpeedRun
	Season1r14 []models.SpeedRun
	Season2r7  []models.SpeedRun
	Season3r7  []models.SpeedRun
	Season4r7  []models.SpeedRun
	Season5r7  []models.SpeedRun
	Season6r7  []models.SpeedRun
	Season7r7  []models.SpeedRun
	// Season8r7 []models.SpeedRun

	// Tournament Stuff
	CurrentTournament bool
	TournamentRaces   []models.TournamentRace
	TournamentInfos   TournamentStats
	AllTournaments    []TournamentInfo
}

/*
	Initialization function
*/

func httpInit() {
	// Create a new Gin HTTP router
	gin.SetMode(gin.ReleaseMode) // Comment this out to debug HTTP stuff
	httpRouter := gin.Default()

	// Read some HTTP server configuration values from environment variables
	// (they were loaded from the .env file in main.go)
	sessionSecret := os.Getenv("SESSION_SECRET")
	if len(sessionSecret) == 0 {
		logger.Info("The \"SESSION_SECRET\" environment variable is blank; aborting HTTP initialization.")
		return
	}
	domain := os.Getenv("DOMAIN")
	if len(domain) == 0 {
		logger.Info("The \"DOMAIN\" environment variable is blank; aborting HTTP initialization.")
		return
	}
	tlsCertFile := os.Getenv("TLS_CERT_FILE")
	tlsKeyFile := os.Getenv("TLS_KEY_FILE")
	useTLS := true
	if len(tlsCertFile) == 0 || len(tlsKeyFile) == 0 {
		useTLS = false
	}

	// Create a session store
	sessionStore = cookie.NewStore([]byte(sessionSecret))
	options := sessions.Options{
		Path:   "/",
		Domain: domain,
		MaxAge: 5, // 5 seconds
		// After getting a cookie via "/login", the client will immediately
		// establish a WebSocket connection via "/ws", so the cookie only needs
		// to exist for that time frame
		Secure: true,
		// Only send the cookie over HTTPS:
		// https://www.owasp.org/index.php/Testing_for_cookies_attributes_(OTG-SESS-002)
		HttpOnly: true,
		// Mitigate XSS attacks:
		// https://www.owasp.org/index.php/HttpOnly
	}
	if !useTLS {
		options.Secure = false
	}
	sessionStore.Options(options)
	httpRouter.Use(sessions.Sessions(HTTPSessionName, sessionStore))

	// Use the Tollbooth Gin middleware for rate limiting
	limiter := tollbooth.NewLimiter(1, nil) // Limit each user to 1 request per second
	// When a user requests "/", they will also request the CSS and images;
	// this middleware is smart enough to know that it is considered part of the first request
	// However, it is still not possible to spam download CSS or image files
	limiterMiddleware := tollbooth_gin.LimitHandler(limiter)
	httpRouter.Use(limiterMiddleware)

	/*
		This was used as an alterate to the Tollbooth middleware when it wasn't working

		// Use the gin-limiter middleware for rate-limiting
		// We only allow 60 request per minute, an average of 1 per second
		// This is because when a user requests "/", they will also request the CSS and images
		// Based on: https://github.com/julianshen/gin-limiter/blob/master/example/web.go
		limiterMiddleware := limiter.NewRateLimiter(time.Second*60, 60, func(c *gin.Context) (string, error) {
			// Local variables
			r := c.Request
			ip, _, _ := net.SplitHostPort(r.RemoteAddr)

			// Just use the IP address as the key
			return ip, nil
		}).Middleware()
		httpRouter.Use(limiterMiddleware)
	*/

	// Use a custom middleware for Google Analytics tracking
	GATrackingID = os.Getenv("GA_TRACKING_ID")
	if len(GATrackingID) != 0 {
		httpRouter.Use(httpMwGoogleAnalytics)
	}

	// Path handlers (for the WebSocket server)
	httpRouter.POST("/login", httpLogin)
	httpRouter.POST("/register", httpRegister)
	httpRouter.GET("/ws", httpWS)

	// Path handlers (for the website)
	httpRouter.GET("/", httpHome)
	httpRouter.GET("/news", httpNews)
	httpRouter.GET("/races", httpRaces)
	httpRouter.GET("/races/:page", httpRaces)
	httpRouter.GET("/race", httpRace)
	httpRouter.GET("/race/:raceid", httpRace)
	httpRouter.GET("/profiles", httpProfiles)
	httpRouter.GET("/profiles/:page", httpProfiles) // Handles extra pages for profiles
	httpRouter.GET("/profile", httpProfile)
	httpRouter.GET("/profile/:player", httpProfile) // Handles profile username
	httpRouter.GET("/tournaments", httpTournament)
	httpRouter.GET("/leaderboards", httpLeaderboards)
	httpRouter.GET("/info", httpInfo)
	httpRouter.GET("/download", httpDownload)
	httpRouter.GET("/halloffame", httpHallOfFame)
	httpRouter.Static("/public", path.Join(projectPath, "public"))

	// Figure out the port that we are using for the HTTP server
	var port int
	if useTLS {
		// We want all HTTP requests to be redirected to HTTPS
		// (but make an exception for Let's Encrypt)
		// The Gin router is using the default serve mux,
		// so we need to create a new fresh one for the HTTP handler
		HTTPServeMux := http.NewServeMux()
		letsEncryptPath := path.Join(projectPath, "letsencrypt")
		HTTPServeMux.Handle("/.well-known/acme-challenge/", http.FileServer(http.FileSystem(http.Dir(letsEncryptPath))))
		HTTPServeMux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			http.Redirect(w, req, "https://"+req.Host+req.URL.String(), http.StatusMovedPermanently)
		}))

		// ListenAndServe is blocking, so start listening on a new goroutine
		go func() {
			// Nothing before the colon implies 0.0.0.0
			if err := http.ListenAndServe(":80", HTTPServeMux); err != nil {
				logger.Fatal("http.ListenAndServe failed:", err)
			}
			logger.Fatal("http.ListenAndServe ended for port 80.", nil)
		}()

		// 443 is the default port for HTTPS
		port = 443
	} else {
		// 80 is the default port for HTTP
		port = 80
	}

	// Start listening and serving requests (which is blocking)
	logger.Info("Listening on port " + strconv.Itoa(port) + ".")
	if useTLS {
		if err := http.ListenAndServeTLS(
			":"+strconv.Itoa(port), // Nothing before the colon implies 0.0.0.0
			tlsCertFile,
			tlsKeyFile,
			httpRouter,
		); err != nil {
			logger.Fatal("http.ListenAndServeTLS failed:", err)
		}
		logger.Fatal("http.ListenAndServeTLS ended prematurely.", nil)
	} else {
		// Listen and serve (HTTP)
		if err := http.ListenAndServe(
			":"+strconv.Itoa(port), // Nothing before the colon implies 0.0.0.0
			httpRouter,
		); err != nil {
			logger.Fatal("http.ListenAndServe failed:", err)
		}
		logger.Fatal("http.ListenAndServe ended prematurely.", nil)
	}
}

/*
	HTTP miscellaneous subroutines
*/

func httpServeTemplate(w http.ResponseWriter, templateName string, data interface{}) {
	lp := path.Join(projectPath, "src", "views", "layout.tmpl")
	fp := path.Join(projectPath, "src", "views", templateName+".tmpl")

	// Return a 404 if the template doesn't exist
	info, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// Create the template
	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		logger.Error("Failed to create the template: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Execute the template and send it to the user
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		if strings.HasSuffix(err.Error(), ": write: broken pipe") ||
			strings.HasSuffix(err.Error(), ": client disconnected") ||
			strings.HasSuffix(err.Error(), ": http2: stream closed") ||
			strings.HasSuffix(err.Error(), ": write: connection timed out") {

			// Broken pipe errors can occur when the user presses the "Stop" button while the template is executing
			// We don't want to reporting these errors to Sentry
			// https://stackoverflow.com/questions/26853200/filter-out-broken-pipe-errors-from-template-execution
			// I don't know exactly what the other errors mean
			logger.Info("Ordinary error when executing the template: " + err.Error())
		} else {
			logger.Error("Failed to execute the template: " + err.Error())
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
