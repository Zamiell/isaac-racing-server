package main

import (
	"html/template"
	"net"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/Zamiell/isaac-racing-server/src/log"
	"github.com/Zamiell/isaac-racing-server/src/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	limiter "github.com/julianshen/gin-limiter"
)

const (
	sessionName = "isaac.sid"
)

var (
	sessionStore sessions.CookieStore
	GATrackingID string
	myHTTPClient = &http.Client{ // We don't want to use the default http.Client structure because it has no default timeout set
		Timeout: 10 * time.Second,
	}
)

/*
	Data structures
*/

type TemplateData struct {
	Title string

	// Races stuff
	RaceResults    []models.RaceHistory
	ResultsRaces   []models.RaceHistory
	TotalRaceCount int
	TotalPages     int
	PreviousPage   int
	NextPage       int

	// Profiles/profile stuff
	ResultsProfiles   []models.ProfilesRow
	ResultsProfile    models.ProfileData
	TotalProfileCount int
	UsersPerPage      int
}

/*
	Initialization function
*/

func httpInit() {
	// Create a new Gin HTTP router
	//gin.SetMode(gin.ReleaseMode) // Comment this out to debug HTTP stuff
	httpRouter := gin.Default()

	// Read some HTTP server configuration values from environment variables
	// (they were loaded from the .env file in main.go)
	sessionSecret := os.Getenv("SESSION_SECRET")
	if len(sessionSecret) == 0 {
		log.Info("The \"SESSION_SECRET\" environment variable is blank; aborting HTTP initalization.")
		return
	}
	domain := os.Getenv("DOMAIN")
	if len(domain) == 0 {
		log.Info("The \"DOMAIN\" environment variable is blank; aborting HTTP initalization.")
		return
	}
	tlsCertFile := os.Getenv("TLS_CERT_FILE")
	tlsKeyFile := os.Getenv("TLS_KEY_FILE")
	useTLS := true
	if len(tlsCertFile) == 0 || len(tlsKeyFile) == 0 {
		useTLS = false
	}

	// Create a session store
	sessionStore = sessions.NewCookieStore([]byte(sessionSecret))
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
	httpRouter.Use(sessions.Sessions(sessionName, sessionStore))

	/*
		Commented out because it doesn't work:
		https://github.com/didip/tollbooth_gin/issues/3

		// Use the Tollbooth Gin middleware for Google Analytics tracking
		limiter := tollbooth.NewLimiter(1, time.Second, nil) // Limit each user to 1 request per second
		httpRouter.Use(tollbooth_gin.LimitHandler(limiter))
	*/

	// Use the gin-limiter middleware for rate-limiting
	// (to only allow one request per second)
	// Based on: https://github.com/julianshen/gin-limiter/blob/master/example/web.go
	limiterMiddleware := limiter.NewRateLimiter(time.Second*60, 60, func(c *gin.Context) (string, error) {
		// Local variables
		r := c.Request
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		// Just use the IP address as the key
		return ip, nil
	}).Middleware()
	httpRouter.Use(limiterMiddleware)

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

	// Path handlers for single profile
	httpRouter.GET("/profile", httpProfile)
	httpRouter.GET("/profile/:player", httpProfile) // Handles profile username

	// Path handlers for all profiles
	httpRouter.GET("/profiles", httpProfiles)
	httpRouter.GET("/profiles/:page", httpProfiles) // Handles extra pages for profiles

	// Path handlers for race page
	httpRouter.GET("/race", httpRace)
	httpRouter.GET("/race/:raceid", httpRace)

	// Path handlers for races page
	httpRouter.GET("/races", httpRaces)
	httpRouter.GET("/races/:page", httpRaces)

	//	httpRouter.GET("/leaderboards", httpLeaderboards)
	httpRouter.GET("/info", httpInfo)
	httpRouter.GET("/download", httpDownload)
	httpRouter.Static("/public", "../public")

	// Figure out the port that we are using for the HTTP server
	var port int
	if useTLS {
		// We want all HTTP requests to be redirected to HTTPS
		// (but make an exception for Let's Encrypt)
		// The Gin router is using the default serve mux, so we need to create a
		// new fresh one for the HTTP handler
		HTTPServeMux := http.NewServeMux()
		HTTPServeMux.Handle("/.well-known/acme-challenge/", http.FileServer(http.FileSystem(http.Dir("letsencrypt"))))
		HTTPServeMux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			http.Redirect(w, req, "https://"+req.Host+req.URL.String(), http.StatusMovedPermanently)
		}))

		// ListenAndServe is blocking, so start listening on a new goroutine
		go func() {
			http.ListenAndServe(":80", HTTPServeMux) // Nothing before the colon implies 0.0.0.0
			log.Fatal("http.ListenAndServe ended for port 80.", nil)
		}()

		// 443 is the default port for HTTPS
		port = 443
	} else {
		// 80 is the defeault port for HTTP
		port = 80
	}

	// Start listening and serving requests (which is blocking)
	log.Info("Listening on port " + strconv.Itoa(port) + ".")
	if useTLS {
		if err := http.ListenAndServeTLS(
			":"+strconv.Itoa(port), // Nothing before the colon implies 0.0.0.0
			tlsCertFile,
			tlsKeyFile,
			httpRouter,
		); err != nil {
			log.Fatal("http.ListenAndServeTLS failed:", err)
		}
		log.Fatal("http.ListenAndServeTLS ended prematurely.", nil)
	} else {
		// Listen and serve (HTTP)
		if err := http.ListenAndServe(
			":"+strconv.Itoa(port), // Nothing before the colon implies 0.0.0.0
			httpRouter,
		); err != nil {
			log.Fatal("http.ListenAndServe failed:", err)
		}
		log.Fatal("http.ListenAndServe ended prematurely.", nil)
	}
}

/*
	HTTP miscellaneous subroutines
*/

func httpServeTemplate(w http.ResponseWriter, templateName string, data interface{}) {
	lp := path.Join("views", "layout.tmpl")
	fp := path.Join("views", templateName+".tmpl")

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
		log.Error("Failed to create the template: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Execute the template and send it to the user
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		if strings.HasSuffix(err.Error(), ": write: broken pipe") {
			// Broken pipe errors can occur when the user presses the "Stop" button while the template is executing
			// We don't want to reporting these errors to Sentry
			// https://stackoverflow.com/questions/26853200/filter-out-broken-pipe-errors-from-template-execution
			log.Info("Failed to execute the template: " + err.Error())
		} else {
			log.Error("Failed to execute the template: " + err.Error())
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
