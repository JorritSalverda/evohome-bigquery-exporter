package main

import (
	"time"

	"cloud.google.com/go/bigquery"
)

func mapLocationsToMeasurements(locations []LocationResponse, outdoorZoneName string, zoneInfoMap map[int64]ZoneInfo) (measurements []BigQueryMeasurement) {
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

			heatDemandValue := bigquery.NullFloat64{Valid: false}
			zoneInfo := getZoneInfoFromMapByName(zoneInfoMap, d.Name)
			if zoneInfo != nil {
				heatDemandValue = bigquery.NullFloat64{Float64: zoneInfo.HeatDemand, Valid: true}
			}

			zone := BigQueryZone{
				Zone:              d.Name,
				TemperatureUnit:   d.Thermostat.Units,
				TemperatureValue:  bigquery.NullFloat64{Float64: d.Thermostat.IndoorTemperature, Valid: true},
				HeatSetPointValue: bigquery.NullFloat64{Float64: d.Thermostat.ChangeableValues.HeatSetpoint.Value, Valid: true},
				HeatDemandValue:   heatDemandValue,
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

func getZoneInfoFromMapByName(zoneInfoMap map[int64]ZoneInfo, zoneName string) *ZoneInfo {
	for _, v := range zoneInfoMap {
		if v.Name == zoneName {
			return &v
		}
	}

	return nil
}
