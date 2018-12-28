package info

import uuid "github.com/satori/go.uuid"

// Info represents the community information
type Info interface {
	ID() *uuid.UUID
	ChunkSizeInBytes() int
}
