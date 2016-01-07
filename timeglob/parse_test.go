package timeglob

import (
	"gopkg.in/check.v1"
	"time"
)

func (suite *MySuite) TestParseGlobParseGood(c *check.C) {
	globs := []string{
		"2015/12/25 19:37:22 America/New_York",
		"2015/12/25 19:37 America/New_York",
		"2015/12/25 19:37 Local",
		"2015/12/25 19:37 UTC",
		"12/25 19:37 UTC",
		"2015/12/25 19:37",
		"12/25 19:37",
		"2015/12/25",
		"12/25",
		"19:37",
		"19:37:22",
		"*/*/* *:*:* UTC",
		"*/*/* *:* UTC",
		"*/*/* *:*",
		",/,/, ,:, UTC",
		"2014,2015/11,12/22,25 8,19:15,37:11,22 America/New_York",
		"2014,2015/11,12/22,25 8,19:15,37 America/New_York",
		"2015,/12/25,25 10,:37 America/New_York",
		"2015/12/25 ,:37 America/New_York",
		",2015/12/25 19:37 America/New_York",
	}

	for _, g := range globs {
		tg, err := Parse(g)
		c.Check(tg, check.NotNil)
		c.Check(err, check.IsNil)
	}
}

func (suite *MySuite) TestParseGlobParseBad(c *check.C) {
	globs := []string{
		"",
		"Known Bad Glob",
		"2015/12/25/12 19:37 America/New_York",
		"12 19:37 America/New_York",
		"2015/12/25 19:37 Bad/Local",
		" 2015/12/25 19:37 America/New_York",
		"2015/12/25 19:37 America/New_York ",
		"Extra 2015/12/25 19:37 America/New_York",
		"2015/12/25 Extra 19:37 America/New_York",
		"2015/12/25 19:37 Extra America/New_York",
		"2015/12/25 19:37 America/New_York Extra",
		"2015/12/25 aa:37 America/New_York",
	}

	for _, g := range globs {
		tg, err := Parse(g)
		c.Check(tg, check.IsNil)
		c.Check(err, check.NotNil)
	}
}

func matchesExpected(c *check.C, glob string, expected *TimeGlob) {
	tg, err := Parse(glob)
	c.Check(tg, check.NotNil)
	c.Check(err, check.IsNil)
	c.Check(tg, check.DeepEquals, expected)
}

func (suite *MySuite) TestParseGlobParseVerify(c *check.C) {
	matchesExpected(c, "2015/12/25 19:37:22 UTC", &TimeGlob{
		[]int{2015}, []int{12}, []int{25},
		[]int{19}, []int{37}, []int{22},
		time.UTC,
	})

	matchesExpected(c, "2015/12/25 UTC", &TimeGlob{
		[]int{2015}, []int{12}, []int{25},
		intRange(0, 0), intRange(0, 0), []int{0},
		time.UTC,
	})

	matchesExpected(c, "12/25 UTC", &TimeGlob{
		nil, []int{12}, []int{25},
		intRange(0, 0), intRange(0, 0), []int{0},
		time.UTC,
	})

	matchesExpected(c, "19:37:22 UTC", &TimeGlob{
		nil, nil, nil,
		[]int{19}, []int{37}, []int{22},
		time.UTC,
	})

	matchesExpected(c, "19:37:* UTC", &TimeGlob{
		nil, nil, nil,
		[]int{19}, []int{37}, nil,
		time.UTC,
	})

	matchesExpected(c, "19:37 UTC", &TimeGlob{
		nil, nil, nil,
		[]int{19}, []int{37}, []int{0},
		time.UTC,
	})

	// matchesExpected(c, "2015/12/25 19:37", &TimeGlob{
	//   []int{2015}, []int{12}, []int{25},
	//   []int{19}, []int{37},
	//   time.UTC,
	// })

	matchesExpected(c, "2015,2016/11,12/22,25 8,19:22,37 UTC", &TimeGlob{
		[]int{2015, 2016}, []int{11, 12}, []int{22, 25},
		[]int{8, 19}, []int{22, 37}, []int{0},
		time.UTC,
	})

	matchesExpected(c, "2015,2016,/11,11,12/25,22,25 19,8:22,37:11,22 UTC", &TimeGlob{
		[]int{2015, 2016}, []int{11, 12}, []int{22, 25},
		[]int{8, 19}, []int{22, 37}, []int{11, 22},
		time.UTC,
	})

	matchesExpected(c, ",/,/, ,:,:, UTC", &TimeGlob{
		[]int{}, []int{}, []int{},
		[]int{}, []int{}, []int{},
		time.UTC,
	})

}

func (suite *MySuite) TestParseGlobParseWildcardEquivalence(c *check.C) {

	testEquivalence := func(globs []string) {
		expected, err := Parse(globs[0])
		c.Check(err, check.IsNil)

		for _, g := range globs {
			matchesExpected(c, g, expected)
		}
	}

	wildcards := []string{
		"*/*/* *:* Local",
		"*/* *:* Local",
		"*/*/* *:*",
		"*/* *:*",
		"*:* Local",
		"*:*",
	}

	full := []string{
		"2015/12/25 19:37 UTC",
		",2015,/,12,/,25, ,19,:,37, UTC",
	}

	time := []string{
		"*/*/* 19:37 UTC",
		"*/* 19:37 UTC",
		"19:37 UTC",
		",19,:,37, UTC",
	}

	testEquivalence(wildcards)
	testEquivalence(full)
	testEquivalence(time)
}
