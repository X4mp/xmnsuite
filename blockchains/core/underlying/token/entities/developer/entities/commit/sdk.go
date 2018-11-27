package commit

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/entities/developer/entities/branch"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/entities/developer/entities/file"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/entities/developer/entities/task"
)

// Commit represents a commit
type Commit interface {
	ID() *uuid.UUID
	To() branch.Branch
	Task() task.Task
	Files() []file.File
	Description() string
}
