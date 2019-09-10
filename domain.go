package main

import (
	"time"

	"cloud.google.com/go/bigquery"
)

// SessionRequest represents the json request body for POST https://tccna.honeywell.com/WebAPI/api/Session
type SessionRequest struct {
	Username      string `json:"Username"`
	Password      string `json:"Password"`
	ApplicationID string `json:"ApplicationId"`
}

// SessionResponse represents the response for POST https://tccna.honeywell.com/WebAPI/api/Session
type SessionResponse struct {
	SessionID string           `json:"sessionId"`
	UserInfo  UserInfoResponse `json:"userInfo"`
}

// UserInfoResponse see SessionResponse
type UserInfoResponse struct {
	UserID             int    `json:"userID"`
	Username           string `json:"username"`
	FirstName          string `json:"firstname"`
	LastName           string `json:"lastname"`
	StreetAddress      string `json:"streetAddress"`
	City               string `json:"city"`
	ZipCode            string `json:"zipcode"`
	Country            string `json:"country"`
	Telephone          string `json:"telephone"`
	UserLanguage       string `json:"userLanguage"`
	IsActivated        bool   `json:"isActivated"`
	DeviceCount        int    `json:"deviceCount"`
	TenantID           int    `json:"tenantID"`
	SecurityQuestion1  string `json:"securityQuestion1"`
	SecurityQuestion2  string `json:"securityQuestion2"`
	SecurityQuestion3  string `json:"securityQuestion3"`
	LatestEulaAccepted bool   `json:"latestEulaAccepted"`
}

// LocationResponse contains information for the evohome device location
type LocationResponse struct {
	LocationID                int                      `json:"locationID"`
	Name                      string                   `json:"name"`
	StreetAddress             string                   `json:"streetAddress"`
	City                      string                   `json:"city"`
	State                     string                   `json:"state"`
	Country                   string                   `json:"country"`
	ZipCode                   string                   `json:"zipcode"`
	Type                      string                   `json:"type"`
	HasStation                bool                     `json:"hasStation"`
	Devices                   []DeviceResponse         `json:"devices"`
	OneTouchButtons           []OneTouchButtonResponse `json:"oneTouchButtons"`
	Weather                   WeatherResponse          `json:"weather"`
	DaylightSavingTimeEnabled bool                     `json:"daylightSavingTimeEnabled"`
	TimeZone                  TimeZoneResponse         `json:"timeZone"`
	OneTouchActionsSuspended  bool                     `json:"oneTouchActionsSuspended"`
	IsLocationOwner           bool                     `json:"isLocationOwner"`
	LocationOwnerID           int                      `json:"locationOwnerID"`
	LocationOwnerName         string                   `json:"locationOwnerName"`
	LocationOwnerUserName     string                   `json:"locationOwnerUserName"`
	CanSearchForContractors   bool                     `json:"canSearchForContractors"`
}

// DeviceResponse contains additional information about the evohome device
type DeviceResponse struct {
	GatewayID           int                   `json:"gatewayId"`
	DeviceID            int                   `json:"deviceID"`
	ThermostatModelType string                `json:"thermostatModelType"`
	DeviceType          int                   `json:"deviceType"`
	Name                string                `json:"name"`
	ScheduleCapable     bool                  `json:"scheduleCapable"`
	HoldUntilCapable    bool                  `json:"holdUntilCapable"`
	Thermostat          ThermostatResponse    `json:"thermostat"`
	AlertSettings       AlertSettingsResponse `json:"alertSettings"`
	IsUpgrading         bool                  `json:"isUpgrading"`
	IsAlive             bool                  `json:"isAlive"`
	ThermostatVersion   string                `json:"thermostatVersion"`
	MacID               string                `json:"macID"`
	LocationID          int                   `json:"locationID"`
	DomainID            int                   `json:"domainID"`
	Instance            int                   `json:"instance"`
}

