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
	DataId  string
}

func initDataSource(apiKey string) cwbDataSource {
	return cwbDataSource{APIKey: apiKey}
}

func (cwb cwbDataSource) loadDataSet(dataId string) (result cwbDataSet) {
	result = cwbDataSet{DataId: dataId}
	response, err := http.Get(fmt.Sprintf("%s?dataid=%s&authorizationKey=%s", ApiUrl, dataId, cwb.APIKey))
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

type weatherElement struct {
	ElementName  string `xml:"elementName"`
	ElementValue string `xml:"elementValue>value"`
}

type stationLocation struct {
	LocationName    string           `xml:"locationName"`
	WeatherElements []weatherElement `xml:"weatherElement"`
}

type xmlStationsStatus struct {
	Location stationLocation `xml:"location"`
}
