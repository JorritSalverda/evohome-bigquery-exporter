package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	stdlog "log"
	"os"
	"runtime"
	"time"

	"github.com/alecthomas/kingpin"
	"github.com/ericchiang/k8s"
	corev1 "github.com/ericchiang/k8s/apis/core/v1"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
	username              = kingpin.Flag("username", "Evohome username.").Envar("EVOHOME_USERNAME").Required().String()
	password              = kingpin.Flag("password", "Evohome password.").Envar("EVOHOME_PASSWORD").Required().String()
	sessionSecretPath     = kingpin.Flag("session-secret-path", "Path to session secret.").Default("/secrets/session.json").OverrideDefaultFromEnvar("SESSION_SECRET_PATH").String()
	sessionSecretName     = kingpin.Flag("session-secret-name", "Name of the session secret.").Default("evohome-bigquery-exporter-credentials").OverrideDefaultFromEnvar("SESSION_SECRET_NAME").String()
	sessionTimeoutMinutes = kingpin.Flag("session-timeout-minutes", "Number of minutes before a session has to be refreshed.").Default("30").OverrideDefaultFromEnvar("SESSION_TIMEOUT_MINUTES").Int()
	stateFilePath         = kingpin.Flag("state-file-path", "Path to file with state from evohome-hgi80-listener.").Default("/state/state.json").OverrideDefaultFromEnvar("STATE_FILE_PATH").String()
	namespace             = kingpin.Flag("namespace", "Namespace the pod runs in.").Envar("NAMESPACE").Required().String()
	bigqueryProjectID     = kingpin.Flag("bigquery-project-id", "Google Cloud project id that contains the BigQuery dataset").Envar("BQ_PROJECT_ID").Required().String()
	bigqueryDataset       = kingpin.Flag("bigquery-dataset", "Name of the BigQuery dataset").Envar("BQ_DATASET").Required().String()
	bigqueryTable         = kingpin.Flag("bigquery-table", "Name of the BigQuery table").Envar("BQ_TABLE").Required().String()
	outdoorZoneName       = kingpin.Flag("outdoor-zone-name", "Name of the zone representing the outdoor temperature and humidity").Default("Outside").OverrideDefaultFromEnvar("OUTDOOR_ZONE_NAME").String()
)

func main() {

	// parse command line parameters
	kingpin.Parse()

	// log as severity for stackdriver logging to recognize the level
	zerolog.LevelFieldName = "severity"

	// set some default fields added to all logs
	log.Logger = zerolog.New(os.Stdout).With().
		Timestamp().
		Str("app", app).
		Str("version", version).
		Logger()

	// use zerolog for any logs sent via standard log library
	stdlog.SetFlags(0)
	stdlog.SetOutput(log.Logger)

	// log startup message
	log.Info().
		Str("branch", branch).
		Str("revision", revision).
		Str("buildDate", buildDate).
		Str("goVersion", goVersion).
		Msgf("Starting %v version %v...", app, version)

	evoClient, err := NewEvohomeClient()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed creating evohome client")
	}

	zoneInfoMap := readStateFromStateFile()

	bigqueryClient, err := NewBigQueryClient(*bigqueryProjectID)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed creating bigquery client")
	}
	initBigqueryTable(bigqueryClient)

	validSessionSecret, sessionSecret := readSessionSecretFromFile()

	if !validSessionSecret {
		sessionSecret = refreshSessionSecret()
	}

	log.Info().Msgf("Retrieving locations for user with id %v...", sessionSecret.UserID)

	locations, err := evoClient.GetLocations(sessionSecret.SessionID, sessionSecret.UserID)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed retrieving locations for userid %v", sessionSecret.UserID)
	}

	log.Debug().Interface("locations", locations).Msgf("Retrieved %v locations: ", len(locations))

	log.Debug().Msg("Mapping locations to measurements")
	measurements := mapLocationsToMeasurements(locations, *outdoorZoneName, zoneInfoMap)

	log.Debug().Msgf("Inserting measurements into table %v.%v.%v...", *bigqueryProjectID, *bigqueryDataset, *bigqueryTable)
	err = bigqueryClient.InsertMeasurements(*bigqueryDataset, *bigqueryTable, measurements)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed inserting measurements into bigquery table")
	}

	// done
	log.Info().Msg("Finished exporting metrics")
}

