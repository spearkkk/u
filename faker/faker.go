package faker

import (
	faker "github.com/brianvoe/gofakeit/v7"
	aw "github.com/deanishe/awgo"
	"log"
)

type Faker struct {
	values []string
}

func NewFaker(values []string) *Faker {
	return &Faker{values: values}
}

func (f *Faker) Key() string {
	return "fk"
}

func (f *Faker) Do(wf *aw.Workflow) {
	wf.Configure(aw.SuppressUIDs(false))

	log.Printf("[Faker] \tprocessing..., values: %+v", f.values)

	result := faker.Noun()

	wf.NewItem(result).
		Arg(result).
		Copytext(result).
		Quicklook(result).
		Valid(true).
		Autocomplete(f.Key())
}
