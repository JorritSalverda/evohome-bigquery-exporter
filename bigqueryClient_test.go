package main

import (
	"os"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/stretchr/testify/assert"
)

func TestCheckIfDatasetExists(t *testing.T) {

	t.Run("ReturnsTrueIfDatasetExists", func(t *testing.T) {

		if testing.Short() {
			t.Skip("skipping test in short mode.")
		}

		client, _ := NewBigQueryClient(os.Getenv("BQ_PROJECT_ID"))

		// act
		exists := client.CheckIfDatasetExists(os.Getenv("BQ_DATASET"))

		assert.True(t, exists)
	})
}

func TestCheckIfTableExists(t *testing.T) {

	t.Run("ReturnsFalseIfTableDoesNotExist", func(t *testing.T) {

		if testing.Short() {
			t.Skip("skipping test in short mode.")
		}

		client, _ := NewBigQueryClient(os.Getenv("BQ_PROJECT_ID"))

		// act
		exists := client.CheckIfTableExists(os.Getenv("BQ_DATASET"), "evohome_test")

		assert.False(t, exists)
	})

	t.Run("ReturnsTrueIfTableExists", func(t *testing.T) {

		if testing.Short() {
			t.Skip("skipping test in short mode.")
		}

		client, _ := NewBigQueryClient(os.Getenv("BQ_PROJECT_ID"))
		client.CreateTable(os.Getenv("BQ_DATASET"), "evohome_test", BigQueryMeasurement{}, "measured_at", true)

		// act
		exists := client.CheckIfTableExists(os.Getenv("BQ_DATASET"), "evohome_test")

		assert.True(t, exists)

		// cleanup
		client.DeleteTable(os.Getenv("BQ_DATASET"), "evohome_test")
	})
}

// export BQ_PROJECT_ID=...
// export BQ_DATASET=...
func TestCreateTable(t *testing.T) {

	t.Run("CreatesTableIfTableDoesNotExistYet", func(t *testing.T) {

		if testing.Short() {
			t.Skip("skipping test in short mode.")
		}

		client, _ := NewBigQueryClient(os.Getenv("BQ_PROJECT_ID"))

		// act
		err := client.CreateTable(os.Getenv("BQ_DATASET"), "evohome_test", BigQueryMeasurement{}, "measured_at", false)

		assert.Nil(t, err)

		// cleanup
		client.DeleteTable(os.Getenv("BQ_DATASET"), "evohome_test")
	})

	t.Run("ErrorsIfTableAlreadyExists", func(t *testing.T) {

		if testing.Short() {
			t.Skip("skipping test in short mode.")
		}

		client, _ := NewBigQueryClient(os.Getenv("BQ_PROJECT_ID"))
		client.CreateTable(os.Getenv("BQ_DATASET"), "evohome_test", BigQueryMeasurement{}, "measured_at", false)

		// act
		err := client.CreateTable(os.Getenv("BQ_DATASET"), "evohome_test", BigQueryMeasurement{}, "measured_at", false)

		assert.NotNil(t, err)

		// cleanup
		client.DeleteTable(os.Getenv("BQ_DATASET"), "evohome_test")
	})
}

