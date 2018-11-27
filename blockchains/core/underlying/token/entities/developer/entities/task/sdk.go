package task

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/entities/developer"
)

// Task represents a task
type Task interface {
	ID() *uuid.UUID
	Creator() developer.Developer
	AssignTo() developer.Developer
	Title() string
	Description() string
}
