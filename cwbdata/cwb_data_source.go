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

func (openData *CwbOpenData) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	original := new(struct {
		Identifier string `xml:"identifier"`
		Sender     string `xml:"sender"`
		Sent       string `xml:"sent"`
		Status     string `xml:"status"`
		Scope      string `xml:"scope"`
		MsgType    string `xml:"msgType"`
		DataID     string `xml:"dataid"`
		Source     string `xml:"source"`
		DataSet    []byte `xml:"dataset,innerXML"`
	})
	err := d.DecodeElement(original, &start)
	if err != nil {
		return err
	}

	openData.Identifier = original.Identifier
	openData.Sender = original.Sender
	openData.Status = original.Status
	openData.Scope = original.Scope
	openData.MsgType = original.MsgType
	openData.DataID = original.DataID
	openData.Source = original.Source
	openData.DataSet = original.DataSet
	err = AssignTime(original.Sent, &openData.Sent)
	if err != nil {
		return err
	}
	return nil
}

// GetOpenDataByData unmarshals data to CwbOpenData.
func GetOpenDataByData(data []byte) (openData CwbOpenData, err error) {
	err = xml.Unmarshal(data, &openData)
	return
}

// GetOpenData makes API request to retrive data then pass it to GetOpenDataByData.
func GetOpenData(apiKey string, dataID string) (openData CwbOpenData, err error) {
	response, err := http.Get(fmt.Sprintf("%s?dataid=%s&authorizationkey=%s", ApiUrl, dataID, apiKey))
	if err != nil {
		return
	}
	buffer := new(bytes.Buffer)
	defer response.Body.Close()
	buffer.ReadFrom(response.Body)
	openData, err = GetOpenDataByData(buffer.Bytes())
	return
}
