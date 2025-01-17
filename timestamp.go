package main

var Timestamp = Function{
	name:        "ts",
	description: "Convert between timestamps and human readable dates",
	commands: []Command{
		{name: "unix(epoch)", description: "Get the current unix timestamp", requiresOption: false},
		{name: "unix(milliseconds)", description: "Get the current unix timestamp with milliseconds", requiresOption: false},
		{name: "fmt", description: "Get the current timestamp as formatted string",
			requiresOption: true, options: []Option{
				{name: "YYYY-MM-DD", description: "2025-01-01"},
				{name: "YYYY-MM-DDThh:mm:ssZ", description: "2025-01-01T12:00:00Z"},
				{name: "YYYY-MM-DDThh:mm:ss.SSSZ", description: "2025-01-01T12:00:00.000Z"},
			},
			defaultOption: &Option{name: "iso8601", description: "ISO 8601 format"}},
		{name: "cvt", description: "Convert between millisecond timestamps and human readable dates",
			requiresOption: true, options: []Option{
				{name: "milliseconds", description: "Convert from milliseconds to human readable date"},
				{name: "unix", description: "Convert from human readable date to milliseconds"},
			},
			defaultOption: &Option{name: "milliseconds", description: "Convert from milliseconds to human readable date"}},
	},
	defaultCommand:  &Command{name: "unix(epoch)", description: "Get the current unix timestamp", requiresOption: false},
	requiresCommand: true,
}
