package twweather

import (
	"encoding/xml"
	"errors"
	"reflect"
)

// StationStatusDataID is dataid of cwb opendata API.
const StationStatusDataID = "O-A0001-001"

// Weather store data source and loaded station status.
type Weather struct {
	stationStatus *StationList
	cwbDataSource *cwbDataSource
}

// New return a initial weather struct without loading anything.
func New(cwbAPIKey string) *Weather {
	// create cwbDataSource
	weather := new(Weather)
	weather.cwbDataSource = &cwbDataSource{cwbAPIKey}
	return weather
}

// GetAvailableStationName returns a slice of available station names.
func (weather *Weather) GetAvailableStationName() []string {
	keys := reflect.ValueOf(weather.stationStatus.Locations).MapKeys()
	names := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		names[i] = keys[i].String()
	}
	return names
}

// LoadStationStatus reload station status.
func (weather *Weather) LoadStationStatus() (err error) {
	if weather.cwbDataSource == nil {
		err = errors.New("cwbDataSource haven't initialized")
		return
	}
	stationDataSet := weather.cwbDataSource.loadDataSet(StationStatusDataID)
	err = weather.UpdateStationStatusWithData(stationDataSet.RawData)
	return
}

func (weather *Weather) GetStation(name string) StationStatus {
	station := weather.stationStatus.Locations[name]
	return station
}

// UpdateStationStatusWithData update station status with a slice of byte.
func (weather *Weather) UpdateStationStatusWithData(data []byte) (err error) {
	newList := new(StationList)
	err = xml.Unmarshal(data, newList)
	if err != nil {
		return
	}
	weather.stationStatus = newList
	return
}
