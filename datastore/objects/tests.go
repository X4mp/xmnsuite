package objects

import (
	uuid "github.com/satori/go.uuid"
)

type objFortests struct {
	ID          *uuid.UUID `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
}

func createObjForTests() *objFortests {
	id := uuid.NewV4()
	out := objFortests{
		ID:          &id,
		Name:        "My name",
		Description: "This is a simple description",
	}

	return &out
}
