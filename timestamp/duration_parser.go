package timestamp

import (
	"fmt"
	"github.com/sosodev/duration"

	"strconv"
	"strings"
)

var durationParser *DurationParser

func init() {
	durationParser = &DurationParser{}
}

type DurationParser struct {
}

func (dp *DurationParser) parseDuration(candidate string) *ParsedDuration {
	isShifting := true
	if strings.HasPrefix(candidate, "-") {
		candidate = candidate[1:]
		isShifting = false
	}

	candidate = strings.TrimPrefix(candidate, "+")

	parsedDuration, err := duration.Parse(candidate)
	if err == nil {
		return &ParsedDuration{isParsed: true, isShifting: isShifting, duration: *parsedDuration}
	}

	candidate = strings.ToLower(candidate)
	switch {
	case strings.HasSuffix(candidate, "y"), strings.HasSuffix(candidate, "year"), strings.HasPrefix(candidate, "y"), strings.HasPrefix(candidate, "year"):
		if n, err := strconv.Atoi(strings.Trim(strings.Trim(candidate, "year"), "y")); err == nil {
			n = n * 365 * 24 * 60 * 60
			return convertToDuration(n, isShifting)
		}
	case strings.HasSuffix(candidate, "mo"), strings.HasSuffix(candidate, "month"), strings.HasPrefix(candidate, "mo"), strings.HasPrefix(candidate, "month"):
		if n, err := strconv.Atoi(strings.Trim(strings.Trim(candidate, "month"), "mo")); err == nil {
			n = n * 30 * 24 * 60 * 60
			return convertToDuration(n, isShifting)
		}
	case strings.HasSuffix(candidate, "d"), strings.HasSuffix(candidate, "day"), strings.HasPrefix(candidate, "d"), strings.HasPrefix(candidate, "day"):
		if n, err := strconv.Atoi(strings.Trim(strings.Trim(candidate, "day"), "d")); err == nil {
			n = n * 24 * 60 * 60
			return convertToDuration(n, isShifting)
		}
	case strings.HasSuffix(candidate, "h"), strings.HasSuffix(candidate, "hour"), strings.HasPrefix(candidate, "h"), strings.HasPrefix(candidate, "hour"):
		if n, err := strconv.Atoi(strings.Trim(strings.Trim(candidate, "hour"), "h")); err == nil {
			n = n * 60 * 60
			return convertToDuration(n, isShifting)
		}
	case strings.HasSuffix(candidate, "m"), strings.HasSuffix(candidate, "minute"), strings.HasPrefix(candidate, "m"), strings.HasPrefix(candidate, "minute"):
		if n, err := strconv.Atoi(strings.Trim(strings.Trim(candidate, "minute"), "m")); err == nil {
			n = n * 60
			return convertToDuration(n, isShifting)
		}
	case strings.HasSuffix(candidate, "s"), strings.HasSuffix(candidate, "second"), strings.HasPrefix(candidate, "s"), strings.HasPrefix(candidate, "second"):
		if n, err := strconv.Atoi(strings.Trim(strings.Trim(candidate, "second"), "s")); err == nil {
			return convertToDuration(n, isShifting)
		}
	}
	return &ParsedDuration{isParsed: false}
}

func convertToDuration(second int, isShifting bool) *ParsedDuration {
	if parsedDuration, err := duration.Parse(fmt.Sprintf("PT%dS", second)); err == nil {
		return &ParsedDuration{isParsed: true, isShifting: isShifting, duration: *parsedDuration}
	}
	return &ParsedDuration{isParsed: false}
}
