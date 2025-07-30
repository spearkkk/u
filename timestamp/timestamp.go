package timestamp

import (
	"fmt"
	"log"
	"time"

	aw "github.com/deanishe/awgo"
)

type Timestamp struct {
	value1  string
	value2  string
	formats []string
}

func NewTimestamp(value1, value2 string, formats ...string) *Timestamp {
	return &Timestamp{value1: value1, value2: value2, formats: formats}
}

func (t *Timestamp) Key() string {
	return "ts"
}

func (t *Timestamp) Do(wf *aw.Workflow) {
	// Prevent to order by UID
	wf.Configure(aw.SuppressUIDs(true))

	log.Printf("[Timestamp] \tprocessing..., values: %s %s", t.value1, t.value2)
	now := time.Now()

	// t.value2 must be empty too
	if t.value1 == "" || t.value1 == "now" {
		t.setMillisecondValue(wf, now)
		return
	}

	parsedTimestamp1 := timestampParser.Parse(t.value1, t.formats)
	isUnaryOperation := t.value2 == ""
	if isUnaryOperation {
		log.Printf("[Timestamp] \tprocessing unary operation..., value1: %s", t.value1)

		if parsedTimestamp1.isParsed {
			log.Printf("[Timestamp] \tprocessing unary operation..., parsedTimestamp1: %+v", parsedTimestamp1)

			if parsedTimestamp1.isNumericTimestamp {
				t.setFormattedValues(wf, t.formatTime(parsedTimestamp1.time))
				return
			}

			t.setMillisecondValue(wf, parsedTimestamp1.time)
			t.setFormattedValues(wf, t.formatTime(parsedTimestamp1.time))
			return
		} else {
			parsedDuration1 := durationParser.parseDuration(t.value1)
			if parsedDuration1.isParsed {
				log.Printf("[Timestamp] \tprocessing unary operation..., parsedDuration1: %+v", parsedDuration1)

				millisecond := now.UnixMilli()
				if parsedDuration1.isShifting {
					millisecond += parsedDuration1.duration.ToTimeDuration().Milliseconds()
				} else {
					millisecond -= parsedDuration1.duration.ToTimeDuration().Milliseconds()
				}

				t.setFormattedValues(wf, t.formatTime(time.UnixMilli(millisecond)))
				return
			}

			t.setInvalidValue(wf)

			// fallback logic for formatting string value
			maybeFormat := t.value1
			maybeFormattedTime := t.formatTimeWithFormat(now, maybeFormat)
			// if formattedTime does not equal the format string, it means the format could be valid
			if maybeFormat != maybeFormattedTime[maybeFormat] {
				t.setFormattedValues(wf, maybeFormattedTime)
			}
			return
		}
	}

	parsedTimestamp2 := timestampParser.Parse(t.value2, t.formats)

	if !parsedTimestamp1.isParsed && !parsedTimestamp2.isParsed {
		t.setInvalidValue(wf)
		return
	}

	if parsedTimestamp1.isParsed && parsedTimestamp2.isParsed {
		result := parsedTimestamp1.time.Sub(parsedTimestamp2.time).String()
		wf.NewItem(result).
			Subtitle("Get diff").
			Arg(result).
			Copytext(result).
			Quicklook(result).
			Valid(true)
		return
	}

	if parsedTimestamp1.isParsed {
		// parsedTimestamp2 is not timestamp
		err := t.processRawValue(wf, parsedTimestamp1.time, t.value2)
		if err != nil {
			t.setInvalidValue(wf)
			return
		}
		return
	}

	// parsedTimestamp1 is not timestamp
	err := t.processRawValue(wf, parsedTimestamp2.time, t.value1)
	if err != nil {
		t.setInvalidValue(wf)
		return
	}
}

func (t *Timestamp) processRawValue(wf *aw.Workflow, timestamp time.Time, rawValue string) error {
	log.Printf("[Timestamp] \tprocessing binary operation..., value1: %s, value2: %s", t.value1, t.value2)

	parsedDuration := durationParser.parseDuration(rawValue)
	log.Printf("[Timestamp] \tprocessing binary operation..., parsed duration: %+v", parsedDuration)

	if parsedDuration.isParsed {
		var processedTime = timestamp

		millisecond := processedTime.UnixMilli()
		if parsedDuration.isShifting {
			millisecond += parsedDuration.duration.ToTimeDuration().Milliseconds()
		} else {
			millisecond -= parsedDuration.duration.ToTimeDuration().Milliseconds()
		}

		t.setFormattedValues(wf, t.formatTime(time.UnixMilli(millisecond)))
		t.setMillisecondValue(wf, time.UnixMilli(millisecond))

		return nil
	}

	reformattedTime := timestampFormatter.Format(rawValue, timestamp)
	if reformattedTime != "" {
		t.setFormattedValues(wf, map[string]string{rawValue: reformattedTime})
		return nil
	}

	return fmt.Errorf("cannot Parse duration or format(%s)", rawValue)
}

func (t *Timestamp) formatTime(timestamp time.Time) map[string]string {
	formatToValue := make(map[string]string)
	for _, format := range t.formats {
		formatToValue[format] = timestampFormatter.Format(format, timestamp)
	}
	return formatToValue
}

func (t *Timestamp) formatTimeWithFormat(timestamp time.Time, format string) map[string]string {
	formatToValue := make(map[string]string)
	formatToValue[format] = timestampFormatter.Format(format, timestamp)
	return formatToValue
}

func (t *Timestamp) setFormattedValues(wf *aw.Workflow, formatToValue map[string]string) {
	for key, result := range formatToValue {
		wf.NewItem(result).
			Subtitle(key).
			Arg(result).
			Copytext(result).
			Quicklook(result).
			Valid(true).
			Autocomplete(fmt.Sprintf("%s %s ", t.Key(), result))
	}
}

func (t *Timestamp) setMillisecondValue(wf *aw.Workflow, timestamp time.Time) {
	result := fmt.Sprintf("%d", timestamp.UnixMilli())
	wf.NewItem(result).
		Subtitle("Get millisecond.").
		Arg(result).
		Copytext(result).
		Quicklook(result).
		Valid(true).
		Autocomplete(fmt.Sprintf("%s %s ", t.Key(), result))
}

func (t *Timestamp) setInvalidValue(wf *aw.Workflow) {
	wf.NewItem(fmt.Sprintf("Invalid Value: %s %s", t.value1, t.value2)).
		Subtitle("u ts [TIMESTAMP|DURATION]").
		Valid(false).
		Autocomplete(fmt.Sprintf("%s now", t.Key()))
}
