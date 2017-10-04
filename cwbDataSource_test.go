package twweather

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var (
	sampleXML   []byte
	locationXML []byte
)

func load(path string) []byte {
	filePath := "./testdata/" + path
	buffer, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Cannot load sample data: %s", filePath)
		os.Exit(1)
	}
	return buffer
}

func TestMain(m *testing.M) {
	// load sample data
	sampleXML = load("sample.xml")
	locationXML = load("location.xml")
	m.Run()
}

// Test if we can unmarshal location xml with struct stationLocation
func TestParseLocation(t *testing.T) {
	location := new(stationLocation)
	xml.Unmarshal(locationXML, &location)
	if location.LocationName != "橫山" {
		t.Fail()
	}
	if count := len(location.WeatherElements); count != 11 {
		t.Logf("weather element count of the sample location should be 11. Got %d", count)
		t.Fail()
	}
	for _, element := range location.WeatherElements {
		t.Logf("%s => %s", element.ElementName, element.ElementValue)
	}
}
