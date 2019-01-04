package request

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/keyname"
)

type request struct {
	UUID  *uuid.UUID      `json:"id"`
	Frm   user.User       `json:"from"`
	Nw    entity.Entity   `json:"new_entity"`
	Rson  string          `json:"reason"`
	Kname keyname.Keyname `json:"keyname"`
}

func createRequest(id *uuid.UUID, frm user.User, nw entity.Entity, reason string, kname keyname.Keyname) Request {
	out := request{
		UUID:  id,
		Frm:   frm,
		Nw:    nw,
		Rson:  reason,
		Kname: kname,
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

	knameIns, knameInsErr := keyname.SDKFunc.CreateMetaData().Denormalize()(normalized.Keyname)
	if knameInsErr != nil {
		return nil, knameInsErr
	}

	if from, ok := fromIns.(user.User); ok {
		if kname, ok := knameIns.(keyname.Keyname); ok {
			ins, insErr := reg.fromJSONToEntity(normalized.NewEntityJS, kname.Name())
			if insErr != nil {
				return nil, insErr
			}

			out := createRequest(&id, from, ins, normalized.Reason, kname)
			return out, nil
		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid Keyname instance", knameIns.ID().String())
		return nil, errors.New(str)
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

// Reason returns the reason
func (req *request) Reason() string {
	return req.Rson
}

// Keyname returns the keyname
func (req *request) Keyname() keyname.Keyname {
	return req.Kname
}
