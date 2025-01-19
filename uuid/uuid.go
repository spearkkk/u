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

func (u *UUID) GetKey() string {
	return "uuid"
}

func (u *UUID) GetName() string {
	return "UUID"
}

func (u *UUID) GetDescription() string {
	return "[copy:⏎] " + u.GetName() + ": To generate a unique identifier (UUID)."
}

func (u *UUID) GetResults(wf *aw.Workflow) {
	result := uuid.New().String()
	wf.NewItem(result).
		Subtitle(u.GetDescription()).
		Arg(result).
		Copytext(result).
		Quicklook(result).
		Valid(true).
		Autocomplete(u.GetKey())
}
