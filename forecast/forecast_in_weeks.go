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
	f = new(WeeklyForecast)
	err = xml.Unmarshal(openData.DataSet, &f)
	if err != nil {
		f = nil
	}
	return
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
	Timeline []*Timed `xml:"time"`
}

type Timeline = []*Timed

func (te TimelineWeatherElement) GetTimeline() Timeline {
	return te.Timeline
}

type WeeklyForecast struct {
	Name            string
	Names           []string
	Geocode         string
	latitude        float32
	longitude       float32
	WeatherElements []TimelineWeatherElement
}
