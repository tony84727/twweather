package forecast

import (
	"strconv"
)

func calcWeeklyForecastDatasetID(base int, city int) string {
	return strconv.Itoa(city*4 + base)
}

func getTwoWeekDatasetID(city int) string {
	return "F-D0047-0" + calcWeeklyForecastDatasetID(1, city)
}

func getOneWeekDatasetID(city int) string {
	return "F-D0047-0" + calcWeeklyForecastDatasetID(3, city)
}

type WeelyForecast struct {
}
