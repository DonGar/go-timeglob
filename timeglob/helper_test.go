package timeglob

import (
	"gopkg.in/check.v1"
	"time"
)

func (suite *MySuite) TestIntRange(c *check.C) {
	result := intRange(0, 3)
	c.Check(result, check.DeepEquals, []int{0, 1, 2, 3})

	result = intRange(1, 3)
	c.Check(result, check.DeepEquals, []int{1, 2, 3})

	result = intRange(-1, 1)
	c.Check(result, check.DeepEquals, []int{-1, 0, 1})

	result = intRange(12006, 12007)
	c.Check(result, check.DeepEquals, []int{12006, 12007})

	// begin == end means one value.
	result = intRange(12006, 12006)
	c.Check(result, check.DeepEquals, []int{12006})

	// end comes before begin.
	result = intRange(1, -1)
	c.Check(result, check.DeepEquals, []int{1, 0, -1})
}

func (suite *MySuite) TestReverseCopy(c *check.C) {
	result := reverseCopy(nil)
	c.Check(result, check.DeepEquals, []int(nil))

	result = reverseCopy([]int{})
	c.Check(result, check.DeepEquals, []int{})

	result = reverseCopy([]int{1})
	c.Check(result, check.DeepEquals, []int{1})

	result = reverseCopy([]int{1, 2})
	c.Check(result, check.DeepEquals, []int{2, 1})
}

func (suite *MySuite) TestDateNoNormalize(c *check.C) {
	tg, err := Parse("2010/1/1 America/New_York")
	c.Assert(err, check.IsNil)

	result := tg.dateNoNormalize(2015, 12, 25, 1, 15, 47)
	c.Check(result.Year(), check.Equals, 2015)
	c.Check(result.Month(), check.Equals, time.December)
	c.Check(result.Day(), check.Equals, 25)
	c.Check(result.Hour(), check.Equals, 1)
	c.Check(result.Minute(), check.Equals, 15)
	c.Check(result.Second(), check.Equals, 47)

	// Out of bounds.
	result = tg.dateNoNormalize(2015, 13, 25, 1, 12, 0)
	c.Check(result, check.Equals, UNKNOWN)

	result = tg.dateNoNormalize(2015, 12, 32, 1, 12, 0)
	c.Check(result, check.Equals, UNKNOWN)

	result = tg.dateNoNormalize(2015, 12, 25, 24, 12, 0)
	c.Check(result, check.Equals, UNKNOWN)

	result = tg.dateNoNormalize(2015, 12, 25, 1, 62, 0)
	c.Check(result, check.Equals, UNKNOWN)

	result = tg.dateNoNormalize(2015, 12, 25, 1, 12, 62)
	c.Check(result, check.Equals, UNKNOWN)

	// Leap Days
	result = tg.dateNoNormalize(2015, 2, 29, 0, 0, 0)
	c.Check(result, check.Equals, UNKNOWN)

	result = tg.dateNoNormalize(2016, 2, 29, 0, 0, 0)
	c.Check(result.Year(), check.Equals, 2016)
	c.Check(result.Month(), check.Equals, time.February)
	c.Check(result.Day(), check.Equals, 29)
	c.Check(result.Hour(), check.Equals, 0)
	c.Check(result.Minute(), check.Equals, 0)
}

func (suite *MySuite) TestAdjustMinutesSeconds(c *check.C) {
	tg, err := Parse("2010/1/1 America/New_York")
	c.Assert(err, check.IsNil)
	base := tg.dateNoNormalize(2015, 2, 3, 4, 0, 0)

	// Valid
	result := tg.adjustMinutesSeconds(base, 5, 10)
	c.Check(result, check.Equals, tg.dateNoNormalize(2015, 2, 3, 4, 5, 10))

	// Out of bounds.
	result = tg.adjustMinutesSeconds(base, 61, 10)
	c.Check(result, check.Equals, UNKNOWN)

	result = tg.adjustMinutesSeconds(base, 5, 61)
	c.Check(result, check.Equals, UNKNOWN)
}
