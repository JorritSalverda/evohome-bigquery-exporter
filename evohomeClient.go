package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sethgrid/pester"
)

// EvohomeClient is the interface for connecting to the evohome api
type EvohomeClient interface {
	// GetBasicAuthorizationToken returns base 64 encoded combination of consumer key and secret
	GetBasicAuthorizationToken(consumerKey, consumerSecret string) (basicAuthorizationToken string)
	// GetAccessToken retrieves an access token for the evohome api using app consumer key and secret and grant type client_credentials
	GetAccessToken(consumerKey, consumerSecret string) (accessToken *AccessToken, err error)
}

type evohomeClientImpl struct {
	baseURL string
}

// NewEvohomeClient returns new EvohomeClient
func NewEvohomeClient() (EvohomeClient, error) {
	return &evohomeClientImpl{
		baseURL: "https://api.honeywell.com",
	}, nil
}

func (ec *evohomeClientImpl) GetBasicAuthorizationToken(consumerKey, consumerSecret string) (basicAuthorizationToken string) {
	credentials := fmt.Sprintf("%v:%v", consumerKey, consumerSecret)
	basicAuthorizationToken = base64.StdEncoding.EncodeToString([]byte(credentials))

	return
}

func (ec *evohomeClientImpl) GetAccessToken(consumerKey, consumerSecret string) (accessToken *AccessToken, err error) {
	// https://developer.honeywell.com/authorization-oauth2/apis/post/accesstoken

	requestURL := ec.baseURL + "/oauth2/accesstoken"

	formData := url.Values{
		"grant_type": []string{"client_credentials"},
	}

	// create client, in order to add headers
	client := pester.New()
	client.MaxRetries = 3
	client.Backoff = pester.ExponentialJitterBackoff
	client.KeepLog = true
	client.Timeout = time.Second * 10
	request, err := http.NewRequest("POST", requestURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return
	}

	// The Authorization HTTP header for this request is a Base64 encoded value of apikey and apiSecret concatenated with a colon. For example if your API Key was 123abc and your Secret was 456def your HTTP header would look like this:
	//curl -X POST -H "Authorization: Basic MTIzYWJjOjQ1NmRlZg==" -H "Content-Type: application/x-www-form-urlencoded" -d 'grant_type=client_credentials'

	basicAuthorizationToken := ec.GetBasicAuthorizationToken(consumerKey, consumerSecret)

	// add headers
	request.Header.Add("Authorization", fmt.Sprintf("Basic %v", basicAuthorizationToken))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

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
		return nil, fmt.Errorf("Request to %v failed with status code %v: %v", requestURL, response.StatusCode, string(body))
	}

	log.Print(body)

	// unmarshal json body
	var accessTokenResponse AccessTokenResponse
	err = json.Unmarshal(body, &accessTokenResponse)
	if err != nil {
		return
	}

	return &AccessToken{Token: accessTokenResponse.AccessToken}, nil
}
