package milestone

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project/feature"
)

type milestone struct {
	UUID  *uuid.UUID      `json:"id"`
	Proj  project.Project `json:"project"`
	Feat  feature.Feature `json:"feature"`
	Wal   wallet.Wallet   `json:"wallet"`
	Shres int             `json:"shares"`
	Titl  string          `json:"title"`
	Desc  string          `json:"description"`
	Det   string          `json:"details"`
}

func createMilestone(id *uuid.UUID, proj project.Project, wal wallet.Wallet, shares int, title string, description string, details string) (Milestone, error) {
	return createMilestoneWithFeature(id, proj, wal, shares, title, description, details, nil)
}

func createMilestoneWithFeature(id *uuid.UUID, proj project.Project, wal wallet.Wallet, shares int, title string, description string, details string, feat feature.Feature) (Milestone, error) {
	out := milestone{
		UUID:  id,
		Proj:  proj,
		Feat:  feat,
		Wal:   wal,
		Shres: shares,
		Titl:  title,
		Desc:  description,
		Det:   details,
	}

	return &out, nil
}

func createMilestoneFromNormalized(normalized *normalizedMilestone) (Milestone, error) {
	id, idErr := uuid.FromString(normalized.ID)
	if idErr != nil {
		return nil, idErr
	}

	projIns, projInsErr := project.SDKFunc.CreateMetaData().Denormalize()(normalized.Project)
	if projInsErr != nil {
		return nil, projInsErr
	}

	walIns, walInsErr := wallet.SDKFunc.CreateMetaData().Denormalize()(normalized.Wallet)
	if walInsErr != nil {
		return nil, walInsErr
	}

	var featIns entity.Entity
	if normalized.Feature != nil {
		featInsDenorm, featInsDenormErr := feature.SDKFunc.CreateMetaData().Denormalize()(normalized.Feature)
		if featInsDenormErr != nil {
			return nil, featInsDenormErr
		}

		featIns = featInsDenorm
	}

	if proj, ok := projIns.(project.Project); ok {
		if wal, ok := walIns.(wallet.Wallet); ok {

			if featIns != nil {
				if feat, ok := featIns.(feature.Feature); ok {
					return createMilestoneWithFeature(&id, proj, wal, normalized.Shares, normalized.Title, normalized.Description, normalized.Details, feat)
				}

				str := fmt.Sprintf("the entity (ID: %s) is not a valid Feature instance", featIns.ID().String())
				return nil, errors.New(str)
			}

			return createMilestone(&id, proj, wal, normalized.Shares, normalized.Title, normalized.Description, normalized.Details)

		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid Wallet instance", walIns.ID().String())
		return nil, errors.New(str)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Project instance", projIns.ID().String())
	return nil, errors.New(str)
}

func createMilestoneFromStorable(storable *storableMilestone, rep entity.Repository) (Milestone, error) {
	id, idErr := uuid.FromString(storable.ID)
	if idErr != nil {
		return nil, idErr
	}

	projID, projIDErr := uuid.FromString(storable.ProjectID)
	if projIDErr != nil {
		return nil, projIDErr
	}

	walletID, walletIDErr := uuid.FromString(storable.WalletID)
	if walletIDErr != nil {
		return nil, walletIDErr
	}

	projIns, projInsErr := rep.RetrieveByID(project.SDKFunc.CreateMetaData(), &projID)
	if projInsErr != nil {
		return nil, projInsErr
	}

	walIns, walInsErr := rep.RetrieveByID(wallet.SDKFunc.CreateMetaData(), &walletID)
	if walInsErr != nil {
		return nil, walInsErr
	}

	var featIns entity.Entity
	if storable.FeatureID != "" {
		featID, featIDErr := uuid.FromString(storable.FeatureID)
		if featIDErr != nil {
			return nil, featIDErr
		}

		retFeatIns, retFeatInsErr := rep.RetrieveByID(feature.SDKFunc.CreateMetaData(), &featID)
		if retFeatInsErr != nil {
			return nil, retFeatInsErr
		}

		featIns = retFeatIns
	}

	if proj, ok := projIns.(project.Project); ok {
		if wal, ok := walIns.(wallet.Wallet); ok {
			if featIns != nil {
				if feat, ok := featIns.(feature.Feature); ok {
					return createMilestoneWithFeature(&id, proj, wal, storable.Shares, storable.Title, storable.Description, storable.Details, feat)
				}

				str := fmt.Sprintf("the entity (ID: %s) is not a valid Feature instance", featIns.ID().String())
				return nil, errors.New(str)
			}

			return createMilestone(&id, proj, wal, storable.Shares, storable.Title, storable.Description, storable.Details)
		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid Wallet instance", walIns.ID().String())
		return nil, errors.New(str)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Project instance", projIns.ID().String())
	return nil, errors.New(str)
}

// ID returns the ID
func (obj *milestone) ID() *uuid.UUID {
	return obj.UUID
}

// Project returns the project
func (obj *milestone) Project() project.Project {
	return obj.Proj
}

// Wallet returns the wallet
func (obj *milestone) Wallet() wallet.Wallet {
	return obj.Wal
}

// Shares returns the shares
func (obj *milestone) Shares() int {
	return obj.Shres
}

// Title returns the title
func (obj *milestone) Title() string {
	return obj.Titl
}

// Description returns the description
func (obj *milestone) Description() string {
	return obj.Desc
}

// Details returns the details
func (obj *milestone) Details() string {
	return obj.Det
}

// HasFeature returns true if there is a feature, false otherwise
func (obj *milestone) HasFeature() bool {
	return obj.Feat != nil
}

// Feature returns the feature
func (obj *milestone) Feature() feature.Feature {
	return obj.Feat
}
