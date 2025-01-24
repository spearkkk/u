package uuid

import (
	"github.com/deanishe/awgo"
	"github.com/google/uuid"
)

var instance *UUID

type UUID struct {
}

func init() {
	instance = &UUID{}
}

func NewUUID() *UUID {
	return instance
}

func (u *UUID) Key() string {
	return "uuid"
}

func (u *UUID) Do(wf *aw.Workflow) {
	result := uuid.New().String()

	wf.NewItem(result).
		Subtitle("Generate a unique identifier (UUID).").
		Arg(result).
		Copytext(result).
		Quicklook(result).
		Valid(true).
		Autocomplete(u.Key())
}
