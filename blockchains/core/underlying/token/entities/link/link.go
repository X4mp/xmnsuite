package link

import (
	uuid "github.com/satori/go.uuid"
)

type link struct {
	UUID *uuid.UUID `json:"id"`
	Titl string     `json:"title"`
	Desc string     `json:"description"`
}

func createLink(id *uuid.UUID, title string, description string) (Link, error) {
	out := link{
		UUID: id,
		Titl: title,
		Desc: description,
	}

	return &out, nil
}

func createLinkFromNormalized(normalized *normalizedLink) (Link, error) {
	id, idErr := uuid.FromString(normalized.ID)
	if idErr != nil {
		return nil, idErr
	}

	return createLink(&id, normalized.Title, normalized.Description)
}

// ID returns the ID
func (obj *link) ID() *uuid.UUID {
	return obj.UUID
}

// Title returns the title
func (obj *link) Title() string {
	return obj.Titl
}

// Description returns the description
func (obj *link) Description() string {
	return obj.Desc
}
