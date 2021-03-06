package timeglob

import (
	"gopkg.in/check.v1"
	"time"
)

func validateNext(c *check.C, tg *TimeGlob, now, expected time.Time) {
	c.Check(tg.Next(now), check.Equals, expected)
}

func validateNextSequence(c *check.C, tg *TimeGlob, start time.Time, expected []time.Time) {
	for _, exp := range expected {
		start = tg.Next(start)
		c.Check(start, check.Equals, exp)
	}
}

func (suite *MySuite) TestNextExplicit(c *check.C) {
	tg, err := Parse("2015/11/25 19:37 America/New_York")
	c.Assert(err, check.IsNil)

	matched := tg.dateNoNormalize(2015, 11, 25, 19, 37, 0)

	// Year
	now := tg.dateNoNormalize(2014, 1, 1, 0, 0, 0)
	validateNext(c, tg, now, matched)

	// Month
	now = tg.dateNoNormalize(2015, 10, 16, 20, 40, 0)
	validateNext(c, tg, now, matched)

	// Day
	now = tg.dateNoNormalize(2015, 11, 20, 20, 40, 0)
	validateNext(c, tg, now, matched)

	// We don't have to be much in the future for this to match.
	now = matched.Add(-1 * time.Nanosecond)
	validateNext(c, tg, now, matched)

	// Exact time.
	validateNext(c, tg, matched, UNKNOWN)

	// Later time.
	now = tg.dateNoNormalize(2015, 11, 25, 19, 38, 0)
	validateNext(c, tg, now, UNKNOWN)

	// Much later time.
	now = tg.dateNoNormalize(2017, 1, 1, 0, 0, 0)
	validateNext(c, tg, now, UNKNOWN)
}

func (suite *MySuite) TestNextDate(c *check.C) {
	tg, err := Parse("*/12/25 America/New_York")
	c.Assert(err, check.IsNil)

	expected2015 := tg.dateNoNormalize(2015, 12, 25, 0, 0, 0)
	expected2016 := tg.dateNoNormalize(2016, 12, 25, 0, 0, 0)

	now := tg.dateNoNormalize(2014, 12, 25, 0, 0, 0)
	validateNext(c, tg, now, expected2015)

	now = tg.dateNoNormalize(2015, 11, 25, 0, 0, 0)
	validateNext(c, tg, now, expected2015)

	now = tg.dateNoNormalize(2015, 12, 25, 19, 37, 0)
	validateNext(c, tg, now, expected2016)

	now = tg.dateNoNormalize(2015, 12, 26, 19, 37, 0)
	validateNext(c, tg, now, expected2016)

	now = tg.dateNoNormalize(2016, 1, 1, 0, 0, 0)
	validateNext(c, tg, now, expected2016)

	now = expected2015
	validateNext(c, tg, now, expected2016)
}

func (suite *MySuite) TestNextDateLeapDay(c *check.C) {
	tg, err := Parse("*/2/29 America/New_York")
	c.Assert(err, check.IsNil)

	expected2016 := tg.dateNoNormalize(2016, 2, 29, 0, 0, 0)
	expected2020 := tg.dateNoNormalize(2020, 2, 29, 0, 0, 0)

	now := tg.dateNoNormalize(2015, 1, 1, 0, 0, 0)
	validateNext(c, tg, now, expected2016)

	now = tg.dateNoNormalize(2015, 12, 25, 19, 37, 0)
	validateNext(c, tg, now, expected2016)

	now = tg.dateNoNormalize(2016, 2, 28, 4, 5, 0)
	validateNext(c, tg, now, expected2016)

	now = tg.dateNoNormalize(2016, 3, 1, 4, 5, 0)
	validateNext(c, tg, now, expected2020)
}

