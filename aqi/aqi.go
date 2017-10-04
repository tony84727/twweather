package aqi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type AQI struct {
	SiteName      string
	County        string
	AQI           uint16 `json:",string"`
	Pollutant     string
	Status        string
	NO            float32 `json:",string"`
	NO2           float32 `json:",string"`
	NOx           float32 `json:",string"`
	O3            float32 `json:",string"`
	O3_8hr        float32 `json:",string"`
	PM10          float32 `json:",string"`
	PM10Avg       float32 `json:"PM10_AVG"`
	PM25          float32 `json:"PM2.5"`
	PublishTime   string
	SO2           float32 `json:",string"`
	WindDirection uint16  `json:"WindDirec,string"`
	WindSpeed     float32 `json:",string"`
}

func (aqi *AQI) String() string {
	bytes, err := json.Marshal(aqi)
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes)
}

const APIUrl = "http://opendata.epa.gov.tw/ws/Data/AQI/?$format=json"

type Client struct {
	Cities *map[string]AQI
}

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "[AQI]", 0)
}

func New() (client *Client) {

	defer func() {
		if r := recover(); r != nil {
			log.Fatal(r)
		}
	}()

	response, err := http.Get(APIUrl)
	if err != nil {
		panic(fmt.Sprintf("Cannot fetch data from %s", APIUrl))
	} else {
		buffer := new(bytes.Buffer)
		defer response.Body.Close()
		length, err := buffer.ReadFrom(response.Body)
		if err != nil {
			panic(err)
		} else {
			logger.Printf("Read %d bytes from data source", length)
		}
		aqiData := make([]AQI, 10)
		json.Unmarshal(buffer.Bytes(), &aqiData)
		cities := make(map[string]AQI)
		for _, cityAQI := range aqiData {
			cities[cityAQI.County] = cityAQI
		}
		client = &Client{Cities: &cities}
		return
	}
}
