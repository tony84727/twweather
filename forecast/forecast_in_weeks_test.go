package forecast

import (
	"fmt"
)

func ExampleGetDatasetID() {
	fmt.Println(getOneWeekDatasetID(HSINCHU_CITY))
	fmt.Println(getTwoWeekDatasetID(HSINCHU_CITY))
	// Output:
	// F-D0047-055
	// F-D0047-053
}
