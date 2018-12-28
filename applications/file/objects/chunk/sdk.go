package chunk

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/applications/file/objects/file"
	"github.com/xmnservices/xmnsuite/hashtree"
)

// Chunk represents a file chunk
type Chunk interface {
	ID() *uuid.UUID
	Hash() hashtree.Hash
	File() file.File
	Size() int
	Path() string
}
