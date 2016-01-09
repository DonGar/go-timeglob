package timeglob

import (
	"gopkg.in/check.v1"
	"time"
)

func validatePrev(c *check.C, tg *TimeGlob, now, expected time.Time) {
	c.Check(tg.Prev(now), check.Equals, expected)
}

func validatePrevSequence(c *check.C, tg *TimeGlob, start time.Time, expected []time.Time) {
	for _, exp := range expected {
		start = tg.Prev(start)
		c.Check(start, check.Equals, exp)
		// Prev uses <= logic so would enver step back without help.
		start = start.Add(-1 * time.Nanosecond)
	}
}

func (suite *MySuite) TestPrevExplicit(c *check.C) {
	tg, err := Parse("2015/11/25 19:37 America/New_York")
	c.Assert(err, check.IsNil)

	matched := tg.dateNoNormalize(2015, 11, 25, 19, 37)

	// Year
	now := tg.dateNoNormalize(2016, 1, 1, 0, 0)
	validatePrev(c, tg, now, matched)

	// Month
	now = tg.dateNoNormalize(2015, 12, 16, 20, 40)
	validatePrev(c, tg, now, matched)

	// Day
	now = tg.dateNoNormalize(2015, 11, 30, 20, 40)
	validatePrev(c, tg, now, matched)

	// We don't have to be much in the future for this to match.
	now = matched.Add(time.Nanosecond)
	validatePrev(c, tg, now, matched)

	// Exact time.
	now = matched
	validatePrev(c, tg, now, matched)

	// Earlier time.
	now = tg.dateNoNormalize(2015, 11, 25, 19, 36)
	validatePrev(c, tg, now, UNKNOWN)

	// Much earlier time.
	now = tg.dateNoNormalize(2014, 1, 1, 0, 0)
	validatePrev(c, tg, now, UNKNOWN)
}

func (suite *MySuite) TestPrevDate(c *check.C) {
	tg, err := Parse("*/12/25 America/New_York")
	c.Assert(err, check.IsNil)

	expected2015 := tg.dateNoNormalize(2015, 12, 25, 0, 0)
	expected2016 := tg.dateNoNormalize(2016, 12, 25, 0, 0)

	// Year change.
	now := tg.dateNoNormalize(2017, 1, 1, 0, 0)
	result := tg.Prev(now)
	c.Check(result, check.Equals, expected2016)

	// Day change.
	now = tg.dateNoNormalize(2016, 12, 30, 20, 40)
	result = tg.Prev(now)
	c.Check(result, check.Equals, expected2016)

	// Time change.
	now = tg.dateNoNormalize(2016, 12, 25, 19, 37)
	result = tg.Prev(now)
	c.Check(result, check.Equals, expected2016)

	// Slightly earlier time.
	now = tg.dateNoNormalize(2016, 12, 24, 19, 37)
	result = tg.Prev(now)
	c.Check(result, check.Equals, expected2015)

	// Exact match.
	now = expected2016
	result = tg.Prev(now)
	c.Check(result, check.Equals, expected2016)
}

func (suite *MySuite) TestPrevDateLeapDay(c *check.C) {
	tg, err := Parse("*/2/29 America/New_York")
	c.Assert(err, check.IsNil)

	expected2012 := tg.dateNoNormalize(2012, 2, 29, 0, 0)
	expected2016 := tg.dateNoNormalize(2016, 2, 29, 0, 0)

	// Year change.
	now := tg.dateNoNormalize(2017, 1, 1, 0, 0)
	result := tg.Prev(now)
	c.Check(result, check.Equals, expected2016)

	// Month change.
	now = tg.dateNoNormalize(2016, 12, 25, 19, 37)
	result = tg.Prev(now)
	c.Check(result, check.Equals, expected2016)

	// Same day.
	now = tg.dateNoNormalize(2016, 2, 29, 4, 5)
	result = tg.Prev(now)
	c.Check(result, check.Equals, expected2016)

	// 4 years back.
	now = tg.dateNoNormalize(2016, 2, 28, 4, 5)
	result = tg.Prev(now)
	c.Check(result, check.Equals, expected2012)
}

