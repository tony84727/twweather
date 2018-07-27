package cwbdata

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	client *cwbDataClient
)

const ApiUrl = "http://opendata.cwb.gov.tw/opendataapi"
const CwbTimeFormat = time.RFC3339

func init() {
	client = newCWBDataClient()
}

func ParseTime(timeString string) (t time.Time, err error) {
	t, err = time.Parse(CwbTimeFormat, strings.TrimSpace(timeString))
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
	DataSet    []byte
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
		DataSet    struct {
			Data []byte `xml:",innerxml"`
		} `xml:"dataset"`
	})
	err := d.DecodeElement(original, &start)
	if err != nil {
		return err
	}

	openData.Identifier = strings.TrimSpace(original.Identifier)
	openData.Sender = strings.TrimSpace(original.Sender)
	openData.Status = strings.TrimSpace(original.Status)
	openData.Scope = strings.TrimSpace(original.Scope)
	openData.MsgType = strings.TrimSpace(original.MsgType)
	openData.DataID = strings.TrimSpace(original.DataID)
	openData.Source = strings.TrimSpace(original.Source)
	dataset := append([]byte("<dataset>"), original.DataSet.Data...)
	dataset = append(dataset, []byte("</dataset>")...)
	openData.DataSet = dataset
	err = AssignTime(original.Sent, &openData.Sent)
	if err != nil {
		return err
	}
	return nil
}

// GetOpenDataByData unmarshals data to CwbOpenData.
func GetOpenDataByData(data []byte) (openData *CwbOpenData, err error) {
	var d CwbOpenData
	openData = &d
	err = xml.Unmarshal(data, &openData)
	if err != nil {
		log.Println(string(data))

		return openData, err
	}
	return openData, nil
}

// GetOpenData makes API request to retrive data then pass it to GetOpenDataByData.
func GetOpenData(apiKey string, dataID string) (openData *CwbOpenData, err error) {
	client.SetAPIKey(apiKey)
	return client.GetOpenData(dataID)
}

func SetAPIKey(apiKey string) {
	client.SetAPIKey(apiKey)
}

type Cache interface {
	Get(dataID string) (data *CwbOpenData, exist bool)
	GetETag(dataID string) string
	Save(dataID, eTag string, data *CwbOpenData)
}

type cacheEntry struct {
	eTag string
	data *CwbOpenData
}

type inMemoryCache map[string]*cacheEntry

func (i inMemoryCache) Get(dataID string) (data *CwbOpenData, exist bool) {
	entry, exist := i[dataID]
	if exist {
		return entry.data, true
	}
	return nil, false
}

func (i inMemoryCache) GetETag(dataID string) string {
	entry, exist := i[dataID]
	if exist {
		return entry.eTag
	}
	return ""
}

func (i inMemoryCache) Save(dataID, eTag string, data *CwbOpenData) {
	i[dataID] = &cacheEntry{
		eTag,
		data,
	}
}

type cwbDataClient struct {
	httpClient *http.Client
	cache      Cache
	apiKey     string
}

func (c *cwbDataClient) SetAPIKey(apiKey string) {
	c.apiKey = apiKey
}

func newCWBDataClient() *cwbDataClient {
	return &cwbDataClient{
		httpClient: http.DefaultClient,
		cache:      make(inMemoryCache, 1),
	}
}

func (c *cwbDataClient) GetOpenData(dataID string) (openData *CwbOpenData, err error) {
	defer func() {
		if r := recover(); r != nil {
			openData = nil
			err = r.(error)
		}
	}()

	req := c.getRequest(dataID)
	eTag := c.cache.GetETag(dataID)
	if len(eTag) > 0 {
		req.Header.Add("If-None-Match", c.cache.GetETag(dataID))
	}
	response, err := c.httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	if response.StatusCode == 304 {
		cached, _ := c.cache.Get(dataID)
		return cached, nil
	}
	if response.StatusCode == 200 {
		buffer := new(bytes.Buffer)
		defer response.Body.Close()
		buffer.ReadFrom(response.Body)
		openData, err = GetOpenDataByData(buffer.Bytes())
		c.cache.Save(dataID, response.Header.Get("etag"), openData)
		return
	}
	return nil, fmt.Errorf("api returns status %s", response.Status)
}

func (c *cwbDataClient) getEndpoint(dataID string) string {
	return fmt.Sprintf(fmt.Sprintf("%s?dataid=%s&authorizationkey=%s", ApiUrl, dataID, c.apiKey))
}

func (c *cwbDataClient) getRequest(dataID string) *http.Request {
	req, err := http.NewRequest(http.MethodGet, c.getEndpoint(dataID), nil)
	if err != nil {
		panic(err)
	}
	return req
}
