package main

type Option struct {
    name        string
    description string
}

type Command struct {
    name           string
    description    string
    options        []Option
    defaultOption  *Option
    requiresOption bool
}

type Function struct {
    name            string
    description     string
    commands        []Command
    defaultCommand  *Command
    requiresCommand bool
}
