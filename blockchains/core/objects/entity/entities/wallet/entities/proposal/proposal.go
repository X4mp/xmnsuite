package proposal

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/category"
)

type proposal struct {
	UUID              *uuid.UUID        `json:"id"`
	Titl              string            `json:"title"`
	Desc              string            `json:"description"`
	Det               string            `json:"details"`
	Cat               category.Category `json:"category"`
	MgrPledgeNeeded   int               `json:"manager_pledge_needed"`
	LnkerPledgeNeeded int               `json:"linker_pledge_needed"`
}

func createProposal(id *uuid.UUID, title string, description string, details string, cat category.Category, mgrPledgeNeeded int, linkerPledgeNeeded int) (Proposal, error) {
	out := proposal{
		UUID:              id,
		Titl:              title,
		Desc:              description,
		Det:               details,
		Cat:               cat,
		MgrPledgeNeeded:   mgrPledgeNeeded,
		LnkerPledgeNeeded: linkerPledgeNeeded,
	}

	return &out, nil
}

func createProposalFromNormalized(normalized *normalizedProposal) (Proposal, error) {
	id, idErr := uuid.FromString(normalized.ID)
	if idErr != nil {
		return nil, idErr
	}

	catIns, catInsErr := category.SDKFunc.CreateMetaData().Denormalize()(normalized.Category)
	if catInsErr != nil {
		return nil, catInsErr
	}

	if cat, ok := catIns.(category.Category); ok {
		return createProposal(&id, normalized.Title, normalized.Description, normalized.Details, cat, normalized.ManagerPledgeNeeded, normalized.LinkerPledgeNeeded)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Category instance", catIns.ID().String())
	return nil, errors.New(str)
}

func createProposalFromStorable(storable *storableProposal, rep entity.Repository) (Proposal, error) {
	id, idErr := uuid.FromString(storable.ID)
	if idErr != nil {
		return nil, idErr
	}

	catID, catIDErr := uuid.FromString(storable.CategoryID)
	if catIDErr != nil {
		return nil, catIDErr
	}

	catIns, catInsErr := rep.RetrieveByID(category.SDKFunc.CreateMetaData(), &catID)
	if catInsErr != nil {
		return nil, catInsErr
	}

	if cat, ok := catIns.(category.Category); ok {
		return createProposal(&id, storable.Title, storable.Description, storable.Details, cat, storable.ManagerPledgeNeeded, storable.LinkerPledgeNeeded)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Category instance", catIns.ID().String())
	return nil, errors.New(str)
}

// ID returns the ID
func (obj *proposal) ID() *uuid.UUID {
	return obj.UUID
}

// Title returns the title
func (obj *proposal) Title() string {
	return obj.Titl
}

// Description returns the description
func (obj *proposal) Description() string {
	return obj.Desc
}

// Details returns the details
func (obj *proposal) Details() string {
	return obj.Det
}

// Category returns the category
func (obj *proposal) Category() category.Category {
	return obj.Cat
}

// ManagerPledgeNeeded returns the manager pledge needed
func (obj *proposal) ManagerPledgeNeeded() int {
	return obj.MgrPledgeNeeded
}

// LinkerPledgeNeeded returns the linker pledge needed
func (obj *proposal) LinkerPledgeNeeded() int {
	return obj.LnkerPledgeNeeded
}
