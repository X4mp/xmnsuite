package keyname

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/group"
)

type keyname struct {
	UUID *uuid.UUID  `json:"id"`
	Grp  group.Group `json:"group"`
	Nme  string      `json:"name"`
}

func createKeyname(id *uuid.UUID, grp group.Group, name string) (Keyname, error) {
	out := keyname{
		UUID: id,
		Grp:  grp,
		Nme:  name,
	}

	return &out, nil
}

func createKeynameFromStorable(rep entity.Repository, storable *storableKeyname) (Keyname, error) {
	id, idErr := uuid.FromString(storable.ID)
	if idErr != nil {
		str := fmt.Sprintf("the storable ID (%s) is invalid: %s", storable.ID, idErr.Error())
		return nil, errors.New(str)
	}

	groupID, groupIDErr := uuid.FromString(storable.GroupID)
	if groupIDErr != nil {
		str := fmt.Sprintf("the storable GroupID (%s) is invalid: %s", storable.GroupID, groupIDErr.Error())
		return nil, errors.New(str)
	}

	// retrieve the group:
	metaData := group.SDKFunc.CreateMetaData()
	insGroup, insGroupErr := rep.RetrieveByID(metaData, &groupID)
	if insGroupErr != nil {
		return nil, insGroupErr
	}

	if grp, ok := insGroup.(group.Group); ok {
		return createKeyname(&id, grp, storable.Name)
	}

	str := fmt.Sprintf("the given entity (ID: %s) is not a valid Group instance", groupID.String())
	return nil, errors.New(str)

}

func createKeynameFromNormalized(normalized *normalizedKeyname) (Keyname, error) {
	id, idErr := uuid.FromString(normalized.ID)
	if idErr != nil {
		return nil, idErr
	}

	grpIns, grpInsErr := group.SDKFunc.CreateMetaData().Denormalize()(normalized.Group)
	if grpInsErr != nil {
		return nil, grpInsErr
	}

	if grp, ok := grpIns.(group.Group); ok {
		return createKeyname(&id, grp, normalized.Name)
	}

	str := fmt.Sprintf("the given entity (ID: %s) is not a valid Keyname instance", grpIns.ID().String())
	return nil, errors.New(str)
}

// ID returns the ID
func (obj *keyname) ID() *uuid.UUID {
	return obj.UUID
}

// Group returns the group
func (obj *keyname) Group() group.Group {
	return obj.Grp
}

// Name returns the name
func (obj *keyname) Name() string {
	return obj.Nme
}
