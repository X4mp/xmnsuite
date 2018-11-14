package category

import uuid "github.com/satori/go.uuid"

type category struct {
	UUID *uuid.UUID `json:"id"`
	Nme  string     `json:"name"`
	Desc string     `json:"description"`
}

func createCategory(id *uuid.UUID, name string, description string) Category {
	out := category{
		UUID: id,
		Nme:  name,
		Desc: description,
	}

	return &out
}

func createCategoryFromStorable(storable *storableCategory) (Category, error) {
	id, idErr := uuid.FromString(storable.ID)
	if idErr != nil {
		return nil, idErr
	}

	out := createCategory(&id, storable.Name, storable.Description)
	return out, nil
}

// ID returns the ID
func (obj *category) ID() *uuid.UUID {
	return obj.UUID
}

// Name returns the name
func (obj *category) Name() string {
	return obj.Nme
}

// Description returns the description
func (obj *category) Description() string {
	return obj.Desc
}
