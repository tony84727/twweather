package twweather

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
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
