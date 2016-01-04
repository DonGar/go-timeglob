package timeglob

import (
	"time"
)

func intRange(begin, end int) []int {
	// Return a slice containing values between begin and end, inclusive.

	var length, dir int

	// Do we build up, or down?
	if begin <= end {
		length = end - begin + 1
		dir = 1
	} else {
		length = begin - end + 1
		dir = -1
	}

	result := make([]int, length)
	for i := 0; i < length; i++ {
		result[i] = begin + (i * dir)
	}
	return result
}

func reverseCopy(src []int) []int {
	if src == nil {
		return src
	}

	result := make([]int, len(src))

	output := len(result)
	for i := 0; i < len(src); i++ {
		output--
		result[output] = src[i]
	}

	return result
}

func (tg *TimeGlob) dateNoNormalize(year, month, day, hour, minute int) time.Time {

	// This is a wrapper around time.Date that ensures no values were normalized.
	// IE: Feb 30 doesn't become Mar 2.

	result := time.Date(
		year, time.Month(month), day,
		hour, minute, 0, 0,
		tg.location,
	)

	if result.Year() != year ||
		result.Month() != time.Month(month) ||
		result.Day() != day ||
		result.Hour() != hour ||
		result.Minute() != minute {
		return UNKNOWN
	}

	return result
}

func (tg *TimeGlob) addMinutesNoNormalize(base time.Time, minute int) time.Time {
	// Add minutes to an even hour, without normalizing.
	result := base.Add(time.Duration(minute) * time.Minute)

	if result.Hour() != base.Hour() || result.Minute() != minute {
		return UNKNOWN
	}

	return result
}
