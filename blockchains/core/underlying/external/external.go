package external

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/request/entities/link"
)

type external struct {
	UUID  *uuid.UUID `json:"id"`
	Lnk   link.Link  `json:"link"`
	ResID *uuid.UUID `json:"resource_id"`
}

func createExternal(id *uuid.UUID, lnk link.Link, resID *uuid.UUID) External {
	out := external{
		UUID:  id,
		Lnk:   lnk,
		ResID: resID,
	}

	return &out
}

// ID returns the ID
func (obj *external) ID() *uuid.UUID {
	return obj.UUID
}

// Link returns the link
func (obj *external) Link() link.Link {
	return obj.Lnk
}

// ResourceID returns the resource ID
func (obj *external) ResourceID() *uuid.UUID {
	return obj.ResID
}