// ThermostatResponse contains information for each thermostat connected to the evohome device
type ThermostatResponse struct {
	Units                       string                   `json:"units"`
	IndoorTemperature           float64                  `json:"indoorTemperature"`
	OutdoorTemperature          float64                  `json:"outdoorTemperature"`
	OutdoorTemperatureAvailable bool                     `json:"outdoorTemperatureAvailable"`
	OutdoorHumidity             float64                  `json:"outdoorHumidity"`
	OutdootHumidityAvailable    bool                     `json:"outdootHumidityAvailable"`
	IndoorHumidity              float64                  `json:"indoorHumidity"`
	IndoorTemperatureStatus     string                   `json:"indoorTemperatureStatus"`
	IndoorHumidityStatus        string                   `json:"indoorHumidityStatus"`
	OutdoorTemperatureStatus    string                   `json:"outdoorTemperatureStatus"`
	OutdoorHumidityStatus       string                   `json:"outdoorHumidityStatus"`
	IsCommercial                bool                     `json:"isCommercial"`
	AllowedModes                []string                 `json:"allowedModes"`
	Deadband                    float64                  `json:"deadband"`
	MinHeatSetpoint             float64                  `json:"minHeatSetpoint"`
	MaxHeatSetpoint             float64                  `json:"maxHeatSetpoint"`
	MinCoolSetpoint             float64                  `json:"minCoolSetpoint"`
	MaxCoolSetpoint             float64                  `json:"maxCoolSetpoint"`
	ChangeableValues            ChangeableValuesResponse `json:"changeableValues"`
	ScheduleCapable             bool                     `json:"scheduleCapable"`
	VacationHoldChangeable      bool                     `json:"vacationHoldChangeable"`
	VacationHoldCancelable      bool                     `json:"vacationHoldCancelable"`
	ScheduleHeatSp              float64                  `json:"scheduleHeatSp"`
	ScheduleCoolSp              float64                  `json:"scheduleCoolSp"`
}

// ChangeableValuesResponse indicates the mode and other variables for a thermostat
type ChangeableValuesResponse struct {
	Mode             string               `json:"mode"`
	HeatSetpoint     HeatSetpointResponse `json:"heatSetpoint"`
	VacationHoldDays int                  `json:"vacationHoldDays"`
}

// HeatSetpointResponse indicates target temperatur
type HeatSetpointResponse struct {
	Value  float64 `json:"value"`
	Status string  `json:"status"`
}

// AlertSettingsResponse contains info about configured alerts
type AlertSettingsResponse struct {
	DeviceID                    int     `json:"deviceID"`
	TempHigherThanActive        bool    `json:"tempHigherThanActive"`
	TempHigherThan              float64 `json:"tempHigherThan"`
	TempHigherThanMinutes       int     `json:"tempHigherThanMinutes"`
	TempLowerThanActive         bool    `json:"tempLowerThanActive"`
	TempLowerThan               float64 `json:"tempLowerThan"`
	TempLowerThanMinutes        int     `json:"tempLowerThanMinutes"`
	FaultConditionExistsActive  bool    `json:"faultConditionExistsActive"`
	FaultConditionExistsHours   int     `json:"faultConditionExistsHours"`
	NormalConditionsActive      bool    `json:"normalConditionsActive"`
	CommunicationLostActive     bool    `json:"communicationLostActive"`
	CommunicationLostHours      int     `json:"communicationLostHours"`
	CommunicationFailureActive  bool    `json:"communicationFailureActive"`
	CommunicationFailureMinutes int     `json:"communicationFailureMinutes"`
	DeviceLostActive            bool    `json:"deviceLostActive"`
	DeviceLostHours             int     `json:"deviceLostHours"`
}

// OneTouchButtonResponse is not something we're interested in
type OneTouchButtonResponse struct {
}

// WeatherResponse contains info about the weather conditions at the location
type WeatherResponse struct {
	Condition   string  `json:"condition"`
	Temperature float64 `json:"temperature"`
	Units       string  `json:"units"`
	Humidity    float64 `json:"humidity"`
	Phrase      string  `json:"phrase"`
}

// TimeZoneResponse contains info about the timezone of a location
type TimeZoneResponse struct {
	ID                      string `json:"id"`
	DisplayName             string `json:"displayName"`
	OffsetMinutes           int    `json:"offsetMinutes"`
	CurrentOffsetMinutes    int    `json:"currentOffsetMinutes"`
	UsingDaylightSavingTime bool   `json:"usingDaylightSavingTime"`
}

type SessionSecret struct {
	SessionID   string
	UserID      int
	RetrievedAt time.Time
}

type BigQueryMeasurement struct {
	Location   string         `bigquery:"location"`
	MeasuredAt time.Time      `bigquery:"measured_at"`
	Zones      []BigQueryZone `bigquery:"zones"`
	InsertedAt time.Time      `bigquery:"inserted_at"`
}

type BigQueryZone struct {
	Zone              string               `bigquery:"location"`
	TemperatureUnit   string               `bigquery:"unit"`
	TemperatureValue  bigquery.NullFloat64 `bigquery:"temperature"`
	HeatSetPointValue bigquery.NullFloat64 `bigquery:"heat_setpoint"`
	HeatDemandValue   bigquery.NullFloat64 `bigquery:"heat_demand"`
	HumidityValue     bigquery.NullFloat64 `bigquery:"humidity"`
}

// State from evohome-hgi80-listener
type State struct {
	ZoneInfoMap map[int64]ZoneInfo
}

type ZoneInfo struct {
	ID             int64
	Name           string
	MinTemperature float64
	MaxTemperature float64
	Temperature    float64
	Setpoint       float64
	HeatDemand     float64
}

func (z ZoneInfo) IsActualZone() bool {
	return z.ID < 12 && z.Name != ""
}
