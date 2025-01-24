package timestamp

import (
	"log"
	"strings"
	"time"

	"github.com/itchyny/timefmt-go"
)

var timestampFormatter *Formatter

type Formatter struct {
}

func init() {
	timestampFormatter = &Formatter{}
}

func (tf *Formatter) Format(format string, timestamp time.Time) string {
	log.Println("format: " + format)
	if strings.ContainsAny(format, "%") {
		return timefmt.Format(timestamp, format)
	}
	if "RFC1123" == format {
		log.Println("RFC1123" + timestamp.Format(time.RFC1123))
		return timestamp.Format(time.RFC1123)
	}
	if "RFC3339" == format {
		return timestamp.Format(time.RFC3339)
	}
	if "UnixDate" == format {
		return timestamp.Format(time.UnixDate)
	}
	return timestamp.Format(iso8601ToGoLayout(format))
}

// Convert ISO 8601-like format to Go layout (handling MM vs mm accurately)
func iso8601ToGoLayout(format string) string {
	// Define a map of ISO 8601 placeholders to Go layout equivalents
	mapping := map[string]string{
		"YYYY": "2006",   // Year (4 digits)
		"YY":   "06",     // Year (2 digits)
		"MM":   "01",     // Month (2 digits)
		"DD":   "02",     // Day (2 digits)
		"hh":   "15",     // Hour (24-hour)
		"mm":   "04",     // Minute
		"ss":   "05",     // Second
		"TZD":  "Z07:00", // Timezone offset
		"Z":    "Z07:00", // Timezone UTC
		".s":   ".999",   // Fractional seconds
	}

	// Preserve case-sensitive replacements (e.g., MM vs mm)
	// Avoid global case conversion and manually check against the mapping
	for key, value := range mapping {
		format = strings.ReplaceAll(format, key, value)
	}

	return format
}
