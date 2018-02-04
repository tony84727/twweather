package cwbdata

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"
)

const ApiUrl = "http://opendata.cwb.gov.tw/opendataapi"
const CwbTimeFormat = time.RFC3339

func ParseTime(timeString string) (t time.Time, err error) {
	t, err = time.Parse(CwbTimeFormat, timeString)
	return
}

func AssignTime(timeString string, to *time.Time) error {
	t, err := ParseTime(timeString)
	if err != nil {
		return err
	}
	*to = t
	return nil
}

type OpenDataSource interface {
	GetOpenData(dateID string) (CwbOpenData, error)
}
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

func GetOpenData(apiKey string, dataID string) (openData CwbOpenData, err error) {
	response, err := http.Get(fmt.Sprintf("%s?dataid=%s&authorizationkey=%s", ApiUrl, dataID, apiKey))
	if err != nil {
		return
	}
	buffer := new(bytes.Buffer)
	defer response.Body.Close()
	buffer.ReadFrom(response.Body)
	openData = CwbOpenData{}
	err = xml.Unmarshal(buffer.Bytes(), &openData)
	return
}
