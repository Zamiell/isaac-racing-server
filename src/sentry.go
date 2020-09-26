package main

import (
	"os"

	sentry "github.com/getsentry/sentry-go"
)

func sentryInit() {
	// We only want to report errors in production
	if isDev {
		return false
	}

	// Read some configuration values from environment variables
	// (they were loaded from the ".env" file in "main.go")
	sentryDSN := os.Getenv("SENTRY_DSN")
	if len(sentrySecret) == 0 {
		logger.Info("The \"SENTRY_SECRET\" environment variable is blank; aborting Sentry initialization.")
		return false
	}

	// Initialize Sentry
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:          sentryDSN,
		IgnoreErrors: commonHTTPErrors,
		Release:      gitCommitOnStart,
		HTTPClient:   HTTPClientWithTimeout,
	}); err != nil {
		logger.Fatal("Failed to initialize Sentry:", err)
		return false
	}

	return true
}
