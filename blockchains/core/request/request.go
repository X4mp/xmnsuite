package request

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/entities/user"
)

type request struct {
	UUID   *uuid.UUID    `json:"id"`
	Frm    user.User     `json:"from"`
	Nw     entity.Entity `json:"new_entity"`
	NwName string        `json:"new_entity_name"`
}

func createRequest(id *uuid.UUID, frm user.User, nw entity.Entity, newName string) Request {
	out := request{
		UUID:   id,
		Frm:    frm,
		Nw:     nw,
		NwName: newName,
	}

	return &out
}

func createRequestFromNormalized(normalized *normalizedRequest) (Request, error) {
	id, idErr := uuid.FromString(normalized.ID)
	if idErr != nil {
		str := fmt.Sprintf("the normalized ID (%s) is invalid: %s", normalized.ID, idErr.Error())
		return nil, errors.New(str)
	}

	fromIns, fromInsErr := user.SDKFunc.CreateMetaData().Denormalize()(normalized.From)
	if fromInsErr != nil {
		return nil, fromInsErr
	}

	ins, insErr := reg.fromJSONToEntity(normalized.NewEntityJS, normalized.NewEntityName)
	if insErr != nil {
		return nil, insErr
	}

	if from, ok := fromIns.(user.User); ok {
		out := createRequest(&id, from, ins, normalized.NewEntityName)
		return out, nil
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid User instance", fromIns.ID().String())
	return nil, errors.New(str)
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

// NewName returns the new entity name
func (req *request) NewName() string {
	return req.NwName
}
