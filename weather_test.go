package twweather

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var (
	weather *Weather
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
	weather = New("FAKEKEY")
	m.Run()
}

func TestParseStationStatus(t *testing.T) {
	err := weather.UpdateStationStatusWithData(sampleXML)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestGetAvailableStationName(t *testing.T) {
	names := weather.GetAvailableStationName()
	t.Log(names)
	if count := len(names); count != 2 {
		t.Logf("Should have 2 station name. Expect 2, got %d", count)
	}
}
