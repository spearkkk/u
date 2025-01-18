package timestamp

import (
	"strings"

	"github.com/sosodev/duration"
)

var durationParser *DurationParser

func init() {
	durationParser = &DurationParser{}
}

type DurationParser struct {
}

func GetDurationParser() *DurationParser {
	return durationParser
}

func (dp *DurationParser) parseDuration(value string) *ParsedDuration {
	isShifting := true
	if strings.HasPrefix(value, "-") {
		value = value[1:]
		isShifting = false
	}

	value = strings.TrimPrefix(value, "+")

	parsedDuration, err := duration.Parse(value)
	if err == nil {
		return &ParsedDuration{isParsed: true, isShifting: isShifting, duration: *parsedDuration}
	}
	return &ParsedDuration{isParsed: false}
}