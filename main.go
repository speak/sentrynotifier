package main

import (
	"errors"
	"flag"
	"log"
	"os"

	"github.com/getsentry/raven-go"
)

var (
	hostname        string
	sentry_dsn      string
	app_environment string
	report          string
)

func main() {
	// Set up the flags for the command
	flag.StringVar(&hostname, "hostname", os.Getenv("HOSTNAME"), "The server name to report to Sentry")
	flag.StringVar(&sentry_dsn, "sentry-dsn", os.Getenv("SENTRY_DSN"), "The Sentry DSN to use for the error report")
	flag.StringVar(&app_environment, "app_environment", os.Getenv("APP_ENVIRONMENT"), "The environment to report for, e.g. staging, production")
	flag.StringVar(&report, "report", os.Getenv("REPORT"), "The error to report")

	// Parse the command line flags
	flag.Parse()

	// Make sure we've got everything
	if hostname == "" {
		log.Fatal("Hostname is not configured")
	}
	if sentry_dsn == "" {
		log.Fatal("Sentry DSN not configured")
	}
	if app_environment == "" {
		log.Fatal("You must set the app environment")
	}
	if report == "" {
		log.Fatal("You must pass an error to report")
	}

	// Create the Raven Client
	ravenClient, err := raven.NewClient(sentry_dsn, map[string]string{"environment": app_environment, "server_name": hostname})
	if err != nil {
		log.Fatal(err)
	}

	// Set up the packet to send
	packet := raven.NewPacket(report, raven.NewException(errors.New(report), raven.NewStacktrace(0, 5, nil)))

	// Send the packet
	_, errch := ravenClient.Capture(packet, nil)

	// Handle Sentry response
	if err := <-errch; err != nil {
		log.Fatalf("Error sending to Raven: %v", err)
	} else {
		log.Printf("Sent error to Sentry: %v", err)
	}

	// Quit
	os.Exit(0)
}
