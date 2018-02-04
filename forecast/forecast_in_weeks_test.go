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
	data := testutil.Load("timeline_weather_element_measurement.xml")
	log.Println(string(data))
	te := &TimelineWeatherElement{}
	err := xml.Unmarshal(data, te)
	if err != nil {
		panic(err)
	}
	for _, timed := range te.Timeline {
		measurement := timed.(*Measurement)
		fmt.Println(measurement.Start.Format(time.RFC3339))
		fmt.Println(measurement.End.Format(time.RFC3339))
	}
	// Output:
	// 2018-02-03T18:00:00+08:00
	// 2018-02-04T06:00:00+08:00
}
