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
	UUID         *uuid.UUID      `json:"id"`
	Frm          user.User       `json:"from"`
	SaveEntity   entity.Entity   `json:"new_entity"`
	DeleteEntity entity.Entity   `json:"delete_entity"`
	Rson         string          `json:"reason"`
	Kname        keyname.Keyname `json:"keyname"`
}

func createRequestWithSaveEntity(id *uuid.UUID, frm user.User, saveEntity entity.Entity, reason string, kname keyname.Keyname) Request {
	out := request{
		UUID:         id,
		Frm:          frm,
		SaveEntity:   saveEntity,
		DeleteEntity: nil,
		Rson:         reason,
		Kname:        kname,
	}

	return &out
}

func createRequestWithDeleteEntity(id *uuid.UUID, frm user.User, deleteEntity entity.Entity, reason string, kname keyname.Keyname) Request {
	out := request{
		UUID:         id,
		Frm:          frm,
		SaveEntity:   nil,
		DeleteEntity: deleteEntity,
		Rson:         reason,
		Kname:        kname,
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

			if normalized.SaveEntityJSON != nil {
				saveIns, saveInsErr := reg.fromJSONToEntity(normalized.SaveEntityJSON, kname.Name())
				if saveInsErr != nil {
					return nil, saveInsErr
				}

				out := createRequestWithSaveEntity(&id, from, saveIns, normalized.Reason, kname)
				return out, nil
			}

			delIns, delInsErr := reg.fromJSONToEntity(normalized.DeleteEntityJSON, kname.Name())
			if delInsErr != nil {
				return nil, delInsErr
			}

			out := createRequestWithDeleteEntity(&id, from, delIns, normalized.Reason, kname)
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

// HasSave returns true if there is a save entity instance
func (req *request) HasSave() bool {
	return req.SaveEntity != nil
}

// Save returns the save entity instance
func (req *request) Save() entity.Entity {
	return req.SaveEntity
}

// HasDelete returns true if there is a delete entity instance
func (req *request) HasDelete() bool {
	return req.DeleteEntity != nil
}

// Delete returns the delete entity instance
func (req *request) Delete() entity.Entity {
	return req.DeleteEntity
}

// Reason returns the reason
func (req *request) Reason() string {
	return req.Rson
}

// Keyname returns the keyname
func (req *request) Keyname() keyname.Keyname {
	return req.Kname
}