func (suite *MySuite) TestNextDay(c *check.C) {
	tg, err := Parse("*/*/29 America/New_York")
	c.Assert(err, check.IsNil)

	now := tg.dateNoNormalize(2015, 1, 1, 0, 0, 0)
	result := tg.Next(now)
	c.Check(result, check.Equals, tg.dateNoNormalize(2015, 1, 29, 0, 0, 0))

	now = tg.dateNoNormalize(2015, 2, 1, 19, 37, 0)
	result = tg.Next(now)
	c.Check(result, check.Equals, tg.dateNoNormalize(2015, 3, 29, 0, 0, 0))

	now = tg.dateNoNormalize(2016, 2, 1, 19, 37, 0)
	result = tg.Next(now)
	c.Check(result, check.Equals, tg.dateNoNormalize(2016, 2, 29, 0, 0, 0))

	now = tg.dateNoNormalize(2015, 12, 30, 4, 5, 0)
	result = tg.Next(now)
	c.Check(result, check.Equals, tg.dateNoNormalize(2016, 1, 29, 0, 0, 0))
}

func (suite *MySuite) TestNextFixedTime(c *check.C) {
	tg, err := Parse("13:12 America/New_York")
	c.Assert(err, check.IsNil)

	expectedA := tg.dateNoNormalize(2015, 1, 1, 13, 12, 0)
	expectedB := tg.dateNoNormalize(2015, 1, 2, 13, 12, 0)

	now := tg.dateNoNormalize(2014, 12, 31, 20, 0, 0)
	validateNext(c, tg, now, expectedA)

	now = tg.dateNoNormalize(2015, 1, 1, 0, 0, 0)
	validateNext(c, tg, now, expectedA)

	now = tg.dateNoNormalize(2015, 1, 1, 12, 40, 0)
	validateNext(c, tg, now, expectedA)

	now = tg.dateNoNormalize(2015, 1, 1, 13, 11, 0)
	validateNext(c, tg, now, expectedA)

	now = tg.dateNoNormalize(2015, 1, 1, 13, 12, 0)
	validateNext(c, tg, now, expectedB)

	now = tg.dateNoNormalize(2015, 1, 1, 13, 13, 0)
	validateNext(c, tg, now, expectedB)
}

func (suite *MySuite) TestNextDslStart(c *check.C) {
	tg, err := Parse("*:12 America/New_York")
	c.Assert(err, check.IsNil)

	// There is no 2 AM of that day, in that timezone.
	expectedA := tg.dateNoNormalize(2016, 3, 13, 1, 12, 0)
	expectedB := tg.dateNoNormalize(2016, 3, 13, 3, 12, 0)
	expectedC := tg.dateNoNormalize(2016, 3, 13, 4, 12, 0)

	now := tg.dateNoNormalize(2016, 3, 13, 1, 10, 0)
	result := tg.Next(now)
	c.Check(result, check.Equals, expectedA)

	now = now.Add(time.Hour)
	result = tg.Next(now)
	c.Check(result, check.Equals, expectedB)

	now = now.Add(time.Hour)
	result = tg.Next(now)
	c.Check(result, check.Equals, expectedC)
}

func (suite *MySuite) TestNextDslStop(c *check.C) {
	tg, err := Parse("*:12 America/New_York")
	c.Assert(err, check.IsNil)

	// 2 AM of that day, in that timezone.
	expectedA := tg.dateNoNormalize(2016, 11, 6, 0, 12, 0)
	expectedB := tg.dateNoNormalize(2016, 11, 6, 1, 12, 0)
	expectedC := expectedB.Add(time.Hour)
	expectedD := tg.dateNoNormalize(2016, 11, 6, 2, 12, 0)
	c.Assert(expectedB, check.Not(check.Equals), expectedC)
	c.Assert(expectedC, check.Not(check.Equals), expectedD)

	now := tg.dateNoNormalize(2016, 11, 6, 0, 10, 0)
	result := tg.Next(now)
	c.Check(result, check.Equals, expectedA)

	now = now.Add(time.Hour)
	result = tg.Next(now)
	c.Check(result, check.Equals, expectedB)

	now = now.Add(time.Hour)
	result = tg.Next(now)
	c.Check(result, check.Equals, expectedC)

	now = tg.dateNoNormalize(2016, 11, 6, 0, 50, 0)
	result = tg.Next(now)
	c.Check(result, check.Equals, expectedB)

	now = tg.dateNoNormalize(2016, 11, 6, 1, 50, 0)
	validateNext(c, tg, now, expectedC)
}

