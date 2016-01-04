package timeglob

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func Parse(glob string) (*TimeGlob, error) {
	result := new()
	sections := strings.SplitN(glob, " ", 3)

	if len(sections) > 0 {
		if result.parseDate(sections[0]) {
			sections = sections[1:]
		}
	}

	if len(sections) > 0 {
		if result.parseTime(sections[0]) {
			sections = sections[1:]
		}
	}

	if len(sections) > 0 {
		if result.parseLocation(sections[0]) {
			sections = sections[1:]
		}
	}

	if len(sections) > 0 {
		return nil, fmt.Errorf("Not a valid TimeGlob: %s", glob)
	}

	return &result, nil
}

func parseIntList(blob string) []int {
	if blob == "*" || blob == "" {
		return nil
	}

	sections := strings.Split(blob, ",")
	values := map[int]bool{}

	for _, s := range sections {
		if s == "" {
			continue
		}

		val, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			panic(err)
		}
		values[int(val)] = true
	}

	// Convert map to sorted slice (nil if empty).
	result := []int{}
	for key := range values {
		result = append(result, key)
	}
	sort.Ints(result)
	return result
}

func (tg *TimeGlob) parseDate(glob string) bool {
	re := regexp.MustCompile(`^(([0-9,]+|\*)/)?([0-9,]+|\*)/([0-9,]+|\*)$`)
	submatches := re.FindStringSubmatch(glob)

	if submatches == nil {
		return false
	}

	tg.year = parseIntList(submatches[2])
	tg.month = parseIntList(submatches[3])
	tg.day = parseIntList(submatches[4])
	return true
}

func (tg *TimeGlob) parseTime(glob string) bool {
	re := regexp.MustCompile(`^([0-9,]+|\*):([0-9,]+|\*)$`)
	submatches := re.FindStringSubmatch(glob)

	if submatches == nil {
		return false
	}

	tg.hour = parseIntList(submatches[1])
	tg.minute = parseIntList(submatches[2])
	return true
}

func (tg *TimeGlob) parseLocation(glob string) bool {
	loc, err := time.LoadLocation(glob)
	if glob != "" && err == nil {
		tg.location = loc
		return true
	}
	return false
}
