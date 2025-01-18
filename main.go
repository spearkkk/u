package main

import (
	"log"
	"regexp"
	"strings"

	aw "github.com/deanishe/awgo"
	ts "github.com/spearkkk/u/timestamp"
	uuid "github.com/spearkkk/u/uuid"
)

type Utility interface {
	GetKey() string
	GetName() string
	GetDescription() string
	GetResults(wf *aw.Workflow)
}

var wf *aw.Workflow
var utilities = []Utility{}

func init() {
	wf = aw.New()
	utilities = append(utilities, uuid.NewUUID())
	utilities = append(utilities, ts.NewTimestamp("", ""))
}

func main() {
	wf.Run(run)
}

func run() {
	log.Println("Running workflow...")
	// Parse input into queries
	var queries []string
	if len(wf.Args()) > 0 {
		rawQuery := strings.TrimSpace(wf.Args()[0])
		queries = parseQueries(rawQuery)
	}

	if len(queries) == 0 {
		for _, utility := range utilities {
			utility.GetResults(wf)
		}
	} else {
		utility := createUtility(queries)
		if utility != nil {
			utility.GetResults(wf)
		} else {
			wf.WarnEmpty("No matching utility found", "Please try again")
			wf.NewItem("No matching utility found").Valid(false)
		}
	}

	wf.SendFeedback()
}

func parseQueries(input string) []string {
	re := regexp.MustCompile(`"([^"]*)"|'([^']*)'|(\S+)`)
	matches := re.FindAllStringSubmatch(input, -1)
	var results []string
	for _, match := range matches {
		for _, group := range match[1:] {
			if group != "" {
				results = append(results, group)
			}
		}
	}
	return results
}

func createUtility(queries []string) Utility {
	if len(queries) < 1 {
		return nil
	}

	key := queries[0]
	value1 := ""
	value2 := ""

	if len(queries) > 1 {
		value1 = queries[1]
	}
	if len(queries) > 2 {
		value2 = queries[2]
	}

	for _, utility := range utilities {
		if utility.GetKey() == key {
			switch key {
			case "uuid":
				return uuid.NewUUID()
			case "ts":
				return ts.NewTimestamp(value1, value2)
			}
		}
	}
	return nil
}
