package timestamp

import (
	"time"

	"github.com/sosodev/duration"
)

type ParsedTimestamp struct {
	isParsed           bool
	isNumericTimestamp bool
	time               time.Time
}

type ParsedDuration struct {
	isParsed   bool
	duration   duration.Duration
	isShifting bool
}