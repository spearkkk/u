package main

import (
	"strings"

	aw "github.com/deanishe/awgo"
)

var wf *aw.Workflow

func init() {
	wf = aw.New()
}

func main() {
	wf.Run(run)
}

func run() {
	// Parse input into queries
	var queries []string
	if len(wf.Args()) > 0 {
		rawQuery := strings.TrimSpace(wf.Args()[0])
		queries = parseQueries(rawQuery)
	}

	functions := getFunctions()

	if len(queries) == 0 {
		// Show all functions with their default command results
		for _, function := range functions {
			if function.defaultCommand != nil {
				wf.NewItem(function.name).
					Subtitle(function.defaultCommand.description).
					Arg(function.defaultCommand.name).
					Valid(true)
			}
		}
	} else {
		// Show matched functions based on the query
		for _, function := range functions {
			if matchesQuery(function, queries) {
				wf.NewItem(function.name).
					Subtitle(function.description).
					Arg(function.name).
					Valid(true)
			}
		}
	}

	wf.SendFeedback()
}

func parseQueries(input string) []string {
	return strings.Fields(input)
}

func matchesQuery(function Function, queries []string) bool {
	for _, query := range queries {
		if strings.Contains(function.name, query) || strings.Contains(function.description, query) {
			return true
		}
		for _, command := range function.commands {
			if strings.Contains(command.name, query) || strings.Contains(command.description, query) {
				return true
			}
			for _, option := range command.options {
				if strings.Contains(option.name, query) || strings.Contains(option.description, query) {
					return true
				}
			}
		}
	}
	return false
}

func getFunctions() []Function {
	// Define your functions, commands, and options here
	return []Function{
		{
			name:        "Function1",
			description: "Description of Function1",
			commands: []Command{
				{
					name:        "Command1",
					description: "Description of Command1",
					options: []Option{
						{name: "Option1", description: "Description of Option1"},
					},
					defaultOption:  &Option{name: "Option1", description: "Description of Option1"},
					requiresOption: true,
				},
			},
			defaultCommand:  &Command{name: "Command1", description: "Description of Command1"},
			requiresCommand: true,
		},
		// Add more functions as needed
	}
}
