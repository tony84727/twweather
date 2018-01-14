package cwbdata

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"time"
)

const ApiUrl = "http://opendata.cwb.gov.tw/opendataapi"
const CwbTimeFormat = "2006-01-02T15:04:05-07:00"

type CwbDataSource struct {
	APIKey string
}

type CwbDataSet struct {
	RawData []byte
	DataID  string
}

type CwbOpenData struct {
	Identifier string    `xml:"identifier"`
	Sender     string    `xml:"sender"`
	Sent       time.Time `xml:"sent"`
	Status     string    `xml:"status"`
	Scope      string    `xml:"scope"`
	MsgType    string    `xml:"msgType"`
	DataID     string    `xml:"dataid"`
	Source     string    `xml:"source"`
	DataSet    []byte    `xml:"dataset,innerXML"`
}
type rawCwbOpenData struct {
	CwbOpenData
	Sent string `xml:"sent"`
}

func (openData *CwbOpenData) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	origin := new(rawCwbOpenData)
	d.DecodeElement(origin, &start)

	openData.Identifier = origin.Identifier
	openData.Sender = origin.Sender
	openData.Status = origin.Status
	openData.Scope = origin.Scope
	openData.MsgType = origin.MsgType
	openData.DataID = origin.DataID
	openData.Source = origin.Source
	openData.DataSet = origin.DataSet
	timeStamp, err := time.Parse(CwbTimeFormat, origin.Sent)
	if err != nil {
		return err
	}
	openData.Sent = timeStamp
	return nil
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
