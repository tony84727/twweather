package twweather

import (
	"github.com/MinecraftXwinP/twweather/cwbdata"
	"encoding/xml"
	"fmt"
	"strconv"
	"time"
)
const StationStatusDataID = "O-A0001-001"
type StationStatusList map[string]StationStatus

func (list *StationStatusList) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	origin := new(struct {
		Locations StationStatus `xml:"location"`
	})

	err := d.DecodeElement(origin, &start)
	if err != nil {
		return err
	}
	for _, location := range origin.Locations {
		list[location.StationName] = location
	}
	return nil
}
func GetStationList(dataSource cwbdata.OpenDataSource) *StationStatusList,error {
	openData,err := dataSource.GetOpenData(StationStatusDataID)
	if err != nil {
		return nil,err
	}
	list := make(StationStatusList, 100)
	err = xml.Unmarshal(openData.DataSet, &list)
	return list,err

}


func (list *StationStatusList) GetAvailableStationNames() []string {
	names := make([]string,100)
	for name,_ := range list {
		names[] = name
	}
	return names
}

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
		temperture = strconv.FormatFloat(_temperture, 'f', -1, 64)
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

type rawWeatherElement struct {
	Name  string
	Value interface{}
}

func (rawElement *rawWeatherElement) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	const timeShowFormat = "2006-01-02T15:04:05-07:00"
	raw := new(struct {
		Name  string `xml:"elementName"`
		Value string `xml:"elementValue>value"`
	})
	err := d.DecodeElement(raw, &start)
	if err != nil {
		return err
	}
	rawElement.Name = raw.Name

	valStr := raw.Value
	timeStamp, err := time.Parse(timeShowFormat, valStr)
	if err != nil {
		f, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			return err
		}
		rawElement.Value = f
	} else {
		rawElement.Value = timeStamp
	}
	return nil
}

func (status *StationStatus) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	raw := new(struct {
		StationName     string              `xml:"locationName"`
		Latitude        float64             `xml:"lat"`
		Longitude       float64             `xml:"lon"`
		WeatherElements []rawWeatherElement `xml:"weatherElement"`
		Parameters      []struct {
			Name  string `xml:"parameterName"`
			Value string `xml:"parameterValue"`
		} `xml:"parameter"`
	})
	err := d.DecodeElement(raw, &start)
	if err != nil {
		return err
	}
	status.StationName = raw.StationName
	status.latitude = raw.Latitude
	status.longitude = raw.Longitude
	// init map
	status.WeatherElements = make(map[string]interface{}, 11)
	for _, element := range raw.WeatherElements {
		status.WeatherElements[element.Name] = element.Value
	}
	for _, parameter := range raw.Parameters {
		switch parameter.Name {
		case "CITY":
			status.CityName = parameter.Value
			break
		case "CITY_SN":
			i, err := strconv.Atoi(parameter.Value)
			if err == nil {
				status.CitySN = i
			}
			break
		case "TOWN":
			status.TownName = parameter.Value
			break
		case "TOWN_SN":
			i, err := strconv.Atoi(parameter.Value)
			if err == nil {
				status.TownSN = i
			}
			break
		}
	}
	return nil
}
