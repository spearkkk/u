package uuid

import (
	"github.com/google/uuid"
	"github.com/deanishe/awgo"
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

func (u *UUID) GetKey() string {
	return "uuid"
}

func (u *UUID) GetName() string {
	return "UUID"
}

func (u *UUID) GetDescription() string {
	return "[copy:⏎] To generate a unique identifier (UUID)."
}

func (u *UUID) GetResults(wf *aw.Workflow) {
	result := uuid.New().String()
	wf.NewItem(u.GetName() + ": " + result).
					Subtitle(u.GetDescription()).
					Arg(result).
					Copytext(result).
					Quicklook(result).
					Valid(true).
					Autocomplete(u.GetKey())
}
