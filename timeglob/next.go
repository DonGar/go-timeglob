package timeglob

import (
	"time"
)

const YEAR_SEARCH_DEPTH = 5

func (tg *TimeGlob) Next(now time.Time) time.Time {
	// Find the closest time which matches the glob and is after now. Returns
	// UNKNOWN if there isn't a match.

	result := tg.nextDate(now.In(tg.location))
	if result != UNKNOWN {
		result = result.In(now.Location())
	}

	return result
}

func (tg *TimeGlob) expandNext(now time.Time) (years, months, days, hours, minutes, seconds []int) {
	// Expand wildcard values out to explict lists of values.

	years = tg.year
	if len(years) == 0 {
		// Expand years wildcard in a limited way to avoid searching forever.
		years = intRange(now.Year(), now.Year()+YEAR_SEARCH_DEPTH)
	}

	months = tg.month
	if len(months) == 0 {
		months = intRange(1, 12)
	}

	days = tg.day
	if len(days) == 0 {
		days = intRange(1, 32)
	}

	hours = tg.hour
	if len(hours) == 0 {
		hours = intRange(0, 24)
	}

	minutes = tg.minute
	if len(minutes) == 0 {
		minutes = intRange(0, 61)
	}

	seconds = tg.second
	if len(seconds) == 0 {
		seconds = intRange(0, 61)
	}

	return years, months, days, hours, minutes, seconds
}

func (tg *TimeGlob) nextDate(now time.Time) time.Time {
	years, months, days, hours, minutes, seconds := tg.expandNext(now)

	dateNow := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, tg.location)

	for _, year := range years {
		for _, month := range months {
			for _, day := range days {

				// For performance validate that each date might be parse of a valid
				// result before searching inside the date.
				searchDate := tg.dateNoNormalize(year, month, day, 0, 0, 0)
				if searchDate == UNKNOWN || searchDate.Before(dateNow) {
					continue
				}

				for _, hour := range hours {
					for _, minute := range minutes {
						for _, second := range seconds {
							result := tg.dateNoNormalize(year, month, day, hour, minute, second)
							if result != UNKNOWN && now.Before(result) {
								return result
							}
						}
					}

					// Cheesy, cheesy daylight savings hack.
					//
					// If we are using hour wildcards, and we can add an hour, but have
					// the same hour value (IE: 1 AM repeating), process minutes from the
					// extra hour in a special loop.
					if len(tg.hour) == 0 {
						base := tg.dateNoNormalize(year, month, day, hour, 0, 0)
						advanced := base.Add(time.Hour)

						if base.Hour() == advanced.Hour() {
							for _, minute := range minutes {
								for _, second := range seconds {
									result := tg.adjustMinutesSeconds(advanced, minute, second)
									if result != UNKNOWN && now.Before(result) {
										return result
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return UNKNOWN
}
