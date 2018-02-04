package cwbdata

import "testing"

func testParseTime(t *testing.T) {
	timeStr := "2018-02-03T18:00:00+08:00"
	time, _ := ParseTime(timeStr)
	if time.Year() != 2018 {
		t.Fail()
	}
}
