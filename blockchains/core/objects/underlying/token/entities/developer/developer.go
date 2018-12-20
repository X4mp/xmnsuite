package developer

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/pledge"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/user"
)

type developer struct {
	UUID  *uuid.UUID    `json:"id"`
	Pldge pledge.Pledge `json:"pledge"`
	Usr   user.User     `json:"user"`
	Nme   string        `json:"name"`
	Res   string        `json:"resume"`
}

func createDeveloper(id *uuid.UUID, usr user.User, pldge pledge.Pledge, name string, resume string) Developer {
	out := developer{
		UUID:  id,
		Usr:   usr,
		Pldge: pldge,
		Nme:   name,
		Res:   resume,
	}

	return &out
}

func createDeveloperFromNormalized(normalized *normalizedDeveloper) (Developer, error) {
	id, idErr := uuid.FromString(normalized.ID)
	if idErr != nil {
		return nil, idErr
	}

	usrIns, usrInsErr := user.SDKFunc.CreateMetaData().Denormalize()(normalized.User)
	if usrInsErr != nil {
		return nil, usrInsErr
	}

	pldgeIns, pldgeInsErr := pledge.SDKFunc.CreateMetaData().Denormalize()(normalized.Pledge)
	if pldgeInsErr != nil {
		return nil, pldgeInsErr
	}

	if usr, ok := usrIns.(user.User); ok {
		if pldge, ok := pldgeIns.(pledge.Pledge); ok {
			out := createDeveloper(&id, usr, pldge, normalized.Name, normalized.Resume)
			return out, nil
		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid Pledge instance", pldgeIns.ID().String())
		return nil, errors.New(str)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid User instance", usrIns.ID().String())
	return nil, errors.New(str)

}

// ID returns the ID
func (obj *developer) ID() *uuid.UUID {
	return obj.UUID
}

// User returns the user
func (obj *developer) User() user.User {
	return obj.Usr
}

// Pledge returns the pledge
func (obj *developer) Pledge() pledge.Pledge {
	return obj.Pldge
}

// Name returns the name
func (obj *developer) Name() string {
	return obj.Nme
}

// Resume returns the resume
func (obj *developer) Resume() string {
	return obj.Res
}
