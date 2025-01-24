package timestamp

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParse(t *testing.T) {
	Convey("Given a timestamp string looks like unix(epoch) timestamp", t, func() {
		numericValue := "1737170867"

		Convey("When the numeric value be parsed", func() {
			parsedTimestamp := timestampParser.Parse(numericValue, []string{})

			Convey("Then the numeric value should be parsed", func() {
				So(parsedTimestamp.isParsed, ShouldBeTrue)
				So(parsedTimestamp.isNumericTimestamp, ShouldBeTrue)

				timestamp := parsedTimestamp.time
				So(timestamp.Year(), ShouldEqual, 2025)
				So(timestamp.Month(), ShouldEqual, 1)
				So(timestamp.Day(), ShouldEqual, 18)
				So(timestamp.Hour(), ShouldEqual, 12)
				So(timestamp.Minute(), ShouldEqual, 27)
				So(timestamp.Second(), ShouldEqual, 47)
			})
		})
	})

	Convey("Given a timestamp string looks like millisecond timestamp", t, func() {
		numericValue := "1737171674235"

		Convey("When the numeric value be parsed", func() {
			parsedTimestamp := timestampParser.Parse(numericValue, []string{})

			Convey("Then the numeric value should be parsed", func() {
				So(parsedTimestamp.isParsed, ShouldBeTrue)
				So(parsedTimestamp.isNumericTimestamp, ShouldBeTrue)

				timestamp := parsedTimestamp.time
				So(timestamp.Year(), ShouldEqual, 2025)
				So(timestamp.Month(), ShouldEqual, 1)
				So(timestamp.Day(), ShouldEqual, 18)
				So(timestamp.Hour(), ShouldEqual, 12)
				So(timestamp.Minute(), ShouldEqual, 41)
				So(timestamp.Second(), ShouldEqual, 14)
			})
		})
	})

	Convey("Given a timestamp string looks like microsecond timestamp", t, func() {
		numericValue := "1737171674235000"

		Convey("When the numeric value be parsed", func() {
			parsedTimestamp := timestampParser.Parse(numericValue, []string{})

			Convey("Then the numeric value should be parsed", func() {
				So(parsedTimestamp.isParsed, ShouldBeTrue)
				So(parsedTimestamp.isNumericTimestamp, ShouldBeTrue)

				timestamp := parsedTimestamp.time
				So(timestamp.Year(), ShouldEqual, 2025)
				So(timestamp.Month(), ShouldEqual, 1)
				So(timestamp.Day(), ShouldEqual, 18)
				So(timestamp.Hour(), ShouldEqual, 12)
				So(timestamp.Minute(), ShouldEqual, 41)
				So(timestamp.Second(), ShouldEqual, 14)
			})
		})
	})

	Convey("Given a timestamp string looks like nanosecond timestamp", t, func() {
		numericValue := "1737171674235000000"

		Convey("When the numeric value be parsed", func() {
			parsedTimestamp := timestampParser.Parse(numericValue, []string{})

			Convey("Then the numeric value should be parsed", func() {
				So(parsedTimestamp.isParsed, ShouldBeTrue)
				So(parsedTimestamp.isNumericTimestamp, ShouldBeTrue)

				timestamp := parsedTimestamp.time
				So(timestamp.Year(), ShouldEqual, 2025)
				So(timestamp.Month(), ShouldEqual, 1)
				So(timestamp.Day(), ShouldEqual, 18)
				So(timestamp.Hour(), ShouldEqual, 12)
				So(timestamp.Minute(), ShouldEqual, 41)
				So(timestamp.Second(), ShouldEqual, 14)
			})
		})
	})

	Convey("Given a timestamp string looks formatted timestamp", t, func() {
		formattedValues := []string{
			"Mon, 02 Jan 2006 15:04:05 PST",
			"2025-01-18 12:41:14",
			"2025-01-18T12:41:14+06:00",
			"2025-02-28",
			// "15:34:58",
			"2020/07/24 09:07:29",
			"20250124140909",
		}

		Convey("When the formatted value be parsed with supported formats", func() {
			for idx, formattedValue := range formattedValues {
				parsedTimestamp := timestampParser.Parse(formattedValue, []string{
					"2006-01-02 15:04:05",
					time.RFC3339,
					time.RFC1123,
					"2006-01-02",
					"15:04:05",
					"%Y/%m/%d %H:%M:%S",
					"%Y-%m-%dT%H:%M:%S%z",
					"%Y-%m-%d %H:%M:%S",
					"%Y-%m-%d",
					"%Y%m%d%H%M%S",
				})

				Convey(fmt.Sprintf("Then the formatted value should be parsed[%d]", idx), func() {
					So(parsedTimestamp.isParsed, ShouldBeTrue)
					So(parsedTimestamp.isNumericTimestamp, ShouldBeFalse)
				})
			}
		})
	})

	Convey("Given a duration string", t, func() {
		duration := "PT1H30M"
		Convey("When the duration be parsed", func() {
			parsedTimestamp := timestampParser.Parse(duration, []string{
				"2006-01-02 15:04:05",
				time.RFC3339,
				time.RFC1123,
				"2006-01-02",
				"15:04:05",
				"%Y/%m/%d %H:%M:%S",
				"%Y-%m-%dT%H:%M:%S%z",
				"%Y-%m-%d %H:%M:%S",
				"%Y-%m-%d",
			})

			Convey("Then the duration value must be not parsed", func() {
				So(parsedTimestamp.isParsed, ShouldBeFalse)
				So(parsedTimestamp.isNumericTimestamp, ShouldBeFalse)

				timestamp := parsedTimestamp.time
				Println(timestamp.String())
			})
		})
	})
}