func TestInsertMeasurements(t *testing.T) {

	t.Run("InsertMeasurementWithNullValuesForTemperatureAndHeatSetpointAndHumidityForAllZones", func(t *testing.T) {

		if testing.Short() {
			t.Skip("skipping test in short mode.")
		}

		client, _ := NewBigQueryClient(os.Getenv("BQ_PROJECT_ID"))
		client.CreateTable(os.Getenv("BQ_DATASET"), "evohome_test", BigQueryMeasurement{}, "measured_at", true)
		measurements := []BigQueryMeasurement{
			BigQueryMeasurement{
				Location:   "here",
				MeasuredAt: time.Now().UTC(),
				Zones: []BigQueryZone{
					BigQueryZone{
						Zone:              "room 1",
						TemperatureUnit:   "Celsius",
						TemperatureValue:  bigquery.NullFloat64{},
						HeatSetPointValue: bigquery.NullFloat64{},
						HumidityValue:     bigquery.NullFloat64{},
					},
					BigQueryZone{
						Zone:              "room 2",
						TemperatureUnit:   "Celsius",
						TemperatureValue:  bigquery.NullFloat64{},
						HeatSetPointValue: bigquery.NullFloat64{},
						HumidityValue:     bigquery.NullFloat64{},
					},
					BigQueryZone{
						Zone:              "room 3",
						TemperatureUnit:   "Celsius",
						TemperatureValue:  bigquery.NullFloat64{},
						HeatSetPointValue: bigquery.NullFloat64{},
						HumidityValue:     bigquery.NullFloat64{},
					},
				},
				InsertedAt: time.Now().UTC(),
			},
		}

		// act
		err := client.InsertMeasurements(os.Getenv("BQ_DATASET"), "evohome_test", measurements)

		assert.Nil(t, err)

		// cleanup
		client.DeleteTable(os.Getenv("BQ_DATASET"), "evohome_test")
	})

	t.Run("InsertMeasurementWithValuesForTemperatureAndHeatSetpointAndHumidityForAllZones", func(t *testing.T) {

		if testing.Short() {
			t.Skip("skipping test in short mode.")
		}

		client, _ := NewBigQueryClient(os.Getenv("BQ_PROJECT_ID"))
		client.CreateTable(os.Getenv("BQ_DATASET"), "evohome_test", BigQueryMeasurement{}, "measured_at", true)
		measurements := []BigQueryMeasurement{
			BigQueryMeasurement{
				Location:   "here",
				MeasuredAt: time.Now().UTC(),
				Zones: []BigQueryZone{
					BigQueryZone{
						Zone:              "room 1",
						TemperatureUnit:   "Celsius",
						TemperatureValue:  bigquery.NullFloat64{Float64: 19.6, Valid: true},
						HeatSetPointValue: bigquery.NullFloat64{Float64: 20.0, Valid: true},
						HumidityValue:     bigquery.NullFloat64{},
					},
					BigQueryZone{
						Zone:              "room 2",
						TemperatureUnit:   "Celsius",
						TemperatureValue:  bigquery.NullFloat64{Float64: 17.3, Valid: true},
						HeatSetPointValue: bigquery.NullFloat64{Float64: 15.0, Valid: true},
						HumidityValue:     bigquery.NullFloat64{},
					},
					BigQueryZone{
						Zone:              "outdoor",
						TemperatureUnit:   "Celsius",
						TemperatureValue:  bigquery.NullFloat64{Float64: 15.6, Valid: true},
						HeatSetPointValue: bigquery.NullFloat64{},
						HumidityValue:     bigquery.NullFloat64{Float64: 45.5, Valid: true},
					},
				},
				InsertedAt: time.Now().UTC(),
			},
		}

		// act
		err := client.InsertMeasurements(os.Getenv("BQ_DATASET"), "evohome_test", measurements)

		assert.Nil(t, err)

		// cleanup
		client.DeleteTable(os.Getenv("BQ_DATASET"), "evohome_test")
	})
}

func TestEndToEnd(t *testing.T) {

	t.Run("RetrieveAndInsertsMeasurements", func(t *testing.T) {

		if testing.Short() {
			t.Skip("skipping test in short mode.")
		}

		bqClient, _ := NewBigQueryClient(os.Getenv("BQ_PROJECT_ID"))
		bqClient.CreateTable(os.Getenv("BQ_DATASET"), "evohome_test", BigQueryMeasurement{}, "measured_at", true)
		evoClient, _ := NewEvohomeClient()

		// act
		sessionID, userID, _ := evoClient.GetSession(os.Getenv("EVOHOME_USERNAME"), os.Getenv("EVOHOME_PASSWORD"))
		locations, _ := evoClient.GetLocations(sessionID, userID)

		measurements := mapLocationsToMeasurements(locations, "outside")

		// act
		err := bqClient.InsertMeasurements(os.Getenv("BQ_DATASET"), "evohome_test", measurements)

		assert.Nil(t, err)

		// cleanup
		bqClient.DeleteTable(os.Getenv("BQ_DATASET"), "evohome_test")
	})
}
