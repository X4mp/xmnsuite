package group

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type group struct {
	UUID *uuid.UUID `json:"id"`
	Nme  string     `json:"name"`
}

func createGroup(id *uuid.UUID, name string) (Group, error) {
	out := group{
		UUID: id,
		Nme:  name,
	}

	return &out, nil
}

func createGroupFromStorable(storable *storableGroup) (Group, error) {
	id, idErr := uuid.FromString(storable.ID)
	if idErr != nil {
		str := fmt.Sprintf("the storable ID (%s) is invalid: %s", storable.ID, idErr.Error())
		return nil, errors.New(str)
	}

	return createGroup(&id, storable.Name)
}

// ID returns the ID
func (obj *group) ID() *uuid.UUID {
	return obj.UUID
}

// Name returns the name
func (obj *group) Name() string {
	return obj.Nme
}
