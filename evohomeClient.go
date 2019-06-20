package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sethgrid/pester"
)

// EvohomeClient is the interface for connecting to the evohome api
type EvohomeClient interface {
	GetSession(username, password string) (sessionID string, userID int, err error)
	GetLocations(sessionID string, userID int) (locations []LocationResponse, err error)
}

type evohomeClientImpl struct {
	baseURL string
}

// NewEvohomeClient returns new EvohomeClient
func NewEvohomeClient() (EvohomeClient, error) {
	return &evohomeClientImpl{
		baseURL: "https://tccna.honeywell.com",
	}, nil
}

func (ec *evohomeClientImpl) GetSession(username, password string) (sessionID string, userID int, err error) {
	// https://tccna.honeywell.com/WebAPI/api/Session

	// using this approach can suffer from rate limiting, see https://github.com/watchforstock/evohome-client/issues/57

	requestURL := ec.baseURL + "/WebAPI/api/Session"

	sessionRequest := SessionRequest{
		Username:      username,
		Password:      password,
		ApplicationID: "91db1612-73fd-4500-91b2-e63b069b185c",
	}
	sessionRequestJSONBytes, err := json.Marshal(sessionRequest)
	if err != nil {
		return
	}

	// create client, in order to add headers
	client := pester.New()
	client.MaxRetries = 3
	client.Backoff = pester.ExponentialJitterBackoff
	client.KeepLog = true
	client.Timeout = time.Second * 10
	request, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(sessionRequestJSONBytes))
	if err != nil {
		return
	}

	// add headers
	request.Header.Add("Content-Type", "application/json")

	// perform actual request
	response, err := client.Do(request)
	if err != nil {
		return
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	if response.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("Request to %v failed with status code %v: %v", requestURL, response.StatusCode, string(body))
	}

	log.Debug().Interface("body", body).Msg("Session response before unmarshalling")

	// unmarshal json body
	var sessionResponse SessionResponse
	err = json.Unmarshal(body, &sessionResponse)
	if err != nil {
		return
	}

	sessionID = sessionResponse.SessionID
	userID = sessionResponse.UserInfo.UserID

	return
}

func (ec *evohomeClientImpl) GetLocations(sessionID string, userID int) (locations []LocationResponse, err error) {
	// https://tccna.honeywell.com/WebAPI/api/locations?userId=%v&allData=True

	requestURL := ec.baseURL + fmt.Sprintf("/WebAPI/api/locations?userId=%v&allData=True", userID)

	// create client, in order to add headers
	client := pester.New()
	client.MaxRetries = 3
	client.Backoff = pester.ExponentialJitterBackoff
	client.KeepLog = true
	client.Timeout = time.Second * 10
	request, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return
	}

	// add headers
	request.Header.Add("sessionID", sessionID)

	// perform actual request
	response, err := client.Do(request)
	if err != nil {
		return
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	if response.StatusCode != http.StatusOK {
		return locations, fmt.Errorf("Request to %v failed with status code %v: %v", requestURL, response.StatusCode, string(body))
	}

	log.Debug().Interface("body", body).Msg("Location response before unmarshalling")

	// unmarshal json body
	err = json.Unmarshal(body, &locations)
	if err != nil {
		return
	}

	return
}
