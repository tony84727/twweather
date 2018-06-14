package twweather

import (
	"encoding/xml"
	"testing"

	"github.com/MinecraftXwinP/twweather/testutil"
	"github.com/stretchr/testify/assert"
)

var (
	sampleXML   []byte
	locationXML []byte
)

func TestMain(m *testing.M) {
	sampleXML = testutil.Load("sample.xml")
	locationXML = testutil.Load("location.xml")
	m.Run()
}

func TestParseObservation(t *testing.T) {
	observation := new(Observation)
	xml.Unmarshal(locationXML, observation)
	assert.Equal(t, float64(227), observation.WeatherElements["ELEV"])
	assert.Equal(t, float64(56), observation.WeatherElements["WDIR"])
	assert.Equal(t, float64(1.9), observation.WeatherElements["WDSD"])
	assert.Equal(t, float64(26.6), observation.WeatherElements["TEMP"])
	assert.Equal(t, float64(0.79), observation.WeatherElements["HUMD"])
	assert.Equal(t, float64(989.1), observation.WeatherElements["PRES"])
	assert.Equal(t, float64(-99), observation.WeatherElements["SUN"])
	assert.Equal(t, float64(0.0), observation.WeatherElements["H_24R"])
	assert.Equal(t, float64(-99), observation.WeatherElements["H_FX"])
	assert.Equal(t, float64(-99), observation.WeatherElements["H_XD"])
	assert.Equal(t, "新竹縣", observation.CityName)
	assert.Equal(t, 10, observation.CitySN)
	assert.Equal(t, "橫山鄉", observation.TownName)
	assert.Equal(t, 78, observation.TownSN)
}

func TestObservationQueries(t *testing.T) {
	ob := new(Observation)
	xml.Unmarshal(locationXML, ob)
	temp, err := ob.GetTemperature(true)
	assert.NoError(t, err)
	assert.Equal(t, float64(26.6), temp)
	temp, err = ob.GetTemperature(false)
	assert.NoError(t, err)
	assert.Equal(t, float64(79.88), temp)
	_, err = ob.GetSunHours()
	assert.Error(t, err)
	rainfail, err := ob.GetDailyRainfall()
	assert.NoError(t, err)
	assert.Equal(t, float64(0.0), rainfail)
	humd, err := ob.GetHumidity()
	assert.NoError(t, err)
	assert.Equal(t, 79, humd)
	pres, err := ob.GetPressure()
	assert.NoError(t, err)
	assert.Equal(t, float64(989.1), pres)
}
