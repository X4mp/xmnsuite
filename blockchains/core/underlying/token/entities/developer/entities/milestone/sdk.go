package milestone

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/entities/developer/entities/project"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/entities/developer/entities/task"
)

// Milestone represents a milestone
type Milestone interface {
	ID() *uuid.UUID
	Project() project.Project
	Title() string
	Description() string
	Tasks() []task.Task
	CreatedOn() time.Time
	DueDate() time.Time
}
