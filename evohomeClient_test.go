package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBasicAuthorizationToken(t *testing.T) {

	t.Run("ReturnsToken", func(t *testing.T) {
		client, _ := NewEvohomeClient()
		consumerKey := "***"
		consumerSecret := "***"

		basicAuthorizationToken := client.GetBasicAuthorizationToken(consumerKey, consumerSecret)

		assert.Equal(t, "", basicAuthorizationToken)
	})
}

func TestGetAccessToken(t *testing.T) {

	t.Run("ReturnsToken", func(t *testing.T) {

		client, _ := NewEvohomeClient()
		consumerKey := "***"
		consumerSecret := "***"

		// act
		accessToken, err := client.GetAccessToken(consumerKey, consumerSecret)

		if assert.Nil(t, err) {
			if assert.NotNil(t, accessToken) {
				assert.NotEqual(t, "", accessToken.Token)
			}
		}
	})
}
