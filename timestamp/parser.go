package timestamp

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/itchyny/timefmt-go"
)

var timestampParser *Parser

type Parser struct {
}

func init() {
	timestampParser = &Parser{}
}

func (tp *Parser) Parse(candidate string, formats []string) *ParsedTimestamp {
	log.Printf("[Timestamp] \tparsing..., candidate: %s", candidate)

	candidate = strings.TrimSpace(candidate)

	if numericTimestamp, err := parseNumericTimestamp(candidate); err == nil {
		return &ParsedTimestamp{isParsed: true, isNumericTimestamp: true, time: numericTimestamp}
	}
	if formattedTimestamp, err := parseFormattedTimestamp(candidate, formats); err == nil {
		return &ParsedTimestamp{isParsed: true, isNumericTimestamp: false, time: formattedTimestamp}
	}

	log.Printf("[Timestamp] \tparsing failed., candidate: %s", candidate)
	return &ParsedTimestamp{isParsed: false}
}

func parseNumericTimestamp(candidate string) (time.Time, error) {
	if strings.HasPrefix(candidate, "2") {
		return time.Time{}, fmt.Errorf("invalid numeric timestamp")
	}

	if len(candidate) == 10 {
		seconds, err := strconv.ParseInt(candidate, 10, 64)
		if err != nil {
			return time.Time{}, err
		}
		return time.Unix(seconds, 0), nil
	} else if len(candidate) == 13 {
		milliseconds, err := strconv.ParseInt(candidate, 10, 64)
		if err != nil {
			return time.Time{}, err
		}
		return time.Unix(0, milliseconds*int64(time.Millisecond)), nil
	} else if len(candidate) == 16 {
		microseconds, err := strconv.ParseInt(candidate, 10, 64)
		if err != nil {
			return time.Time{}, err
		}
		return time.Unix(0, microseconds*int64(time.Microsecond)), nil
	} else if len(candidate) == 19 {
		nanoseconds, err := strconv.ParseInt(candidate, 10, 64)
		if err != nil {
			return time.Time{}, err
		}
		return time.Unix(0, nanoseconds), nil
	}
	return time.Time{}, fmt.Errorf("invalid numeric timestamp")
}

func parseFormattedTimestamp(candidate string, formats []string) (time.Time, error) {
	for _, format := range formats {
		if parsedTime, err := timefmt.ParseInLocation(candidate, format, time.Local); err == nil {
			return parsedTime, nil
		}
		if "RFC1123" == format {
			if parsed, err := time.Parse(time.RFC1123, candidate); err == nil {
				return parsed, nil
			}
		}
		if "RFC3339" == format {
			if parsed, err := time.Parse(time.RFC3339, candidate); err == nil {
				return parsed, nil
			}
		}
		if "UnixDate" == format {
			if parsed, err := time.Parse(time.UnixDate, candidate); err == nil {
				return parsed, nil
			}
		}
		if parsedTime, err := time.ParseInLocation(format, candidate, time.Local); err == nil {
			return parsedTime, nil
		}
		log.Printf("[Timestamp] \tfailed to parse formatted timestamp, candidate: %s, format: %s", candidate, format)
	}

	return time.Time{}, fmt.Errorf("invalid formatted timestamp")
}
