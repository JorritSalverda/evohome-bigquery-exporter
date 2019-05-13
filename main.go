package main

import (
	"log"
	"os"
	"runtime"

	"github.com/alecthomas/kingpin"
)

var (
	// set when building the application
	app       string
	version   string
	branch    string
	revision  string
	buildDate string
	goVersion = runtime.Version()

	// application specific config
	username = kingpin.Flag("username", "Evohome username.").Envar("EVOHOME_USERNAME").Required().String()
	password = kingpin.Flag("password", "Evohome password.").Envar("EVOHOME_PASSWORD").Required().String()
)

func main() {

	// parse command line parameters
	kingpin.Parse()

	// log to stdout and hide timestamp
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	// log startup message
	log.Printf("Starting %v version %v...", app, version)

	client, err := NewEvohomeClient()
	if err != nil {
		log.Fatal("Failed creating evohome client: ", err)
	}

	sessionID, userID, err := client.GetSession(*username, *password)
	if err != nil {
		log.Fatalf("Failed retrieving session id for username %v: %v", *username, err)
	}

	locations, err := client.GetLocations(sessionID, userID)
	if err != nil {
		log.Fatalf("Failed retrieving locations for userid %v: %v", userID, err)
	}

	log.Printf("Retrieved %v locations: ", len(locations))

	// done
	log.Printf("Finished exporting metrics")
}
