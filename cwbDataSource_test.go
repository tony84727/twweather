package twweather

import (
	"encoding/xml"
	"fmt"
	"os"
	"testing"
)

var (
	sampleXML       []byte
	locationXML     []byte
	exampleElements = make(map[string]float64)
)

func init() {
	exampleElements["ELEV"] = 227
	exampleElements["WDIR"] = 56
	exampleElements["WDSD"] = 1.9
	exampleElements["TEMP"] = 26.6
	exampleElements["HUMD"] = 0.79
	exampleElements["PRES"] = 989.1
	exampleElements["SUN"] = -99
	exampleElements["H_24R"] = 0.0
	exampleElements["H_FX"] = -99
	exampleElements["H_XD"] = -99
	exampleElements["H_FXT"] = -99
}

type TestError struct {
	Message string
}

func (err TestError) Error() string {
	return err.Message
}

func createTestError(format string, params ...interface{}) *TestError {
	return &TestError{fmt.Sprintf(format, params)}
}

func matchExampleElements(t *testing.T, location *rawStationStatus) *TestError {
	convertedLocation := location.Convert()

	for name, expected := range exampleElements {
		element, ok := convertedLocation.WeatherElements[name]
		if !ok {
			return createTestError("Element %s not found!", name)
		}
		if element != expected {
			return createTestError("Element %s should be %f got %v!", name, expected, element)
		}
		t.Logf("Element match %s => %f = %f", name, expected, element)
	}
	return nil
}

// Test if we can unmarshal location xml with struct stationLocation
func TestParseLocation(t *testing.T) {
	station := new(rawStationStatus)
	xml.Unmarshal(locationXML, &station)
	if station.LocationName != "橫山" {
		t.Fail()
	}
	if count := len(station.RawWeatherElement); count != 11 {
		t.Logf("weather element count of the sample location should be 11. Got %d", count)
		t.Fail()
	}
	matchExampleElements(t, station)
}

func TestLoadData(t *testing.T) {
	t.Skip()
	weather.cwbDataSource = &cwbDataSource{os.Getenv("cwbAPIKey")}
	dataSet := weather.cwbDataSource.loadDataSet(StationStatusDataId)
	t.Log(string(dataSet.RawData))
}
