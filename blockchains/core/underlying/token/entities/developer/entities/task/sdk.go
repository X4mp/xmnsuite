package task

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/entities/developer"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/entities/developer/entities/milestone"
)

// Task represents a task
type Task interface {
	ID() *uuid.UUID
	Milestone() milestone.Milestone
	Creator() developer.Developer
	Title() string
	Description() string
	CreatedOn() time.Time
	DueOn() time.Time
}
