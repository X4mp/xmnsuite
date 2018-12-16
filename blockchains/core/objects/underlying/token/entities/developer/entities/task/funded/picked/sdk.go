package picked

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/pledge"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/developer"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/developer/entities/task/funded"
)

// Task represents a task picked by a developer
type Task interface {
	ID() *uuid.UUID
	Task() funded.Task
	Pledge() pledge.Pledge
	Developer() developer.Developer
	CreatedOn() time.Time
}
