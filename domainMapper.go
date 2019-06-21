package main

import (
	"time"

	"cloud.google.com/go/bigquery"
)

func mapLocationsToMeasurements(locations []LocationResponse, outdoorZoneName string) (measurements []BigQueryMeasurement) {
	measurements = []BigQueryMeasurement{}

	for _, l := range locations {
		measurement := BigQueryMeasurement{
			Location:   l.Name,
			MeasuredAt: time.Now().UTC(),
			Zones:      []BigQueryZone{},
			InsertedAt: time.Now().UTC(),
		}

		// loop devices
		for _, d := range l.Devices {
			zone := BigQueryZone{
				Zone:              d.Name,
				TemperatureUnit:   d.Thermostat.Units,
				TemperatureValue:  bigquery.NullFloat64{Float64: d.Thermostat.IndoorTemperature, Valid: true},
				HeatSetPointValue: bigquery.NullFloat64{Float64: d.Thermostat.ChangeableValues.HeatSetpoint.Value, Valid: true},
			}
			measurement.Zones = append(measurement.Zones, zone)
		}

		// add weather as zone
		measurement.Zones = append(measurement.Zones, BigQueryZone{
			Zone:             outdoorZoneName,
			TemperatureUnit:  l.Weather.Units,
			TemperatureValue: bigquery.NullFloat64{Float64: l.Weather.Temperature, Valid: true},
			HumidityValue:    bigquery.NullFloat64{Float64: l.Weather.Humidity, Valid: true},
		})

		measurements = append(measurements, measurement)
	}

	return
}
