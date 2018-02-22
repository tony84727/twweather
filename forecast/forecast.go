package forecast

import (
	"encoding/xml"
	"strconv"
	"time"

	"github.com/MinecraftXwinP/twweather/cwbdata"
)

func calcWeeklyForecastDatasetID(base int, city int) string {
	return strconv.Itoa(city*4 + base)
}

func getTwoWeekDatasetID(city int) string {
	return "F-D0047-0" + calcWeeklyForecastDatasetID(1, city)
}

func getOneWeekDatasetID(city int) string {
	return "F-D0047-0" + calcWeeklyForecastDatasetID(3, city)
}

func GetWeeklyForecast(apiKey string, city int) (f *WeeklyForecast, err error) {
	openData, err := cwbdata.GetOpenData(apiKey, getOneWeekDatasetID(city))
	if err != nil {
		return
	}
	f, err = OpenDataToWeeklyForecast(openData)
	return
}

func OpenDataToWeeklyForecast(opendata cwbdata.CwbOpenData) (*WeeklyForecast, error) {
	wf := new(WeeklyForecast)
	err := xml.Unmarshal(opendata.DataSet, wf)
	if err != nil {
		return nil, err
	}
	return wf, nil
}

type Forecast struct {
	Description        string
	Language           string
	IssueTime          time.Time
	UpdateTime         time.Time
	ContentDescription string
	LocationName       string
}

func (f *Forecast) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	original := new(struct {
		Description        string `xml:"datasetInfo>datasetDescription"`
		Language           string `xml:"datasetInfo>datasetLanguage"`
		IssueTime          string `xml:"datasetInfo>issueTime"`
		UpdateTime         string `xml:"datasetInfo>update"`
		ContentDescription string `xml:"contents>contentDescription"`
		LocationName       string `xml:"locations>locationsName"`
	})
	err := d.DecodeElement(original, &start)
	if err != nil {
		return err
	}
	f.Description = original.Description
	f.Language = original.Language
	err = cwbdata.AssignTime(original.IssueTime, &f.IssueTime)
	if err != nil {
		return err
	}
	err = cwbdata.AssignTime(original.UpdateTime, &f.UpdateTime)
	f.ContentDescription = original.ContentDescription
	f.LocationName = original.LocationName
	return nil
}

type WeeklyForecast struct {
	Forecast
	Locations []LocationForecast
}

func (wf *WeeklyForecast) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	original := new(struct {
		Description        string             `xml:"datasetInfo>datasetDescription"`
		Language           string             `xml:"datasetInfo>datasetLanguage"`
		IssueTime          string             `xml:"datasetInfo>issueTime"`
		UpdateTime         string             `xml:"datasetInfo>update"`
		ContentDescription string             `xml:"contents>contentDescription"`
		LocationName       string             `xml:"locations>locationsName"`
		Locations          []LocationForecast `xml:"locations>location"`
	})
	err := d.DecodeElement(original, &start)
	if err != nil {
		return err
	}
	wf.Description = original.Description
	wf.Language = original.Language
	err = cwbdata.AssignTime(original.IssueTime, &wf.IssueTime)
	if err != nil {
		return err
	}
	err = cwbdata.AssignTime(original.UpdateTime, &wf.UpdateTime)
	if err != nil {
		return err
	}
	wf.ContentDescription = original.ContentDescription
	wf.Locations = original.Locations
	return nil
}

type Location struct {
	Name      string  `xml:"locationName"`
	Geocode   string  `xml:"geocode"`
	Latitude  float32 `xml:"lat"`
	Longitude float32 `xml:"lon"`
}

type LocationForecast struct {
	Location
	WeatherElements []TimelineWeatherElement `xml:"weatherElement"`
}
type Timed struct {
	Start time.Time
	End   time.Time
	Data  []interface{}
}

func (t *Timed) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var data interface{}
	for {
		tok, err := d.Token()
		if err != nil {
			return err
		}
		switch v := tok.(type) {
		case xml.StartElement:
			switch v.Name.Local {
			case "elementValue":
				data = &Measurement{}
			case "parameter":
				data = &Parameter{}
			case "startTime", "endTime":
				timeStr := new(string)
				err := d.DecodeElement(timeStr, &v)
				if err != nil {
					return err
				}
				if v.Name.Local == "startTime" {
					cwbdata.AssignTime(*timeStr, &t.Start)
					continue
				}
				cwbdata.AssignTime(*timeStr, &t.End)
				continue
			}
			if data != nil {
				err = d.DecodeElement(data, &v)
				if err != nil {
					return err
				}
				t.Data = append(t.Data, data)
				data = nil
			}
			break
		case xml.EndElement:
			if v == start.End() {
				return nil
			}
		}
	}
}

type Timeline = []*Timed
type Measurement struct {
	Value string `xml:"value"`
	Unit  string `xml:"measures"`
}

type Parameter struct {
	Name  string `xml:"parameterName"`
	Value string `xml:"parameterValue"`
	Unit  string `xml:"parameterUnit"`
}

type TimelineWeatherElement struct {
	Name     string   `xml:"elementName"`
	Timeline Timeline `xml:"time"`
}
