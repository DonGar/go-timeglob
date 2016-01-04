# go-timeglob

This package defines an easy syntax for defining repeating times, in a mannor similar to cron, but with a more friendly syntax.

Standard usage is to parse a text string to create a TimeGlob structure, then ask the structure for Next or Previous matches from a given moment. This structure does NOT match ranges of time, only descrete moments, and has no finer resolution than minutes.

## Format ##

Example globs:
* "2015/12/25 19:37 America/New_York"
* "*/1 19:37"
*	"19:37",
*	"12/25",
*	"12/25 *:0,15,30,45",

Any field can be a wildcard '*', which matches any possible value. Any field can contain multiple values seperated by ',', which means any value in the list is a match.

The date can be specified as year/month/day, or month/day. If the year isn't specified, it defaults to '*'. If the date isn't present, it defaults to '*/*/*'.

The time is specified as hour:minute, where hour is in 24 hour time. If time isn't present, it defaults to 0:0 (midnight at the start of the day). During special cases in which an hour can repeat (like start/end of daylight savings times) , a '*' will match both instances of the hour, but an explicit value will only match the first.

The timezone is any timezone name supported by Go's [time.LoadLocation](https://golang.org/pkg/time/#LoadLocation) is supported (including 'Local' and 'UTC'). If not present 'Local' is used.

## Next/Prev ##

After parsing a glob, the operations available are Next, and Prev. Both methods require a 'now' argument that they calculate from. they return timeglob.UNKNOWN if there is no valid match. The result will be in the same timezone as the 'now' argument, which does not have to match the timezone of the glob.

* Next: > now
* Prev: <= now

## Notes ##

This library was knocked together during a holiday vacation, it probabably has bugs ([time handling is hard](http://infiniteundo.com/post/25326999628/falsehoods-programmers-believe-about-time)). I welcome bug reports, or better yet, test cases that reproduce problems.

This library 'searches' for Next/Prev in a way that could provide really terrible performance for some edge cases. I haven't yet done anything to address this.

I'll probably add something similar to [time.Ticker](https://golang.org/pkg/time/#Ticker) to this library in time.
