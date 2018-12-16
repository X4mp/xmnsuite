package version

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/developer/entities/branch"
)

// Version represents a version
type Version interface {
	ID() *uuid.UUID
	From() branch.Branch
	Name() (int, int, int, int, int, int)
	NameAsString() string
}
