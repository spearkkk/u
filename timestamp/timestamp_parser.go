package timestamp

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/itchyny/timefmt-go"
)

var timestampParser *TimestampParser

type TimestampParser struct {
}

func init() {
	timestampParser = &TimestampParser{}
}

func GetTimestampParser() *TimestampParser {
	return timestampParser
}

func (tp *TimestampParser) parse(timestamp string, supportedFormats []string) *ParsedTimestamp {
	log.Printf("Parsing..., timestamp: %s", timestamp)

	numericTimestamp, err := parseNumericTimestamp(timestamp)
	if err == nil {
		return &ParsedTimestamp{isParsed: true, isNumericTimestamp: true, time: numericTimestamp}
	}

	formattedTimestamp, err := parseFormattedTimestamp(timestamp, supportedFormats)
	if err == nil {
		return &ParsedTimestamp{isParsed: true, isNumericTimestamp: false, time: formattedTimestamp}
	}

	return &ParsedTimestamp{isParsed: false}
}

func parseNumericTimestamp(value string) (time.Time, error) {
	if len(value) == 10 {
		seconds, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return time.Time{}, err
		}
		return time.Unix(seconds, 0), nil
	} else if len(value) == 13 {
		milliseconds, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return time.Time{}, err
		}
		return time.Unix(0, milliseconds*int64(time.Millisecond)), nil
	} else if len(value) == 16 {
		microseconds, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return time.Time{}, err
		}
		return time.Unix(0, microseconds*int64(time.Microsecond)), nil
	} else if len(value) == 19 {
		nanoseconds, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return time.Time{}, err
		}
		return time.Unix(0, nanoseconds), nil
	}
	return time.Time{}, fmt.Errorf("invalid numeric timestamp")
}

func parseFormattedTimestamp(value string, supportedFormats []string) (time.Time, error) {
	value = strings.TrimPrefix(strings.TrimSuffix(value, "'"), "'")
	for _, format := range supportedFormats {
		if parsedTime, err := time.ParseInLocation(format, value, time.Local); err == nil {
			return parsedTime, nil
		}

		if parsedTime, err := timefmt.ParseInLocation(value, format, time.Local); err == nil {
			return parsedTime, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid formatted timestamp")
}
