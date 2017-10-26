package twweather

import "testing"

func TestGetTemperture(t *testing.T) {
	weather.UpdateStationStatusWithData(sampleXML)
	station := weather.GetStation("橫山")

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
