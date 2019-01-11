package feature

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
)

type normalizedFeature struct {
	ID        string             `json:"id"`
	Project   project.Normalized `json:"project"`
	Title     string             `json:"title"`
	Details   string             `json:"details"`
	CreatedBy user.Normalized    `json:"created_by"`
}

func createNormalizedFeature(ins Feature) (*normalizedFeature, error) {
	proj, projErr := project.SDKFunc.CreateMetaData().Normalize()(ins.Project())
	if projErr != nil {
		return nil, projErr
	}

	usr, usrErr := user.SDKFunc.CreateMetaData().Normalize()(ins.CreatedBy())
	if usrErr != nil {
		return nil, usrErr
	}

	out := normalizedFeature{
		ID:        ins.ID().String(),
		Project:   proj,
		Title:     ins.Title(),
		Details:   ins.Details(),
		CreatedBy: usr,
	}

	return &out, nil
}