func readSessionSecretFromFile() (validSessionSecret bool, sessionSecret SessionSecret) {
	// check if session key exists in secret
	validSessionSecret = false
	if _, err := os.Stat(*sessionSecretPath); !os.IsNotExist(err) {

		log.Info().Msgf("File %v exists, reading contents...", *sessionSecretPath)

		// read secret
		data, err := ioutil.ReadFile(*sessionSecretPath)
		if err != nil {
			log.Fatal().Err(err).Msgf("Failed reading file from path %v", *sessionSecretPath)
		}

		log.Info().Msgf("Unmarshalling file %v contents...", *sessionSecretPath)

		// unmarshal secret
		if err := json.Unmarshal(data, &sessionSecret); err != nil {
			log.Fatal().Err(err).Interface("data", data).Msg("Failed unmarshalling session secret")
		}

		log.Info().Interface("RetrievedAt", sessionSecret.RetrievedAt).Msgf("Unmarshalled session secret, checking age...")

		// check if session secret isn't too old
		if sessionSecret.RetrievedAt.Add(time.Minute * time.Duration(*sessionTimeoutMinutes)).After(time.Now().UTC()) {
			validSessionSecret = true
			log.Info().Msg("Session secret is still valid...")
		}
	}
	return
}

func refreshSessionSecret() SessionSecret {
	log.Info().Msg("No valid session secret, retrieving new session id...")

	sessionID, userID, err := evoClient.GetSession(*username, *password)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed retrieving session id for username %v", *username)
	}

	sessionSecret := SessionSecret{
		SessionID:   sessionID,
		UserID:      userID,
		RetrievedAt: time.Now().UTC(),
	}

	// create kubernetes api client
	kubeClient, err := k8s.NewInClusterClient()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed creating Kubernetes API client")
	}

	log.Info().Msg("Retrieved new session id, storing it in secret for using it in the next scheduled pod...")

	// retrieve secret
	var secret corev1.Secret
	err = kubeClient.Get(context.Background(), *namespace, *sessionSecretName, &secret)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed retrieving secret %v", *sessionSecretName)
	}

	// marshal session secret to json
	sessionSecretData, err := json.Marshal(sessionSecret)

	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}

	secret.Data["session.json"] = sessionSecretData

	// update secret to have session information available when the application runs the next time
	err = kubeClient.Update(context.Background(), &secret)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed updating secret %v", *sessionSecretName)
	}

	log.Info().Msgf("Stored session secret in secret %v...", *sessionSecretName)

	return sessionSecret
}

func initBigqueryTable(bigqueryClient BigQueryClient) {

	log.Debug().Msgf("Checking if table %v.%v.%v exists...", *bigqueryProjectID, *bigqueryDataset, *bigqueryTable)
	tableExist := bigqueryClient.CheckIfTableExists(*bigqueryDataset, *bigqueryTable)
	if !tableExist {
		log.Debug().Msgf("Creating table %v.%v.%v...", *bigqueryProjectID, *bigqueryDataset, *bigqueryTable)
		err := bigqueryClient.CreateTable(*bigqueryDataset, *bigqueryTable, BigQueryMeasurement{}, "measured_at", true)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed creating bigquery table")
		}
	} else {
		log.Debug().Msgf("Trying to update table %v.%v.%v schema...", *bigqueryProjectID, *bigqueryDataset, *bigqueryTable)
		err := bigqueryClient.UpdateTableSchema(*bigqueryDataset, *bigqueryTable, BigQueryMeasurement{})
		if err != nil {
			log.Fatal().Err(err).Msg("Failed updating bigquery table schema")
		}
	}
}

func readStateFromStateFile() map[int64]ZoneInfo {

	zoneInfoMap := map[int64]ZoneInfo{}

	// check if state file exists in configmap
	var state State
	if _, err := os.Stat(*stateFilePath); !os.IsNotExist(err) {

		log.Info().Msgf("File %v exists, reading contents...", *stateFilePath)

		// read state file
		data, err := ioutil.ReadFile(*stateFilePath)
		if err != nil {
			log.Fatal().Err(err).Msgf("Failed reading file from path %v", *stateFilePath)
		}

		log.Info().Msgf("Unmarshalling file %v contents...", *stateFilePath)

		// unmarshal state file
		if err := json.Unmarshal(data, &state); err != nil {
			log.Fatal().Err(err).Interface("data", data).Msg("Failed unmarshalling state")
		}

		zoneInfoMap = state.ZoneInfoMap
	}

	return zoneInfoMap
}
