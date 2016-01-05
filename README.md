# go-timeglob

This package defines an easy syntax for defining repeating times, in a mannor similar to cron, but with a more friendly syntax.

Standard usage is to parse a text string to create a TimeGlob structure, then ask the structure for Next or Previous matches from a given moment. This structure does NOT match ranges of time, only descrete moments, and has no finer resolution than minutes.

## Format ##

Example globs:
* "2015/12/25 19:37 America/New_York"

  Matche a single point in time, as specified in a given timezone.

* "*/1 19:37"

  Match the first of the month at 7:37 PM in the local timezone.

*	"19:37",

  Match 7:37 PM every day.

*	"12/25",

  Match midnight at the start of December 25th of every year.

*	"12/25 \*:0,15,30,45",

  Match every 15 minutes all day on December 25th of every year.

Any field can be a wildcard \*, which matches any possible value. Any field can
contain multiple values seperated by comma. Any value in the list is a match.

The date can be specified as year/month/day, or month/day. If the year isn't
specified, it defaults to \*. If the date isn't present, it defaults to
\*/\*/\*.

The time is specified as hour:minute, where hour is in 24 hour time. If time
isn't present, it defaults to 0:0 (midnight at the start of the day). During
special cases in which an hour can repeat (like start/end of daylight savings
times) , a \* will match both instances of the hour, but an explicit value will
only match the first. This is intended to match the intuitive expectations of
'run once an hour, every hour' or 'run once a day at given time'.

The timezone is any timezone name supported by Go's
[time.LoadLocation](https://golang.org/pkg/time/#LoadLocation) is supported
(including 'Local' and 'UTC'). If not present 'Local' is used.

## Next/Prev ##

After parsing a glob, the operations available are Next, and Prev. Both methods
require a 'now' argument that they calculate from. they return timeglob.UNKNOWN
if there is no valid match. The result will be in the same timezone as the 'now'
argument, which does not have to match the timezone of the glob.

* Next: > now
* Prev: <= now

## Ticker ##

Calling Ticker() on a TimeGlob returns an object that sends time values on
ticker.C for each time glob match. It must be stopped with "Stop()" to release
resources.

## TODOs ##
* Improve parsing error messages.
* Add value bounds checking during parsing.
* Add day of week handling.
* Performance is generally good, but can degrade badly in some edge cases.
  Address.

## Notes ##

This library was knocked together during a vacation, it probabably has bugs
([time handling is hard](http://infiniteundo.com/post/25326999628/falsehoods-
programmers-believe-about-time)). I welcome bug reports, or better yet, test
cases that reproduce problems.
