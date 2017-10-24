package twweather

import (
	"fmt"
	"testing"
)

func TestGetStationByTownName(t *testing.T) {
	weather.UpdateStationStatusWithData(sampleXML)
	station, err := weather.GetStationByTownName("橫山鄉")
	if err != nil || station.TownSN != 78 {
		t.Fail()
	}
}

func testHasStation(list map[string]StationStatus, name string) error {
	_, ok := list[name]
	if !ok {
		return fmt.Errorf("Missing station %s", name)
	}
	return nil
}

func TestGetStationsByCityName(t *testing.T) {
	weather.UpdateStationStatusWithData(sampleXML)
	stations, err := weather.GetStationsByCityName("新竹縣")
	logIfError := func(err error) {
		if err != nil {
			t.Error(err)
		}
	}
	if err != nil || len(stations) != 2 {
		t.Fail()
	} else {
		logIfError(testHasStation(stations, "橫山"))
		logIfError(testHasStation(stations, "新豐"))
	}
}
