package twweather

type StationMap map[string]StationStatus

type StationFilter func(StationStatus) bool

func (weather *Weather) GetStationBy(filter StationFilter) StationMap {
	candidate := make(StationMap, 2)
	for _, location := range weather.stationList.Locations {
		if filter(location) {
			// return a copy
			cp := location
			candidate[location.StationName] = cp
		}
	}
	return candidate
}

// GetStationByCityName returns StationMap that contains stations matched by town name.
func (weather *Weather) GetStationByTownName(townName string) StationMap {
	return weather.GetStationBy(func(station StationStatus) bool {
		return station.TownName == townName
	})
}

// GetStationByCityName returns StationMap that contains stations matched by city name.
func (weather *Weather) GetStationsByCityName(cityName string) StationMap {
	return weather.GetStationBy(func(station StationStatus) bool {
		return station.CityName == cityName
	})
}
