package twweather

import (
	"fmt"
)

// GetStationByTownName returns StationStatus that is matched by provided town name
func (weather *Weather) GetStationByTownName(townName string) (*StationStatus, error) {
	for _, location := range weather.stationList.Locations {
		if location.TownName == townName {
			// copy
			ret := location
			return &ret, nil
		}
	}
	return nil, fmt.Errorf("Cannot find station by town name = %s", townName)
}

func (weather *Weather) GetStationsByCityName(cityName string) (map[string]StationStatus, error) {
	stations := make(map[string]StationStatus, 2)
	for _, location := range weather.stationList.Locations {
		if location.CityName == cityName {
			stations[location.StationName] = location
		}
	}
	if len(stations) == 0 {
		return stations, fmt.Errorf("Cannot find station by city name = %s", cityName)
	}
	return stations, nil
}
