package main

/*
	Imports
*/

import (
	"html/template"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/Zamiell/isaac-racing-server/src/log"
	"github.com/Zamiell/isaac-racing-server/src/models"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/thirdparty/tollbooth_gin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

/*
	Constants
*/

const (
	sessionName = "isaac.sid"
)

/*
	Global variables
*/

var (
	sessionStore sessions.CookieStore
)

/*
	Data structures
*/

type TemplateData struct {
	Title string

	// Races stuff
	ResultsRaces   []models.RaceHistory
	TotalRaceCount int
	TotalPages     int
	PreviousPage   int
	NextPage       int

	// Profiles/profile stuff
	ResultsProfiles   []models.UserProfilesRow
	ResultsProfile    models.UserProfileData
	TotalProfileCount int
	UsersPerPage      int
}

/*
	Initialization function
*/

func httpInit() {
	// Create a new Gin HTTP router
	httpRouter := gin.Default()
	gin.SetMode(gin.ReleaseMode) // Uncomment this while debugging HTTP stuff

	// Create a session store
	sessionSecret := os.Getenv("SESSION_SECRET")
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
	if !useSSL {
		options.Secure = false
	}
	if useLocalhost {
		options.Domain = "localhost"
	}
	sessionStore.Options(options)
	httpRouter.Use(sessions.Sessions(sessionName, sessionStore))

	// Define middleware
	limiter := tollbooth.NewLimiter(1, time.Second) // Limit each user to 1 request per second
	httpRouter.Use(tollbooth_gin.LimitHandler(limiter))
	httpRouter.Use(httpMwGoogleAnalytics)

	// Path handlers (for the WebSocket server)
	httpRouter.POST("/login", httpLogin)
	httpRouter.POST("/register", httpRegister)
	httpRouter.GET("/ws", httpWS)

	// Path handlers (for the website)
	httpRouter.GET("/", httpHome)
	httpRouter.GET("/races", httpRaces)
	//httpRouter.GET("/profile", httpProfile)
	//httpRouter.GET("/profile/:player", httpProfile) // Handles profile username
	//httpRouter.GET("/profiles", httpProfiles)
	//httpRouter.GET("/profiles/:page", httpProfiles) // Handles extra pages for profiles
	//httpRouter.GET("/leaderboards", httpLeaderboards)
	httpRouter.GET("/info", httpInfo)
	httpRouter.GET("/download", httpDownload)
	httpRouter.Static("/public", "../public")

	// Figure out the port that we are using for the HTTP server
	var port int
	if useSSL {
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
	if useSSL {
		if err := http.ListenAndServeTLS(
			":"+strconv.Itoa(port), // Nothing before the colon implies 0.0.0.0
			sslCertFile,
			sslKeyFile,
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
		log.Error("Failed to create the template:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Execute the template and send it to the user
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Error("Failed to execute the template:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