func (suite *MySuite) TestNextMinute(c *check.C) {
	tg, err := Parse("*:12 America/New_York")
	c.Assert(err, check.IsNil)

	expectedA := tg.dateNoNormalize(2015, 1, 1, 0, 12, 0)
	expectedB := tg.dateNoNormalize(2015, 1, 1, 13, 12, 0)
	expectedC := tg.dateNoNormalize(2015, 1, 1, 14, 12, 0)

	now := tg.dateNoNormalize(2014, 12, 31, 23, 13, 0)
	validateNext(c, tg, now, expectedA)

	now = tg.dateNoNormalize(2015, 1, 1, 0, 0, 0)
	validateNext(c, tg, now, expectedA)

	now = tg.dateNoNormalize(2015, 1, 1, 12, 40, 0)
	validateNext(c, tg, now, expectedB)

	now = tg.dateNoNormalize(2015, 1, 1, 13, 11, 0)
	validateNext(c, tg, now, expectedB)

	now = tg.dateNoNormalize(2015, 1, 1, 13, 12, 0)
	validateNext(c, tg, now, expectedC)

	now = tg.dateNoNormalize(2015, 1, 1, 13, 13, 0)
	validateNext(c, tg, now, expectedC)
}

func (suite *MySuite) TestNextWild(c *check.C) {
	tg, err := Parse("*:* America/New_York")
	c.Assert(err, check.IsNil)

	now := tg.dateNoNormalize(2014, 12, 31, 23, 59, 0)
	result := tg.Next(now)
	c.Check(result, check.Equals, tg.dateNoNormalize(2015, 1, 1, 0, 0, 0))

	now = tg.dateNoNormalize(2015, 1, 1, 0, 0, 0)
	result = tg.Next(now)
	c.Check(result, check.Equals, tg.dateNoNormalize(2015, 1, 1, 0, 1, 0))

	now = tg.dateNoNormalize(2015, 1, 1, 12, 59, 0)
	result = tg.Next(now)
	c.Check(result, check.Equals, tg.dateNoNormalize(2015, 1, 1, 13, 0, 0))
}

