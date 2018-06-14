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

// func (list *StationStatusList) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
// 	origin := new(struct {
// 		Locations StationStatus `xml:"location"`
// 	})

// 	err := d.DecodeElement(origin, &start)
// 	if err != nil {
// 		return err
// 	}
// 	for _, location := range origin.Locations {
// 		list[location.StationName] = location
// 	}
// 	return nil
// }
// func GetStationList(dataSource cwbdata.OpenDataSource) (*StationStatusList,error) {
// 	openData,err := dataSource.GetOpenData(StationStatusDataID)
// 	if err != nil {
// 		return nil,err
// 	}
// 	list := make(StationStatusList, 100)
// 	err = xml.Unmarshal(openData.DataSet, &list)
// 	return list,err

// }


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



