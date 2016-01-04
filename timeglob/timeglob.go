package timeglob

import (
	"time"
)

// Used when no valid matching time exists.
var UNKNOWN time.Time

// nil values are used to represent wildcards.
type TimeGlob struct {
	year     []int
	month    []int
	day      []int
	hour     []int
	minute   []int
	location *time.Location
}

func new() TimeGlob {
	return TimeGlob{
		nil, nil, nil,
		[]int{0}, []int{0},
		time.Local}
}
