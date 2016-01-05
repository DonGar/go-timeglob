package timeglob

import (
	"gopkg.in/check.v1"
)

func (suite *MySuite) TestTickerStartStop(c *check.C) {
	tg, err := Parse("2012/11/25 19:37 America/New_York")
	c.Assert(err, check.IsNil)

	// Create a ticker, and then stop it right away.
	ticker := tg.Ticker()
	ticker.Stop()
}

// Proper tests require a way to mock out time.Now(). Delaying them until
// I wrap an interface around it.
