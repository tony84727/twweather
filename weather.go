package twweather

import (
	"encoding/xml"
	"reflect"
)

const StationStatusDataId = "O-A0001-001"

type Weather struct {
	stationStatus *stationList
	cwbDataSource *cwbDataSource
}

// Create a weather object.
func New(cwbAPIKey string) *Weather {
	// create cwbDataSource
	weather := new(Weather)
	weather.cwbDataSource = &cwbDataSource{cwbAPIKey}
	return weather
}

func (weather *Weather) GetAvailableStationName() []string {
	keys := reflect.ValueOf(weather.stationStatus.Locations).MapKeys()
	names := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		names[i] = keys[i].String()
	}
	return names
}

func (weather *Weather) LoadStationStatus() (err error) {
	if weather.cwbDataSource == nil {
		err = LoadError{"cwbDataSource haven't initialized"}
		return
	}
	stationDataSet := weather.cwbDataSource.loadDataSet(StationStatusDataId)
	err = weather.UpdateStationStatusWithData(stationDataSet.RawData)
	return
}

// Decode the data (xml) and update station stauts
func (weather *Weather) UpdateStationStatusWithData(data []byte) (err error) {
	_rawStationList := new(rawStationList)
	err = xml.Unmarshal(data, _rawStationList)
	if err != nil {
		return
	}
	weather.stationStatus = _rawStationList.Convert()
	return
}

type LoadError struct {
	Message string
}

func (err LoadError) Error() string {
	return err.Message
}
