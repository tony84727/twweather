package forecast

import (
	"encoding/xml"
	"fmt"
	"log"
	"time"

	"github.com/MinecraftXwinP/twweather/cwbdata"

	"github.com/MinecraftXwinP/twweather/testutil"
)

func ExampleGetDatasetID() {
	fmt.Println(getOneWeekDatasetID(HsinchuCity))
	fmt.Println(getTwoWeekDatasetID(HsinchuCity))
	// Output:
	// F-D0047-055
	// F-D0047-053
}

func ExampleGetMeasurement() {
	te, err := getTimelineWeatherElementFrom("timeline_weather_element_measurement.xml")
	if err != nil {
		log.Panic(err.Error())
	}
	for _, timed := range te.GetTimeline() {
		printTimedTimestamps(timed)
		measurement := timed.Data[0].(*Measurement)
		fmt.Println(measurement.Value)
		fmt.Println(measurement.Unit)
	}
	// Output:
	// 2018-02-03T18:00:00+08:00
	// 2018-02-04T06:00:00+08:00
	// 9
	// C
}
func getTimelineWeatherElementFrom(path string) (*TimelineWeatherElement, error) {
	data := testutil.Load(path)
	te := &TimelineWeatherElement{}
	err := xml.Unmarshal(data, te)
	if err != nil {
		return nil, err
	}
	return te, nil
}

func printTimedTimestamps(t *Timed) {
	fmt.Println(t.Start.Format(time.RFC3339))
	fmt.Println(t.End.Format(time.RFC3339))
}

func printParameter(p Parameter) {
	fmt.Println(p.Name)
	fmt.Println(p.Value)
	fmt.Println(p.Unit)
}

func ExampleGetParameter() {
	te, err := getTimelineWeatherElementFrom("timeline_weather_element_description.xml")
	if err != nil {
		log.Panic(err.Error())
	}
	for _, timed := range te.GetTimeline() {
		p := timed.Data[0].(*Parameter)
		printTimedTimestamps(timed)
		fmt.Println(p.Name)
		fmt.Println(p.Unit)
		fmt.Println(p.Value)
	}
	// Output:
	// 2018-02-03T18:00:00+08:00
	// 2018-02-04T06:00:00+08:00
	// 風向縮寫
	// 16方位
	// NE
}

func ExampleOpenDataToWeeklyForecast() {
	data := testutil.Load("weekly_forecast_sample.xml")
	opendata, err := cwbdata.GetOpenDataByData(data)
	if err != nil {
		log.Panic(err)
	}
	wf, err := OpenDataToWeeklyForecast(opendata)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(wf.Description)
	fmt.Println(wf.Language)
	fmt.Println(wf.IssueTime)
	fmt.Println(wf.UpdateTime)
	// Output:
	// 宜蘭縣未來1週天氣預報
	// zh-TW
	// 2018-01-13 17:00:00 +0800 CST
	// 2018-01-13 17:06:00 +0800 CST
}
