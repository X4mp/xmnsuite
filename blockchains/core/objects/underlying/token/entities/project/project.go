package project

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	proposal "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal"
)

type project struct {
	UUID *uuid.UUID        `json:"id"`
	Prop proposal.Proposal `json:"proposal"`
}

func createProject(
	id *uuid.UUID,
	prop proposal.Proposal,
) (Project, error) {
	out := project{
		UUID: id,
		Prop: prop,
	}

	return &out, nil
}

func createProjectFromNormalized(normalized *normalizedProject) (Project, error) {
	id, idErr := uuid.FromString(normalized.ID)
	if idErr != nil {
		return nil, idErr
	}

	propIns, propInsErr := proposal.SDKFunc.CreateMetaData().Denormalize()(normalized.Proposal)
	if propInsErr != nil {
		return nil, propInsErr
	}

	if prop, ok := propIns.(proposal.Proposal); ok {
		return createProject(&id, prop)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Proposal instance", propIns.ID().String())
	return nil, errors.New(str)
}

func createProjectFromStorable(storable *storableProject, rep entity.Repository) (Project, error) {
	id, idErr := uuid.FromString(storable.ID)
	if idErr != nil {
		return nil, idErr
	}

	propID, propIDErr := uuid.FromString(storable.ProposalID)
	if propIDErr != nil {
		return nil, propIDErr
	}

	propIns, propInsErr := rep.RetrieveByID(proposal.SDKFunc.CreateMetaData(), &propID)
	if propInsErr != nil {
		return nil, propInsErr
	}

	if prop, ok := propIns.(proposal.Proposal); ok {
		return createProject(&id, prop)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Proposal instance", propIns.ID().String())
	return nil, errors.New(str)
}

// ID returns the ID
func (obj *project) ID() *uuid.UUID {
	return obj.UUID
}

// Proposal returns the proposal
func (obj *project) Proposal() proposal.Proposal {
	return obj.Prop
}
