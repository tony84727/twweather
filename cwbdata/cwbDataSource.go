package cwbdata

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

const ApiUrl = "http://opendata.cwb.gov.tw/opendataapi"

type CwbDataSource struct {
	APIKey string
}

type CwbDataSet struct {
	RawData []byte
	DataID  string
}

func InitDataSource(apiKey string) CwbDataSource {
	return CwbDataSource{APIKey: apiKey}
}

func (cwb CwbDataSource) LoadDataSet(dataID string) (result CwbDataSet) {
	result = CwbDataSet{DataID: dataID}
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
