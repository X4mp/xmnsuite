package file

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/hashtree"
)

// File represents a file
type File interface {
	ID() *uuid.UUID
	HashTree() hashtree.HashTree
}
