package twweather

import (
	"fmt"
	"strconv"
)

type StationStatus struct {
	StationName string
	CityName    string
	CitySN      int
	TownName    string
	TownSN      int

	latitude  float64
	longitude float64

	WeatherElements map[string]interface{}
}

func (station StationStatus) String() string {
	_temperture, err := station.GetTemperture(true)
	temperture := "no data"
	if err == nil {
		temperture = strconv.FormatFloat(_temperture,'f',-1, 64)
	}

	_humidity, err := station.GetHumidity()
	humidity := "no data"
	if err == nil {
		humidity = strconv.Itoa(_humidity) + "%"
	}
		
	return fmt.Sprintf(`
		[[Station Info]]
		StationName: %s
		CityName: %s
		CitySN: %d
		TownName: %s
		TownSN: %d

		[[Weather Info]]
		Temperture: %s ℃
		Humidity: %s
`,
		station.StationName,
		station.CityName,
		station.CitySN,
		station.TownName,
		station.TownSN,
		temperture,
		humidity)
}

// testWeatherElementValid tests if there's data for a weather element.
func (s *StationStatus) testWeatherElementValid(key string) (result bool) {
	we, ok := s.WeatherElements[key]
	if !ok {
		return false
	}
	switch v := we.(type) {
	case float64:
		return v != -99
	}
	return true
}

func (s *StationStatus) getWeatherElement(key string, humanReadable string) (element interface{}, err error) {
	if !s.testWeatherElementValid(key) {
		panic(fmt.Errorf("no %s data", humanReadable))
	}
	element, _ = s.WeatherElements[key]
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	return
}

func (s *StationStatus) GetTemperture(celsius bool) (tempture float64, err error) {
	const key = "TEMP"
	we, err := s.getWeatherElement(key, "temperture")
	if err != nil {
		tempture = -99
		return
	}
	tempture = we.(float64)
	if !celsius {
		tempture = tempture*1.8 + 32
	}
	return
}

func (s *StationStatus) GetPressure() (hPa float64, err error) {
	const key = "PRES"
	we, err := s.getWeatherElement(key, "pressure")
	if err != nil {
		hPa = -99
		return
	}

	hPa = we.(float64)
	return
}

// GetHumidity returns percentage of relative humidity.
func (s *StationStatus) GetHumidity() (rh int, err error) {
	const key = "HUMD"
	we, err := s.getWeatherElement(key, "humidity")
	if err != nil {
		rh = -99
		return
	}
	// convert to percentage
	rh = int(we.(float64) * 100)
	return
}

// GetSunHours returns Sun hours
func (s *StationStatus) GetSunHours() (hours int, err error) {
	const key = "SUN"
	we, err := s.getWeatherElement(key, "sun hours")
	if err != nil {
		hours = -99
		return
	}
	hours = int(we.(float64))
	return
}

// GetDailyRainfall returns daily rainfall of the station in millimeters.
func (s *StationStatus) GetDailyRainfall(mm float64, err error) {
	const key = "H_24R"
	we, err := s.getWeatherElement(key, "daily rainfall")
	if err != nil {
		mm = -99
		return
	}
	mm = we.(float64)
	return
}

// GetMaximumWindSpeed returns maximum wind speed in meter per second (m/s).
func (s *StationStatus) GetMaximumWindSpeed() (speed float64, err error) {
	const key = "H_FX"
	we, err := s.getWeatherElement(key, "maximum wind speed")
	if err != nil {
		speed = -99
		return
	}
	speed = we.(float64)
	return
}
