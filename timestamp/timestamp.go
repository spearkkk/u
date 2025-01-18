package timestamp

import (
	"fmt"
	"log"
	"strings"
	"time"

	aw "github.com/deanishe/awgo"
)

var supportedFormats []string

type Timestamp struct {
	value1 string
	value2 string
}

func init() {
	supportedFormats = []string{
		"2006-01-02 15:04:05",
		time.RFC3339,
		time.RFC1123,
		"2006-01-02",
	}
}

func NewTimestamp(value1, value2 string, formats ...string) *Timestamp {
	if len(formats) != 0 {
		supportedFormats = formats
	}
	return &Timestamp{value1: value1, value2: value2}
}

func (t *Timestamp) GetKey() string {
	return "ts"
}

func (t *Timestamp) GetName() string {
	return "Timestamp"
}

func (t *Timestamp) GetDescription() string {
	return "[copy:⏎, next:↹] To get/convert/calculate timestamp."
}
func (t *Timestamp) GetResults(wf *aw.Workflow) {
	log.Printf("Processing..., values: %s %s", t.value1, t.value2)
	// t.value2 must be empty too
	if t.value1 == "" {
		now := time.Now()
		currentMilliseconds := time.Now().UnixMilli()

		t.setResult(wf, append([]string{fmt.Sprintf("%d", currentMilliseconds)}, formatTime(now)...))
		return
	}

	parsedValue1 := timestampParser.parse(t.value1, supportedFormats)

	isUnaryOperation := t.value2 == ""
	if isUnaryOperation {
		if parsedValue1.isParsed {
			if parsedValue1.isNumericTimestamp {
				t.setResult(wf, formatTime(parsedValue1.time))
				return
			}
			t.setResult(wf, []string{fmt.Sprintf("%d", parsedValue1.time.UnixMilli())})
			return
		}
		t.setResult(wf, []string{fmt.Sprintf("Invalid value(%s)", t.value1)})
		return
	}

	parsedValue2 := timestampParser.parse(t.value2, supportedFormats)
	if !parsedValue1.isParsed && !parsedValue2.isParsed {
		t.setResult(wf, []string{parsedValue1.time.Sub(parsedValue2.time).String()})
		return
	}

	if parsedValue1.isParsed {
		// parsedValue2 is not timestamp
		results, err := processRawValue(parsedValue1.time, t.value2)
		if err != nil {
			t.setResult(wf, []string{fmt.Sprintf("Invalid values(%s %s)", t.value1, t.value2)})
			return
		}
		t.setResult(wf, results)
		return
	}

	// parsedValue1 is not timestamp
	results, err := processRawValue(parsedValue2.time, t.value1)
	if err != nil {
		t.setResult(wf, []string{fmt.Sprintf("Invalid values(%s %s)", t.value1, t.value2)})
		return
	}
	t.setResult(wf, results)
}

func (t *Timestamp) setResult(wf *aw.Workflow, results []string) {
	for _, result := range results {
		escapedResult := result
		if strings.ContainsAny(result, " \t\n\r") {
			escapedResult = "'" + result + "'"
		}

		wf.NewItem(t.GetName() + ": " + result).
			Subtitle(t.GetDescription()).
			Arg(result).
			Copytext(result).
			Quicklook(result).
			Valid(true).
			Autocomplete(t.GetKey() + " " + escapedResult + " ")
	}
}

func processRawValue(time time.Time, rawValue string) ([]string, error) {
	log.Printf("Processing raw value..., time: %s, rawValue: %s", time.String(), rawValue)

	parseDuration := durationParser.parseDuration(rawValue)
	log.Printf("Parsed duration: %+v", parseDuration)
	if parseDuration.isParsed {
		var processedTime = time
		if parseDuration.isShifting {
			processedTime = time.Add(parseDuration.duration.ToTimeDuration())
		} else {
			processedTime = time.Add(-parseDuration.duration.ToTimeDuration())
		}

		return append([]string{fmt.Sprintf("%d", processedTime.UnixMilli())}, formatTime(processedTime)...), nil
	}

	reformattedTime := timestampFormatter.Format(rawValue, time)
	if reformattedTime != "" {
		return []string{reformattedTime}, nil
	}

	log.Printf("Invalid value(%s)", rawValue)
	return []string{}, fmt.Errorf("invalid value(%s)", rawValue)
}

func formatTime(timestamp time.Time) []string {
	var formattedTimes []string
	for _, format := range supportedFormats {
		tmp := timestamp.Format(format)
		formattedTimes = append(formattedTimes, tmp)
	}
	return formattedTimes
}
