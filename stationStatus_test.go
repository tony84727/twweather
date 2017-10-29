package twweather

import "testing"

func TestGetTemperture(t *testing.T) {
	weather.UpdateStationStatusWithData(sampleXML)
	station, err := weather.GetStation("橫山")

	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}

	temperture, err := station.GetTemperture(true)

	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}

	if temperture != 26.6 {
		t.Logf("Should got 26.6, got %f", temperture)
		t.Fail()
	}
}

func TesttestWeatherElementValid(t *testing.T) {
	weather.UpdateStationStatusWithData(sampleXML)
	station, err := weather.GetStation("橫山")
	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}
	isValid := station.testWeatherElementValid("SUM")
	if isValid {
		t.Errorf("SUM data of 橫山 should be invalid")
	}
	isValid = station.testWeatherElementValid("H_FXT")
	if !isValid {
		t.Errorf("H_FXT data of 橫山 should be valid")
	}
}
