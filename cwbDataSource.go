package twweather

import (
	"bytes"
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
	response, err := http.Get(fmt.Sprintf("%s?dataid=%s&authorizationKey=%s", ApiUrl, dataID, cwb.APIKey))
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

type rawWeatherElement struct {
	ElementName  string  `xml:"elementName"`
	ElementValue float64 `xml:"elementValue>value"`
}

type stationStatus struct {
	LocationName    string
	WeatherElements map[string]float64
}

type rawStationStatus struct {
	LocationName      string              `xml:"locationName"`
	RawWeatherElement []rawWeatherElement `xml:"weatherElement"`
}

type rawStationList struct {
	Locations []rawStationStatus `xml:"location"`
}

type stationList struct {
	Locations map[string]stationStatus
}

func (raw *rawStationList) Convert() *stationList {
	list := make(map[string]stationStatus, 11)
	for _, rawElem := range raw.Locations {
		list[rawElem.LocationName] = rawElem.Convert()
	}
	return &stationList{list}
}

func (status *rawStationStatus) Convert() (converted stationStatus) {
	converted.LocationName = status.LocationName
	converted.WeatherElements = status.ToMap()
	return
}

func (status rawStationStatus) ToMap() (elemMap map[string]float64) {
	elemMap = make(map[string]float64)
	for _, element := range status.RawWeatherElement {
		elemMap[element.ElementName] = element.ElementValue
	}
	return
}
