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

type Timed interface {
	Start() time.Time
	End() time.Time
}

type timed struct {
	start time.Time `xml:"startTime"`
	end   time.Time `xml:"endTime"`
}

func (t timed) Start() time.Time {
	return t.start
}

func (t timed) End() time.Time {
	return t.end
}

type stringTimed struct {
	Start string `xml:"startTime"`
	End   string `xml:"endTime"`
}

func (t *timed) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	original := new(struct {
		stringTimed
	})
	d.DecodeElement(original, &start)
	startTime, err := time.Parse(cwbdata.CwbTimeFormat, original.Start)
	if err != nil {
		return err
	}
	t.start = startTime

	endTime, err := time.Parse(cwbdata.CwbTimeFormat, original.End)
	if err != nil {
		return err
	}
	t.end = endTime
	return nil
}

func isDescription(start xml.StartElement) bool {
	for _, attr := range start.Attr {
		if attr.Name.Local == "parameter" {
			return true
		}
	}
	return false
}

type Measurement struct {
	timed
	Value string `xml:"elementValue>value"`
	Unit  string `xml:"elementValue>measures"`
}

func (m *Measurement) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	original := new(struct {
		stringTimed
		Value string `xml:"elementValue>value"`
		Unit  string `xml:"elementValue>measures"`
	})
	err := d.DecodeElement(original, &start)
	if err != nil {
		return err
	}
	err = cwbdata.AssignTime(original.Start, &m.start)
	if err != nil {
		return err
	}
	err = cwbdata.AssignTime(original.End, &m.end)
	if err != nil {
		return err
	}
	m.Unit = original.Unit
	m.Value = original.Value
	return nil
}

type Parameter struct {
	Name  string `xml:"parameterName"`
	Value string `xml:"parameterValue"`
	Unit  string `xml:"parameterUnit"`
}

type Description struct {
	timed
	Parameters []Parameter
}

func (desc *Description) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	original := new(struct {
		stringTimed
		Parameters []Parameter `xml:"parameter"`
	})
	err := d.DecodeElement(original, &start)
	if err != nil {
		return err
	}
	err = cwbdata.AssignTime(original.Start, &desc.start)
	if err != nil {
		return err
	}
	err = cwbdata.AssignTime(original.End, &desc.end)
	if err != nil {
		return err
	}
	desc.Parameters = original.Parameters
	return nil
}

type Timeline []Timed
type TimelineWeatherElement struct {
	name     string
	timeline Timeline
}

func (te TimelineWeatherElement) GetTimeline() Timeline {
	return te.timeline
}

type timelinePart struct {
	Data Timed
}

func (tp *timelinePart) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// try unmarshal measurement
	measurement := &Measurement{}
	err := d.DecodeElement(measurement, &start)
	if err == nil && len(measurement.Unit) != 0 {
		tp.Data = measurement
		return nil
	}
	description := &Description{}
	err = d.DecodeElement(description, &start)
	if err != nil {
		return err
	}
	tp.Data = description
	return nil
}

func (te *TimelineWeatherElement) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	original := new(struct {
		Name string         `xml:"elementName"`
		Time []timelinePart `xml:"time"`
	})
	err := d.DecodeElement(original, &start)
	if err != nil {
		return err
	}
	timeline := make([]Timed, 0, 1)
	for _, tp := range original.Time {
		timeline = append(timeline, tp.Data)
	}
	te.timeline = timeline
	return nil
}

type WeeklyForecast struct {
	Name            string
	Names           []string
	Geocode         string
	latitude        float32
	longitude       float32
	WeatherElements []TimelineWeatherElement
}
