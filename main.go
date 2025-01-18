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
var utilities = []Utility{}

func init() {
	flag.BoolVar(&doCheck, "check", false, "Check for updates")

	wf = aw.New(update.GitHub(repo))

	utilities = append(utilities, uuid.NewUUID())
	utilities = append(utilities, ts.NewTimestamp("", ""))
}

func main() {
	wf.Run(run)
}

func run() {
	log.Println("Running workflow...")

	flag.Parse()
	query = flag.Arg(0)

	if wf.UpdateCheckDue() && !wf.IsRunning(updateJobName) {
		log.Println("Running update check in background...")

		cmd := exec.Command(os.Args[0], "-check")
		if err := wf.RunInBackground(updateJobName, cmd); err != nil {
			log.Printf("Error starting update check: %s", err)
		}
	}

	// Only show update status if query is empty.
	if query == "" && wf.UpdateAvailable() {
		// Turn off UIDs to force this item to the top.
		// If UIDs are enabled, Alfred will apply its "knowledge"
		// to order the results based on your past usage.
		wf.Configure(aw.SuppressUIDs(true))

		// Notify user of update. As this item is invalid (Valid(false)),
		// actioning it expands the query to the Autocomplete value.
		// "workflow:update" triggers the updater Magic Action that
		// is automatically registered when you configure Workflow with
		// an Updater.
		//
		// If executed, the Magic Action downloads the latest version
		// of the workflow and asks Alfred to install it.
		wf.NewItem("Update available!").
			Subtitle("↩ to install").
			Autocomplete("workflow:update").
			Valid(false)
	}

	// Add an extra item to reset update status for demo purposes.
	// As with the update notification, this item triggers a Magic
	// Action that deletes the cached list of releases.
	wf.NewItem("Reset update status").
		Autocomplete("workflow:delcache").
		Valid(false)

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
