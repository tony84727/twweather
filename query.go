package twweather

// GetStationByTownName returns StationStatus that is matched by provided town name
func (weather *Weather) GetStationByTownName(townName string) (*StationStatus, error) {
	for _, location := range weather.stationList.Locations {
		if location.TownName == townName {
			// copy
			ret := location
			return &ret, nil
		}
	}
	return nil, nil
}
