package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/update"
	ts "github.com/spearkkk/u/timestamp"
	uuid "github.com/spearkkk/u/uuid"
)

type Utility interface {
	GetKey() string
	GetName() string
	GetDescription() string
	GetResults(wf *aw.Workflow)
}

var updateJobName string = "update"
var repo string = "spearkkk/u"
var doCheck bool
var query string

var wf *aw.Workflow
var utilities []Utility

func init() {
	flag.BoolVar(&doCheck, "check", false, "Check for updates")
	wf = aw.New(update.GitHub(repo))

	log.Println("cache dir: ", wf.CacheDir())
	log.Println("data dir: ", wf.DataDir())
	log.Println("bundle id: ", wf.BundleID())
	log.Println("Config: ", wf.Config)

	utilities = append(utilities, uuid.NewUUID())
	utilities = append(utilities, ts.NewTimestamp("", ""))
}

func main() {
	wf.Run(run)
}

func run() {
	log.Println("Running workflow...")

	if wf.UpdateCheckDue() && !wf.IsRunning(updateJobName) {
		log.Println("Running update check in background...")

		cmd := exec.Command(os.Args[0], "-check")
		if err := wf.RunInBackground(updateJobName, cmd); err != nil {
			log.Printf("Error starting update check: %s", err)
		}
	}

	if doCheck {
		if err := wf.CheckForUpdate(); err != nil {
			log.Printf("Error checking for update: %s", err)
		}
		return
	}

	if query == "" && wf.UpdateAvailable() {
		wf.Configure(aw.SuppressUIDs(true))
		wf.NewItem("Update available!").
			Subtitle("↩ to install").
			Autocomplete("workflow:update").
			Valid(false)
	}

	// Access Alfred configuration variables
	keyToEnabled := map[string]bool{
		"ts":   wf.Config.GetBool("ts"),
		"uuid": wf.Config.GetBool("uuid"),
	}
	tsFormats := wf.Config.Get("ts_formats")

	log.Printf("Utility enabled: %v", keyToEnabled)
	log.Printf("Timestamp formats: %s", tsFormats)

	// Parse input into queries
	var queries []string
	if len(wf.Args()) > 0 {
		rawQuery := strings.TrimSpace(wf.Args()[0])
		queries = parseQueries(rawQuery)
	}

	if len(queries) == 0 {
		for _, utility := range utilities {
			if keyToEnabled[utility.GetKey()] {
				utility.GetResults(wf)
			}
		}
	} else {
		utility := createUtility(queries)
		if utility != nil && keyToEnabled[utility.GetKey()] {
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
