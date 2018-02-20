package cwbdata

import (
	"fmt"
	"log"
	"testing"

	"github.com/MinecraftXwinP/twweather/testutil"
)

func testParseTime(t *testing.T) {
	timeStr := "2018-02-03T18:00:00+08:00"
	time, _ := ParseTime(timeStr)
	if time.Year() != 2018 {
		t.Fail()
	}
}

func ExampleGetOpenDataByData() {
	data := testutil.Load("fake_dataset.xml")
	opendata, err := GetOpenDataByData(data)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(opendata.DataID)
	fmt.Println(opendata.Identifier)
	fmt.Println(opendata.MsgType)
	fmt.Println(opendata.Scope)
	fmt.Println(opendata.Sender)
	fmt.Println(opendata.Sent)
	fmt.Println(opendata.Source)
	fmt.Println(opendata.Status)
	fmt.Println(string(opendata.DataSet))
	// Output:
	// D0047-003
	// 12ef1673-921a-bfdf-39b5-6f57f1d61a5a
	// Issue
	// Public
	// weather@cwb.gov.tw
	// 2018-01-13 17:06:00 +0800 CST
	// MFC
	// Actual
	// <dataset><a><b>ab</b><c>ac</c></a></dataset>
}
