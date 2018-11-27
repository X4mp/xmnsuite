package project

import uuid "github.com/satori/go.uuid"

// Project represents a project
type Project interface {
	ID() *uuid.UUID
	Title() string
	Description() string
}
