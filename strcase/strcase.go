package strcase

import (
	"fmt"
	aw "github.com/deanishe/awgo"
	str "github.com/iancoleman/strcase"
	"log"
)

type StrCase struct {
	value string
}

func NewStrCase(value string) *StrCase {
	return &StrCase{value: value}
}

func (s *StrCase) Key() string {
	return "case"
}

func (s *StrCase) Do(wf *aw.Workflow) {
	// Prevent to order by UID
	wf.Configure(aw.SuppressUIDs(true))

	candidate := s.value
	if "" == candidate {
		return
	}

	log.Printf("[Case] \tprocessing..., values: %s", candidate)

	example := "muzik tiger"
	values := make(map[string]string)
	values[fmt.Sprintf("snake => %s", str.ToSnake(example))] = str.ToSnake(candidate)
	values[fmt.Sprintf("SCREAMING_SNAKE => %s", str.ToScreamingSnake(example))] = str.ToScreamingSnake(candidate)
	values[fmt.Sprintf("lowerCamel => %s", str.ToLowerCamel(example))] = str.ToLowerCamel(candidate)
	values[fmt.Sprintf("UpperCamel => %s", str.ToCamel(example))] = str.ToCamel(candidate)
	values[fmt.Sprintf("kebab => %s", str.ToKebab(example))] = str.ToKebab(candidate)
	values[fmt.Sprintf("SCREAMING-KEBAB => %s", str.ToScreamingKebab(example))] = str.ToScreamingKebab(candidate)

	for k, v := range values {
		wf.NewItem(v).
			Subtitle(k).
			Arg(v).
			Copytext(v).
			Quicklook(v).
			Valid(true)
	}
}
