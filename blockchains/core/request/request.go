package request

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/user"
)

type request struct {
	UUID *uuid.UUID    `json:"id"`
	Frm  user.User     `json:"from"`
	Nw   entity.Entity `json:"new"`
}

func createRequest(id *uuid.UUID, frm user.User, nw entity.Entity) Request {
	out := request{
		UUID: id,
		Frm:  frm,
		Nw:   nw,
	}

	return &out
}

// ID returns the ID
func (req *request) ID() *uuid.UUID {
	return req.UUID
}

// From returns the from user
func (req *request) From() user.User {
	return req.Frm
}

// New returns the new entity to be created
func (req *request) New() entity.Entity {
	return req.Nw
}