func (suite *MySuite) TestPrevDay(c *check.C) {
	tg, err := Parse("*/*/29 America/New_York")
	c.Assert(err, check.IsNil)

	// Previous year.
	now := tg.dateNoNormalize(2015, 1, 1, 0, 0)
	result := tg.Prev(now)
	c.Check(result, check.Equals, tg.dateNoNormalize(2014, 12, 29, 0, 0))

	// Previous month.
	now = tg.dateNoNormalize(2015, 2, 1, 19, 37)
	result = tg.Prev(now)
	c.Check(result, check.Equals, tg.dateNoNormalize(2015, 1, 29, 0, 0))

	// Earlier day.
	now = tg.dateNoNormalize(2015, 1, 30, 4, 5)
	result = tg.Prev(now)
	c.Check(result, check.Equals, tg.dateNoNormalize(2015, 1, 29, 0, 0))

	// Same day, later time.
	now = tg.dateNoNormalize(2015, 1, 29, 19, 37)
	result = tg.Prev(now)
	c.Check(result, check.Equals, tg.dateNoNormalize(2015, 1, 29, 0, 0))
}

func (suite *MySuite) TestPrevFixedTime(c *check.C) {
	tg, err := Parse("13:12 America/New_York")
	c.Assert(err, check.IsNil)

	// Previous year.
	now := tg.dateNoNormalize(2015, 1, 1, 0, 0)
	result := tg.Prev(now)
	c.Check(result, check.Equals, tg.dateNoNormalize(2014, 12, 31, 13, 12))

	// Previous month.
	now = tg.dateNoNormalize(2015, 2, 1, 12, 3)
	result = tg.Prev(now)
	c.Check(result, check.Equals, tg.dateNoNormalize(2015, 1, 31, 13, 12))

	// Earlier day.
	now = tg.dateNoNormalize(2015, 1, 30, 4, 5)
	result = tg.Prev(now)
	c.Check(result, check.Equals, tg.dateNoNormalize(2015, 1, 29, 13, 12))

	// Same day, later time.
	now = tg.dateNoNormalize(2015, 1, 29, 19, 37)
	result = tg.Prev(now)
	c.Check(result, check.Equals, tg.dateNoNormalize(2015, 1, 29, 13, 12))
}

func (suite *MySuite) TestPrevDslStart(c *check.C) {
	tg, err := Parse("*:12 America/New_York")
	c.Assert(err, check.IsNil)

	// There is no 2 AM of this day, in that timezone.
	expectedA := tg.dateNoNormalize(2016, 3, 13, 3, 12)
	expectedB := tg.dateNoNormalize(2016, 3, 13, 1, 12)
	expectedC := tg.dateNoNormalize(2016, 3, 13, 0, 12)

	now := tg.dateNoNormalize(2016, 3, 13, 4, 10)
	result := tg.Prev(now)
	c.Check(result, check.Equals, expectedA)

	now = now.Add(-1 * time.Hour)
	result = tg.Prev(now)
	c.Check(result, check.Equals, expectedB)

	now = now.Add(-1 * time.Hour)
	result = tg.Prev(now)
	c.Check(result, check.Equals, expectedC)
}

func (suite *MySuite) TestPrevDslStop(c *check.C) {
	tg, err := Parse("*:12 America/New_York")
	c.Assert(err, check.IsNil)

	// 2 AM repeats on this day, in that timezone.
	expectedA := tg.dateNoNormalize(2016, 11, 6, 3, 12)
	expectedB := expectedA.Add(-1 * time.Hour)
	expectedC := expectedB.Add(-1 * time.Hour)
	expectedD := expectedC.Add(-1 * time.Hour)
	c.Assert(expectedB, check.Not(check.Equals), expectedC)
	c.Assert(expectedC, check.Not(check.Equals), expectedD)
	c.Assert(expectedD, check.Equals, tg.dateNoNormalize(2016, 11, 6, 1, 12))

	now := tg.dateNoNormalize(2016, 11, 6, 3, 14)
	result := tg.Prev(now)
	c.Check(result, check.Equals, expectedA)

	now = now.Add(-1 * time.Hour)
	result = tg.Prev(now)
	c.Check(result, check.Equals, expectedB)

	now = now.Add(-1 * time.Hour)
	result = tg.Prev(now)
	c.Check(result, check.Equals, expectedC)

	now = now.Add(-1 * time.Hour)
	result = tg.Prev(now)
	c.Check(result, check.Equals, expectedD)
}

