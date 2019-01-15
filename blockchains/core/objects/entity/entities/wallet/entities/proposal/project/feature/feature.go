package feature

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
)

type feature struct {
	UUID *uuid.UUID      `json:"id"`
	Proj project.Project `json:"project"`
	Titl string          `json:"title"`
	Det  string          `json:"details"`
	CrBy user.User       `json:"created_by"`
}

func createFeature(id *uuid.UUID, proj project.Project, title string, details string, createdBy user.User) (Feature, error) {
	out := feature{
		UUID: id,
		Proj: proj,
		Titl: title,
		Det:  details,
		CrBy: createdBy,
	}

	return &out, nil
}

func createFeatureFromNormalized(normalized *normalizedFeature) (Feature, error) {
	id, idErr := uuid.FromString(normalized.ID)
	if idErr != nil {
		return nil, idErr
	}

	projIns, projInsErr := project.SDKFunc.CreateMetaData().Denormalize()(normalized.Project)
	if projInsErr != nil {
		return nil, projInsErr
	}

	usrIns, usrInsErr := user.SDKFunc.CreateMetaData().Denormalize()(normalized.CreatedBy)
	if usrInsErr != nil {
		return nil, usrInsErr
	}

	if proj, ok := projIns.(project.Project); ok {
		if usr, ok := usrIns.(user.User); ok {
			return createFeature(&id, proj, normalized.Title, normalized.Details, usr)
		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid User instance", usrIns.ID().String())
		return nil, errors.New(str)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Project instance", projIns.ID().String())
	return nil, errors.New(str)
}

func createFeatureFromStorable(storable *storableFeature, rep entity.Repository) (Feature, error) {
	id, idErr := uuid.FromString(storable.ID)
	if idErr != nil {
		return nil, idErr
	}

	projID, projIDErr := uuid.FromString(storable.ProjectID)
	if projIDErr != nil {
		return nil, projIDErr
	}

	usrID, usrIDErr := uuid.FromString(storable.CreatedByUserID)
	if usrIDErr != nil {
		return nil, usrIDErr
	}

	projIns, projInsErr := rep.RetrieveByID(project.SDKFunc.CreateMetaData(), &projID)
	if projInsErr != nil {
		return nil, projInsErr
	}

	usrIns, usrInsErr := rep.RetrieveByID(user.SDKFunc.CreateMetaData(), &usrID)
	if usrInsErr != nil {
		return nil, usrInsErr
	}

	if proj, ok := projIns.(project.Project); ok {
		if usr, ok := usrIns.(user.User); ok {
			return createFeature(&id, proj, storable.Title, storable.Details, usr)
		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid User instance", usrIns.ID().String())
		return nil, errors.New(str)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Project instance", projIns.ID().String())
	return nil, errors.New(str)
}

// ID returns the ID
func (obj *feature) ID() *uuid.UUID {
	return obj.UUID
}

// Project returns the project
func (obj *feature) Project() project.Project {
	return obj.Proj
}

// Title returns the title
func (obj *feature) Title() string {
	return obj.Titl
}

// Details returns the details
func (obj *feature) Details() string {
	return obj.Det
}

// CreatedBy returns the user that created the feature
func (obj *feature) CreatedBy() user.User {
	return obj.CrBy
}
