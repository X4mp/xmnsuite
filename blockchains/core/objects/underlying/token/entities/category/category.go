package category

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
)

type category struct {
	UUID *uuid.UUID `json:"id"`
	Titl string     `json:"title"`
	Desc string     `json:"description"`
	Par  Category   `json:"parent_category"`
}

func createCategory(id *uuid.UUID, title string, description string) (Category, error) {
	return createCategoryWithParent(id, title, description, nil)
}

func createCategoryWithParent(id *uuid.UUID, title string, description string, par Category) (Category, error) {
	out := category{
		UUID: id,
		Titl: title,
		Desc: description,
		Par:  par,
	}

	return &out, nil
}

func createCategoryFromNormalized(normalized *normalizedCategory) (Category, error) {
	id, idErr := uuid.FromString(normalized.ID)
	if idErr != nil {
		return nil, idErr
	}

	if normalized.Parent == nil {
		return createCategory(&id, normalized.Title, normalized.Description)
	}

	parIns, parInsErr := createMetaData().Denormalize()(normalized.Parent)
	if parInsErr != nil {
		return nil, parInsErr
	}

	if par, ok := parIns.(Category); ok {
		return createCategoryWithParent(&id, normalized.Title, normalized.Description, par)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Category instance", parIns.ID().String())
	return nil, errors.New(str)
}

func createCategoryFromStorable(storable *storableCategory, rep entity.Repository) (Category, error) {
	id, idErr := uuid.FromString(storable.ID)
	if idErr != nil {
		return nil, idErr
	}

	if storable.ParentID == "" {
		return createCategory(&id, storable.Title, storable.Description)
	}

	parID, parIDErr := uuid.FromString(storable.ParentID)
	if parIDErr != nil {
		return nil, parIDErr
	}

	parIns, parInsErr := rep.RetrieveByID(createMetaData(), &parID)
	if parInsErr != nil {
		return nil, parInsErr
	}

	if par, ok := parIns.(Category); ok {
		return createCategoryWithParent(&id, storable.Title, storable.Description, par)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Category instance", parIns.ID().String())
	return nil, errors.New(str)
}

// ID returns the ID
func (obj *category) ID() *uuid.UUID {
	return obj.UUID
}

// Title returns the title
func (obj *category) Title() string {
	return obj.Titl
}

// Description returns the description
func (obj *category) Description() string {
	return obj.Desc
}

// HasParent returns true if there is a parent category, false otherwise
func (obj *category) HasParent() bool {
	return obj.Par != nil
}

// Parent returns the parent category
func (obj *category) Parent() Category {
	return obj.Par
}
