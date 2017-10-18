package twweather

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
)

const ApiUrl = "http://opendata.cwb.gov.tw/opendataapi"

type cwbDataSource struct {
	APIKey string
}

type cwbDataSet struct {
	RawData []byte
	DataID  string
}

func initDataSource(apiKey string) cwbDataSource {
	return cwbDataSource{APIKey: apiKey}
}

func (cwb cwbDataSource) loadDataSet(dataID string) (result cwbDataSet) {
	result = cwbDataSet{DataID: dataID}
	response, err := http.Get(fmt.Sprintf("%s?dataid=%s&authorizationkey=%s", ApiUrl, dataID, cwb.APIKey))
	if err != nil {
		log.Fatal(err)
		return
	}
	buffer := new(bytes.Buffer)
	defer response.Body.Close()
	buffer.ReadFrom(response.Body)
	result.RawData = buffer.Bytes()
	return
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

type StationList struct {
	Locations map[string]StationStatus `xml:"location"`
}

func (status *StationStatus) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	raw := new(struct {
		StationName     string  `xml:"locationName"`
		Latitude        float64 `xml:"lat"`
		Longitude       float64 `xml:"lon"`
		WeatherElements []struct {
			Name  string      `xml:"elementName"`
			Value interface{} `xml:"elementValue>value"`
		} `xml:"weatherElement"`
		Parameters []struct {
			Name  string      `xml:"parameterName"`
			Value interface{} `xml:"parameterValue"`
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
		log.Printf("Add elment %s => %v", element.Name, element)
	}
	// for _, parameter := range raw.Parameters {
	// 	switch parameter.Name {
	// 	case "CITY":
	// 		status.CityName = parameter.Value.(string)
	// 		break
	// 	case "CITY_SN":
	// 		status.CitySN = parameter.Value.(int)
	// 		break
	// 	case "TOWN":
	// 		status.TownName = parameter.Value.(string)
	// 		break
	// 	case "TOWN_SN":
	// 		status.TownSN = parameter.Value.(int)
	// 		break
	// 	}
	// }
	log.Printf("%v", raw)
	return nil
}

func (list *StationList) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	stations := new(struct {
		Locations []StationStatus `xml:"location"`
	})
	err := d.DecodeElement(stations, &start)
	if err != nil {
		return err
	}
	list.Locations = make(map[string]StationStatus, 150)
	for _, station := range stations.Locations {
		list.Locations[station.StationName] = station
	}
	return nil
}
