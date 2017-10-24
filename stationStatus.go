package twweather

import "errors"

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

func (s *StationStatus) GetTemperture(celsius bool) (tempture float64, err error) {
	const key = "TEMP"

	if !s.testWeatherElementValid(key) {
		panic(errors.New("no temperture data"))
	}
	we, _ := s.WeatherElements[key]
	tempture = we.(float64)

	if !celsius {
		tempture = tempture*1.8 + 32
	}
	return
}

func (s *StationStatus) GetPressure() (hPa float64, err error) {
	const key = "PRES"

	if !s.testWeatherElementValid(key) {
		panic(errors.New("no pressure data"))
	}
	we, _ := s.WeatherElements[key]

	hPa = we.(float64)
	return
}
