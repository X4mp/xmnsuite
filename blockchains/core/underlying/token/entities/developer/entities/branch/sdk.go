package branch

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/entities/developer/entities/file"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/entities/developer/entities/project"
)

// Branch represents a branch
type Branch interface {
	ID() *uuid.UUID
	Name() string
	Project() project.Project
	Files() []file.File
	Parent() Branch
}
