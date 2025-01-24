package json

import (
	"fmt"
	aw "github.com/deanishe/awgo"
	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
	"log"
	"os/exec"
	"strings"
)

type JSON struct {
	value1 string
	value2 string
}

func NewJson(value1, value2 string) *JSON {
	return &JSON{value1: value1, value2: value2}
}

func (j *JSON) Key() string {
	return "json"
}

func (j *JSON) Do(wf *aw.Workflow) {
	log.Printf("[JSON] \tprocessing..., values: %s %s\n", j.value1, j.value2)

	command := "p"
	if "m" == j.value1 {
		command = "m"
	}

	candidate := j.value2
	if "p" != j.value1 && "m" != j.value1 {
		candidate = j.value1
	}

	if "" == candidate {
		tmp, err := exec.Command("pbpaste").Output()
		if err != nil {
			log.Fatalf("error getting clipboard content: %s", err)
		}
		candidate = string(tmp)
	}

	log.Printf("[JSON] \tcommand: %s, candidate: %s\n", command, candidate)

	if !gjson.Valid(candidate) {
		if "" == j.value1 && "" == j.value2 {
			wf.NewItem("json p/m [JSON]").
				Subtitle("Prettify/minify JSON string.").
				Valid(false).Autocomplete(fmt.Sprintf("%s p ", j.Key()))
			return
		}

		wf.NewItem("Invalid JSON string.").
			Valid(false)
		return
	}

	var result string
	switch command {
	case "p":
		result = string(pretty.Pretty([]byte(candidate)))
	case "m":
		result = string(pretty.Ugly([]byte(candidate)))
	}

	wf.NewItem(strings.ReplaceAll(result, "\n", " ")).
		Subtitle("Prettify/minify JSON string.").
		Arg(result).
		Copytext(result).
		Quicklook(result).
		Valid(true).
		Autocomplete(fmt.Sprintf("%s %s ", j.Key(), command))
}
