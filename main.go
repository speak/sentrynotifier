package main

import (
	"errors"
	"flag"
	"log"
	"os"

	"github.com/getsentry/raven-go"
)

func main() {
	var hostname = flag.String("hostname", os.Getenv("HOSTNAME"), "The server name to report to Sentry")
	var sentry_dsn = flag.String("sentry-dsn", os.Getenv("SENTRY_DSN"), "The Sentry DSN to use for the error report")
	var app_environment = flag.String("app_environment", os.Getenv("APP_ENVIRONMENT"), "The environment to report for, e.g. staging, production")
	var report = flag.String("err", os.Getenv("APP_ERROR"), "The error to report")

	if *hostname == "" {
		log.Fatal("Hostname is not configured")
	}
	if *sentry_dsn == "" {
		log.Fatal("Sentry DSN not configured")
	}
	if *app_environment == "" {
		log.Fatal("You must set the app environment")
	}
	if *report == "" {
		log.Fatal("You must pass an error to report")
	}

	ravenClient, err := raven.NewClient(*sentry_dsn, map[string]string{"environment": *app_environment, "server_name": *hostname})
	if err != nil {
		log.Fatal(err)
	}
	packet := raven.NewPacket(*report, raven.NewException(errors.New(*report), raven.NewStacktrace(0, 5, nil)))
	_, errch := ravenClient.Capture(packet, nil)
	if err := <-errch; err != nil {
		log.Fatalf("Error sending to Raven: %v", err)
	} else {
		log.Printf("Sent error to Sentry: %v", err)
	}
	os.Exit(0)
}
