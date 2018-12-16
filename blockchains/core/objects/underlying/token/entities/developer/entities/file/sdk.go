package file

import uuid "github.com/satori/go.uuid"

// File represents a file
type File interface {
	ID() *uuid.UUID
	Path() string
	Content() string
}
