package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSession(t *testing.T) {

	t.Run("ReturnsSessionIDAndUserID", func(t *testing.T) {

		client, _ := NewEvohomeClient()
		username := "jorrit.salverda@gmail.com"
		password := "hRLoHsd8Arca93sNxH$k%kf5z17ta6"

		// act
		sessionID, userID, err := client.GetSession(username, password)

		if assert.Nil(t, err) {
			assert.NotEqual(t, "", sessionID)
			assert.Equal(t, 2625379, userID)
		}
	})
}

func TestGetLocations(t *testing.T) {

	t.Run("ReturnsAllLocationsForUser", func(t *testing.T) {

		if testing.Short() {
			t.Skip("skipping test in short mode.")
		}

		client, _ := NewEvohomeClient()
		username := "***"
		password := "***"
		sessionID, userID, _ := client.GetSession(username, password)

		// act
		locations, err := client.GetLocations(sessionID, userID)

		if assert.Nil(t, err) {
			assert.Equal(t, 1, len(locations))
			assert.Equal(t, "Thuis", locations[0].Name)
			assert.Equal(t, 6, len(locations[0].Devices))
			assert.Equal(t, "Badkamers", locations[0].Devices[0].Name)
			assert.Equal(t, "Logeerkamer", locations[0].Devices[1].Name)
			assert.Equal(t, "Slaapkamers", locations[0].Devices[2].Name)
			assert.Equal(t, "Studeerkamer", locations[0].Devices[3].Name)
			assert.Equal(t, "Washok", locations[0].Devices[4].Name)
			assert.Equal(t, "Woonkamer", locations[0].Devices[5].Name)
			assert.Equal(t, "Celsius", locations[0].Devices[5].Thermostat.Units)
			assert.Equal(t, 20.7800, locations[0].Devices[5].Thermostat.IndoorTemperature)
			assert.Equal(t, 20.0, locations[0].Devices[5].Thermostat.ChangeableValues.HeatSetpoint.Value)
		}
	})
}
