package twweather

import (
	"encoding/xml"
	"fmt"
	"testing"
	"time"

	"github.com/MinecraftXwinP/twweather/testutil"
)

var (
	sampleXML       []byte
	locationXML     []byte
	exampleElements = make(map[string]interface{})
)

func TestMain(m *testing.M) {
	// load sample data
	sampleXML = testutil.Load("sample.xml")
	locationXML = testutil.Load("location.xml")
	weather = New("FAKEKEY")
	m.Run()
}

func init() {
	exampleElements["ELEV"] = float64(227)
	exampleElements["WDIR"] = float64(56)
	exampleElements["WDSD"] = float64(1.9)
	exampleElements["TEMP"] = float64(26.6)
	exampleElements["HUMD"] = float64(0.79)
	exampleElements["PRES"] = float64(989.1)
	exampleElements["SUN"] = float64(-99)
	exampleElements["H_24R"] = float64(0.0)
	exampleElements["H_FX"] = float64(-99)
	exampleElements["H_XD"] = float64(-99)
	exampleElements["H_FXT"] = time.Date(2017, 10, 19, 7, 29, 0, 0, time.FixedZone("CST", 8*60*60))
}

func createTestError(format string, params ...interface{}) error {
	return fmt.Errorf(format, params...)
}

func matchExampleElements(t *testing.T, station *StationStatus) error {
	for name, expected := range exampleElements {
		element, ok := station.WeatherElements[name]
		if !ok {
			return createTestError("Element %s not found!", name)
		}
		switch v := expected.(type) {
		case time.Time:
			if !v.Equal(element.(time.Time)) {
				return createTestError("Element %s should be %v got %v!", name, expected, element)
			}
			break
		default:
			if element != expected {
				return createTestError("Element %s should be %v got %v!", name, expected, element)
			}
			break
		}

		t.Logf("Element match %s => %v = %v", name, expected, element)
	}
	return nil
}

// Test if we can unmarshal location xml with struct stationLocation
func TestParseLocation(t *testing.T) {
	location := new(StationStatus)
	err := xml.Unmarshal(locationXML, location)
	if err != nil {
		t.Fatal(err)
	}
	if location.StationName != "橫山" {
		t.Fail()
	}
	if count := len(location.WeatherElements); count != 11 {
		t.Logf("weather element count of the sample location should be 11. Got %d", count)
		t.Fail()
	}
	err = matchExampleElements(t, location)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

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
