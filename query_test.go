package twweather

import "testing"

func TestGetStationByTownName(t *testing.T) {
	weather.UpdateStationStatusWithData(sampleXML)
	station, err := weather.GetStationByTownName("橫山鄉")
	if err != nil || station.TownSN != 78 {
		t.Fail()
	}
}
