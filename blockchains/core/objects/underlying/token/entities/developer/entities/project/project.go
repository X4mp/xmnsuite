package project

import uuid "github.com/satori/go.uuid"

type project struct {
	UUID *uuid.UUID `json:"id"`
	Titl string     `json:"title"`
	Desc string     `json:"description"`
}

func createProject(id *uuid.UUID, title string, description string) Project {
	out := project{
		UUID: id,
		Titl: title,
		Desc: description,
	}

	return &out
}

func createProjectFromStorable(storable *storableProject) (Project, error) {
	id, idErr := uuid.FromString(storable.ID)
	if idErr != nil {
		return nil, idErr
	}

	out := createProject(&id, storable.Title, storable.Description)
	return out, nil
}

// ID returns the ID
func (obj *project) ID() *uuid.UUID {
	return obj.UUID
}

// Title returns the title
func (obj *project) Title() string {
	return obj.Titl
}

// Description returns the description
func (obj *project) Description() string {
	return obj.Desc
}
