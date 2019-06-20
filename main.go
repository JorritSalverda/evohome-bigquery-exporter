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
	namespace             = kingpin.Flag("namespace", "Namespace the pod runs in.").Envar("NAMESPACE").Required().String()
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

	// check if session key exists in secret
	validSessionSecret := false
	var sessionSecret SessionSecret
	if _, err := os.Stat(*sessionSecretPath); !os.IsNotExist(err) {

		// read secret
		data, err := ioutil.ReadFile(*sessionSecretPath)
		if err != nil {
			log.Fatal().Err(err).Msgf("Failed reading file from path %v", *sessionSecretPath)
		}

		// unmarshal secret
		if err := json.Unmarshal(data, &sessionSecret); err != nil {
			log.Fatal().Err(err).Interface("data", data).Msg("Failed unmarshalling session secret")
		}

		// check if session secret isn't too old
		if sessionSecret.RetrievedAt.Add(time.Minute * time.Duration(*sessionTimeoutMinutes)).After(time.Now().UTC()) {
			validSessionSecret = true
		}
	}

	if !validSessionSecret {
		sessionID, userID, err := evoClient.GetSession(*username, *password)
		if err != nil {
			log.Fatal().Err(err).Msgf("Failed retrieving session id for username %v", *username)
		}

		sessionSecret = SessionSecret{
			SessionID:   sessionID,
			UserID:      userID,
			RetrievedAt: time.Now().UTC(),
		}

		// create kubernetes api client
		kubeClient, err := k8s.NewInClusterClient()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed creating Kubernetes API client")
		}

		// retrieve secret
		var secret corev1.Secret
		err = kubeClient.Get(context.Background(), *sessionSecretName, *namespace, &secret)
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
	}

	locations, err := evoClient.GetLocations(sessionSecret.SessionID, sessionSecret.UserID)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed retrieving locations for userid %v", sessionSecret.UserID)
	}

	log.Debug().Interface("locations", locations).Msgf("Retrieved %v locations: ", len(locations))

	// done
	log.Info().Msg("Finished exporting metrics")
}
