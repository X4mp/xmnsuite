package category

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
)

type category struct {
	UUID *uuid.UUID `json:"id"`
	Par  Category   `json:"parent"`
	Nme  string     `json:"name"`
	Desc string     `json:"description"`
}

func createCategory(id *uuid.UUID, name string, description string) (Category, error) {

	if len(name) > maxAountOfCharactersForName {
		str := fmt.Sprintf("the name (%s) contains %d characters, the limit is: %d", name, len(name), maxAountOfCharactersForName)
		return nil, errors.New(str)
	}

	if len(description) > maxAmountOfCharactersForDescription {
		str := fmt.Sprintf("the description (%s) contains %d characters, thelimit is: %d", description, len(description), maxAmountOfCharactersForDescription)
		return nil, errors.New(str)
	}

	out := category{
		UUID: id,
		Par:  nil,
		Nme:  name,
		Desc: description,
	}

	return &out, nil
}

func createCategoryWithParent(id *uuid.UUID, parent Category, name string, description string) (Category, error) {
	out := category{
		UUID: id,
		Par:  parent,
		Nme:  name,
		Desc: description,
	}

	return &out, nil
}

func fromNormalizedToCategory(normalized *normalizedCategory) (Category, error) {
	id, idErr := uuid.FromString(normalized.ID)
	if idErr != nil {
		return nil, idErr
	}

	if normalized.Parent != nil {
		metaData := createMetaData()
		parIns, parInsErr := metaData.Denormalize()(normalized.Parent)
		if parInsErr != nil {
			return nil, parInsErr
		}

		if par, ok := parIns.(Category); ok {
			return createCategoryWithParent(&id, par, normalized.Name, normalized.Description)
		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid Category instance", parIns.ID().String())
		return nil, errors.New(str)
	}

	return createCategory(&id, normalized.Name, normalized.Description)
}

func fromStorableToCategory(storable *storableCategory, rep entity.Repository) (Category, error) {
	id, idErr := uuid.FromString(storable.ID)
	if idErr != nil {
		return nil, idErr
	}

	if storable.ParentID != "" {
		parentID, parentIDErr := uuid.FromString(storable.ParentID)
		if parentIDErr != nil {
			return nil, parentIDErr
		}

		metaData := createMetaData()
		parIns, parInsErr := rep.RetrieveByID(metaData, &parentID)
		if parInsErr != nil {
			return nil, parInsErr
		}

		if par, ok := parIns.(Category); ok {
			return createCategoryWithParent(&id, par, storable.Name, storable.Description)
		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid Category instance", parIns.ID().String())
		return nil, errors.New(str)
	}

	return createCategory(&id, storable.Name, storable.Description)
}

// ID returns the ID
func (obj *category) ID() *uuid.UUID {
	return obj.UUID
}

// HasParent returns true if there is a parent category, false otherwise
func (obj *category) HasParent() bool {
	return obj.Par != nil
}

// Parent returns the parent category, if any
func (obj *category) Parent() Category {
	return obj.Par
}

// Name returns the name
func (obj *category) Name() string {
	return obj.Nme
}

// Description returns the description
func (obj *category) Description() string {
	return obj.Desc
}
