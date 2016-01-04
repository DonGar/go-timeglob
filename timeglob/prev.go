package timeglob

import (
	"time"
)

func (tg *TimeGlob) Prev(now time.Time) time.Time {
	// Find the closest time which matches the glob and is before, or equal to
	// now. Returns UNKNOWN if there isn't a match.

	result := tg.prevDate(now.In(tg.location))
	if result != UNKNOWN {
		result = result.In(now.Location())
	}

	return result
}

func (tg *TimeGlob) expandPrev(now time.Time) (years, months, days, hours, minutes []int) {
	// Expand wildcard values out to explict lists of values.

	years = reverseCopy(tg.year)
	if len(years) == 0 {
		// Expand years wildcard in a limited way to avoid searching forever.
		years = intRange(now.Year(), now.Year()-YEAR_SEARCH_DEPTH)
	}

	months = reverseCopy(tg.month)
	if len(months) == 0 {
		months = intRange(12, 1)
	}

	days = reverseCopy(tg.day)
	if len(days) == 0 {
		days = intRange(32, 1)
	}

	hours = reverseCopy(tg.hour)
	if len(hours) == 0 {
		hours = intRange(24, 0)
	}

	minutes = reverseCopy(tg.minute)
	if len(minutes) == 0 {
		minutes = intRange(61, 0)
	}

	return years, months, days, hours, minutes
}

func (tg *TimeGlob) prevDate(now time.Time) time.Time {
	years, months, days, hours, minutes := tg.expandPrev(now)

	for _, year := range years {
		for _, month := range months {
			for _, day := range days {
				for _, hour := range hours {

					// Cheesy, cheesy daylight savings hack.
					//
					// If we are using hour wildcards, and we can add an hour, but have
					// the same hour value (IE: 1 AM repeating), process minutes from the
					// extra hour in a special loop.
					if len(tg.hour) == 0 {
						base := tg.dateNoNormalize(year, month, day, hour, 0)
						advanced := base.Add(time.Hour)

						if base.Hour() == advanced.Hour() {
							for _, minute := range minutes {
								result := tg.addMinutesNoNormalize(advanced, minute)
								if result != UNKNOWN && (now.Equal(result) || now.After(result)) {
									return result
								}
							}
						}
					}

					for _, minute := range minutes {
						result := tg.dateNoNormalize(year, month, day, hour, minute)
						if result != UNKNOWN && (now.Equal(result) || now.After(result)) {
							return result
						}
					}
				}
			}
		}
	}

	return UNKNOWN
}
