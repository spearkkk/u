package main

import (
	aw "github.com/deanishe/awgo"
	c "github.com/spearkkk/u/color"
	json "github.com/spearkkk/u/json"
	"github.com/spearkkk/u/strcase"
	ts "github.com/spearkkk/u/timestamp"
	"github.com/spearkkk/u/uuid"
)

type Utility interface {
	Key() string
	Do(wf *aw.Workflow)
}

func createUtility(queries []string, config map[string]interface{}, queryWithoutKey string) Utility {
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

	switch key {
	case "uuid":
		return uuid.NewUUID()
	case "ts":
		tsFormats, ok := config["ts_formats"].([]string)
		if !ok {
			tsFormats = []string{"'%Y-%M-%D %H-%m-%s %z'"}
		}
		return ts.NewTimestamp(value1, value2, tsFormats...)
	case "json":
		return json.NewJson(value1, value2)
	case "color":
		return c.NewColor(queries[1:])
	case "case":
		return strcase.NewStrCase(queryWithoutKey)
	}

	return nil
}