func (suite *MySuite) TestPrevMinute(c *check.C) {
	tg, err := Parse("*:12 America/New_York")
	c.Assert(err, check.IsNil)

	expectedA := tg.dateNoNormalize(2015, 12, 31, 23, 12)
	expectedB := tg.dateNoNormalize(2015, 12, 31, 22, 12)
	expectedC := tg.dateNoNormalize(2015, 12, 31, 21, 12)

	now := tg.dateNoNormalize(2016, 1, 1, 0, 11)
	result := tg.Prev(now)
	c.Check(result, check.Equals, expectedA)

	now = expectedA.Add(-1 * time.Second)
	result = tg.Prev(now)
	c.Check(result, check.Equals, expectedB)

	now = expectedB.Add(-1 * time.Second)
	result = tg.Prev(now)
	c.Check(result, check.Equals, expectedC)
}

func (suite *MySuite) TestPrevWild(c *check.C) {
	tg, err := Parse("*:* America/New_York")
	c.Assert(err, check.IsNil)

	now := tg.dateNoNormalize(2015, 1, 1, 0, 0)
	result := tg.Prev(now)
	c.Check(result, check.Equals, now)

	now = tg.dateNoNormalize(2014, 12, 31, 23, 59)
	result = tg.Prev(now)
	c.Check(result, check.Equals, now)

	now = tg.dateNoNormalize(2014, 12, 31, 0, 0)
	result = tg.Prev(now)
	c.Check(result, check.Equals, now)
}

func (suite *MySuite) TestPrevComma(c *check.C) {
	tg, err := Parse("2015,2016/2,3/10,20 2,14:0,30 UTC")
	c.Assert(err, check.IsNil)

	// Walk through the entire sequence.
	validatePrevSequence(c, tg,
		tg.dateNoNormalize(2017, 3, 1, 0, 0),
		[]time.Time{
			tg.dateNoNormalize(2016, 3, 20, 14, 30),
			tg.dateNoNormalize(2016, 3, 20, 14, 0),
			tg.dateNoNormalize(2016, 3, 20, 2, 30),
			tg.dateNoNormalize(2016, 3, 20, 2, 0),
			tg.dateNoNormalize(2016, 3, 10, 14, 30),
			tg.dateNoNormalize(2016, 3, 10, 14, 0),
			tg.dateNoNormalize(2016, 3, 10, 2, 30),
			tg.dateNoNormalize(2016, 3, 10, 2, 0),
			tg.dateNoNormalize(2016, 2, 20, 14, 30),
			tg.dateNoNormalize(2016, 2, 20, 14, 0),
			tg.dateNoNormalize(2016, 2, 20, 2, 30),
			tg.dateNoNormalize(2016, 2, 20, 2, 0),
			tg.dateNoNormalize(2016, 2, 10, 14, 30),
			tg.dateNoNormalize(2016, 2, 10, 14, 0),
			tg.dateNoNormalize(2016, 2, 10, 2, 30),
			tg.dateNoNormalize(2016, 2, 10, 2, 0),
			tg.dateNoNormalize(2015, 3, 20, 14, 30),
			tg.dateNoNormalize(2015, 3, 20, 14, 0),
			tg.dateNoNormalize(2015, 3, 20, 2, 30),
			tg.dateNoNormalize(2015, 3, 20, 2, 0),
			tg.dateNoNormalize(2015, 3, 10, 14, 30),
			tg.dateNoNormalize(2015, 3, 10, 14, 0),
			tg.dateNoNormalize(2015, 3, 10, 2, 30),
			tg.dateNoNormalize(2015, 3, 10, 2, 0),
			tg.dateNoNormalize(2015, 2, 20, 14, 30),
			tg.dateNoNormalize(2015, 2, 20, 14, 0),
			tg.dateNoNormalize(2015, 2, 20, 2, 30),
			tg.dateNoNormalize(2015, 2, 20, 2, 0),
			tg.dateNoNormalize(2015, 2, 10, 14, 30),
			tg.dateNoNormalize(2015, 2, 10, 14, 0),
			tg.dateNoNormalize(2015, 2, 10, 2, 30),
			tg.dateNoNormalize(2015, 2, 10, 2, 0),
			UNKNOWN,
		})
}