func (suite *MySuite) TestNextComma(c *check.C) {
	tg, err := Parse("2015,2016/2,3/10,20 2,14:0,30:15,45 UTC")
	c.Assert(err, check.IsNil)

	// Walk through the entire sequence.
	validateNextSequence(c, tg,
		tg.dateNoNormalize(2012, 3, 1, 0, 0, 0),
		[]time.Time{
			tg.dateNoNormalize(2015, 2, 10, 2, 0, 15),
			tg.dateNoNormalize(2015, 2, 10, 2, 0, 45),
			tg.dateNoNormalize(2015, 2, 10, 2, 30, 15),
			tg.dateNoNormalize(2015, 2, 10, 2, 30, 45),
			tg.dateNoNormalize(2015, 2, 10, 14, 0, 15),
			tg.dateNoNormalize(2015, 2, 10, 14, 0, 45),
			tg.dateNoNormalize(2015, 2, 10, 14, 30, 15),
			tg.dateNoNormalize(2015, 2, 10, 14, 30, 45),
			tg.dateNoNormalize(2015, 2, 20, 2, 0, 15),
			tg.dateNoNormalize(2015, 2, 20, 2, 0, 45),
			tg.dateNoNormalize(2015, 2, 20, 2, 30, 15),
			tg.dateNoNormalize(2015, 2, 20, 2, 30, 45),
			tg.dateNoNormalize(2015, 2, 20, 14, 0, 15),
			tg.dateNoNormalize(2015, 2, 20, 14, 0, 45),
			tg.dateNoNormalize(2015, 2, 20, 14, 30, 15),
			tg.dateNoNormalize(2015, 2, 20, 14, 30, 45),
			tg.dateNoNormalize(2015, 3, 10, 2, 0, 15),
			tg.dateNoNormalize(2015, 3, 10, 2, 0, 45),
			tg.dateNoNormalize(2015, 3, 10, 2, 30, 15),
			tg.dateNoNormalize(2015, 3, 10, 2, 30, 45),
			tg.dateNoNormalize(2015, 3, 10, 14, 0, 15),
			tg.dateNoNormalize(2015, 3, 10, 14, 0, 45),
			tg.dateNoNormalize(2015, 3, 10, 14, 30, 15),
			tg.dateNoNormalize(2015, 3, 10, 14, 30, 45),
			tg.dateNoNormalize(2015, 3, 20, 2, 0, 15),
			tg.dateNoNormalize(2015, 3, 20, 2, 0, 45),
			tg.dateNoNormalize(2015, 3, 20, 2, 30, 15),
			tg.dateNoNormalize(2015, 3, 20, 2, 30, 45),
			tg.dateNoNormalize(2015, 3, 20, 14, 0, 15),
			tg.dateNoNormalize(2015, 3, 20, 14, 0, 45),
			tg.dateNoNormalize(2015, 3, 20, 14, 30, 15),
			tg.dateNoNormalize(2015, 3, 20, 14, 30, 45),
			tg.dateNoNormalize(2016, 2, 10, 2, 0, 15),
			tg.dateNoNormalize(2016, 2, 10, 2, 0, 45),
			tg.dateNoNormalize(2016, 2, 10, 2, 30, 15),
			tg.dateNoNormalize(2016, 2, 10, 2, 30, 45),
			tg.dateNoNormalize(2016, 2, 10, 14, 0, 15),
			tg.dateNoNormalize(2016, 2, 10, 14, 0, 45),
			tg.dateNoNormalize(2016, 2, 10, 14, 30, 15),
			tg.dateNoNormalize(2016, 2, 10, 14, 30, 45),
			tg.dateNoNormalize(2016, 2, 20, 2, 0, 15),
			tg.dateNoNormalize(2016, 2, 20, 2, 0, 45),
			tg.dateNoNormalize(2016, 2, 20, 2, 30, 15),
			tg.dateNoNormalize(2016, 2, 20, 2, 30, 45),
			tg.dateNoNormalize(2016, 2, 20, 14, 0, 15),
			tg.dateNoNormalize(2016, 2, 20, 14, 0, 45),
			tg.dateNoNormalize(2016, 2, 20, 14, 30, 15),
			tg.dateNoNormalize(2016, 2, 20, 14, 30, 45),
			tg.dateNoNormalize(2016, 3, 10, 2, 0, 15),
			tg.dateNoNormalize(2016, 3, 10, 2, 0, 45),
			tg.dateNoNormalize(2016, 3, 10, 2, 30, 15),
			tg.dateNoNormalize(2016, 3, 10, 2, 30, 45),
			tg.dateNoNormalize(2016, 3, 10, 14, 0, 15),
			tg.dateNoNormalize(2016, 3, 10, 14, 0, 45),
			tg.dateNoNormalize(2016, 3, 10, 14, 30, 15),
			tg.dateNoNormalize(2016, 3, 10, 14, 30, 45),
			tg.dateNoNormalize(2016, 3, 20, 2, 0, 15),
			tg.dateNoNormalize(2016, 3, 20, 2, 0, 45),
			tg.dateNoNormalize(2016, 3, 20, 2, 30, 15),
			tg.dateNoNormalize(2016, 3, 20, 2, 30, 45),
			tg.dateNoNormalize(2016, 3, 20, 14, 0, 15),
			tg.dateNoNormalize(2016, 3, 20, 14, 0, 45),
			tg.dateNoNormalize(2016, 3, 20, 14, 30, 15),
			tg.dateNoNormalize(2016, 3, 20, 14, 30, 45),
			UNKNOWN,
		})
}

func (suite *MySuite) TestNextPerformance(c *check.C) {
	// Demonstrate performance issues.

	// This used to search through every second of the year for each step.
	tg, err := Parse("*:*:* UTC")
	c.Assert(err, check.IsNil)

	validateNextSequence(c, tg,
		tg.dateNoNormalize(2016, 12, 31, 23, 59, 57),
		[]time.Time{
			tg.dateNoNormalize(2016, 12, 31, 23, 59, 58),
			tg.dateNoNormalize(2016, 12, 31, 23, 59, 59),
			tg.dateNoNormalize(2017, 1, 1, 0, 0, 0),
		})
}
