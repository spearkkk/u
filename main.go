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
)

var updateJobName = "update"
var doCheck bool

var wf *aw.Workflow
var utilities []Utility

func init() {
	flag.BoolVar(&doCheck, "check", true, "Check for updates")
	wf = aw.New(update.GitHub("spearkkk/u"))

	utilities = append(utilities, createUtility([]string{"uuid", "", ""}, map[string]interface{}{}))
	utilities = append(utilities, createUtility([]string{"ts", "", ""}, map[string]interface{}{}))
	utilities = append(utilities, createUtility([]string{"json", "", ""}, map[string]interface{}{}))
	utilities = append(utilities, createUtility([]string{"c", "", ""}, map[string]interface{}{}))
}

func main() {
	wf.Run(run)
}

func run() {
	log.Println("########################################################################################")
	log.Println()
	log.Println("Hello, this is utility workflow!")
	log.Println("Feel free to use the workflow!")
	log.Println("Please report any issues to here, https://github.com/spearkkk/u/issues.")
	log.Println("Thanks, Happy coding!")
	log.Println()
	log.Println("########################################################################################")

	if doCheck {
		//wf.Configure(aw.TextErrors(true))
		log.Println("Checking for updates...")
		if err := wf.CheckForUpdate(); err != nil {
			log.Printf("Error checking for update: %s", err)
		}
	}

	if wf.UpdateCheckDue() && !wf.IsRunning(updateJobName) {
		log.Println("Running in background for checking update...")

		cmd := exec.Command(os.Args[0], "-check")
		if err := wf.RunInBackground(updateJobName, cmd); err != nil {
			log.Printf("Error starting update check: %s", err)
		}
	}

	if wf.UpdateAvailable() {
		wf.Configure(aw.SuppressUIDs(true))
		wf.NewItem("New Version Available(↩)").
			Autocomplete("workflow:update").
			Valid(false)
	}

	// Access Alfred configuration variables
	keyToEnabled := map[string]bool{
		"ts":   wf.Config.GetBool("ts", true),
		"uuid": wf.Config.GetBool("uuid", true),
		"json": wf.Config.GetBool("json", true),
		"c":    wf.Config.GetBool("c", true),
	}
	globalConfig := map[string]interface{}{
		"ts_formats": mapStrings(strings.Split(wf.Config.Get("ts_formats", "%Y-%M-%D %H-%m-%s %z"), ",")),
	}

	log.Printf("Utility enabled: %v", keyToEnabled)
	log.Printf("Global configuration: %v", globalConfig)

	// Parse input into queries
	var queries []string
	if len(wf.Args()) > 0 {
		rawQuery := strings.TrimSpace(wf.Args()[0])
		queries = parseQueries(rawQuery)
		log.Printf("Parsed queries: %v", queries)
	}

	if len(queries) == 0 {
		for _, utility := range utilities {
			if keyToEnabled[utility.Key()] {
				utility.Do(wf)
			}
		}
	} else {
		utility := createUtility(queries, globalConfig)
		if utility != nil && keyToEnabled[utility.Key()] {
			utility.Do(wf)
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
				results = append(results, strings.Trim(group, `"'`))
			}
		}
	}
	return results
}

func mapStrings(escapedValues []string) []string {
	var holder []string
	for _, value := range escapedValues {
		holder = append(holder, strings.Trim(strings.Trim(value, "'"), "\""))
	}
	return holder
}
