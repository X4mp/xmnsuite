package completed

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request"
)

// Request represents a completed request
type Request interface {
	ID() *uuid.UUID
	Request() request.Request
	ConcensusNeeded() int
	Approved() int
	Disapproved() int
	Neutral() int
}
