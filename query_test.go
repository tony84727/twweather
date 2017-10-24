package twweather

import (
	"fmt"
	"testing"
)

func TestGetStationByTownName(t *testing.T) {
	weather.UpdateStationStatusWithData(sampleXML)
	stationMap := weather.GetStationByTownName("橫山鄉")
	station, ok := stationMap["橫山"]
	if !ok {
		t.Log(stationMap)
		t.Error("Cannot find station 橫山")
		if station.TownSN != 78 {
			t.Logf("Got Town number %v", station.TownSN)
			t.Fail()
		}
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
	stationMap := weather.GetStationsByCityName("新竹縣")
	logIfError := func(err error) {
		if err != nil {
			t.Error(err)
		}
	}
	if len(stationMap) != 2 {
		t.Fail()
	} else {
		logIfError(testHasStation(stationMap, "橫山"))
		logIfError(testHasStation(stationMap, "新豐"))
	}
}
