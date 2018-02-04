package forecast

import (
	"encoding/xml"
	"fmt"
	"log"
	"time"

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
		measurement := timed.(*Measurement)
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

func printTimedTimestamps(t Timed) {
	fmt.Println(t.Start().Format(time.RFC3339))
	fmt.Println(t.End().Format(time.RFC3339))
}

func printParameter(p Parameter) {
	fmt.Println(p.Name)
	fmt.Println(p.Value)
	fmt.Println(p.Unit)
}

func ExampleGetDescription() {
	te, err := getTimelineWeatherElementFrom("timeline_weather_element_description.xml")
	if err != nil {
		log.Panic(err.Error())
	}
	for _, timed := range te.GetTimeline() {
		desc := timed.(*Description)
		printTimedTimestamps(desc)
		for _, parameter := range desc.Parameters {
			printParameter(parameter)
		}
	}
	// Output:
	// 2018-02-03T18:00:00+08:00
	// 2018-02-04T06:00:00+08:00
	// 風向縮寫
	// NE
	// 16方位
}
