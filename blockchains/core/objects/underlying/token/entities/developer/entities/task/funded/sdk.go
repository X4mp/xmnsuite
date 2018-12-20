package funded

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/pledge"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/developer/entities/task"
)

// Task represents a funded task
type Task interface {
	ID() *uuid.UUID
	Task() task.Task
	Funds() []pledge.Pledge
	CreatedOn() time.Time
}
