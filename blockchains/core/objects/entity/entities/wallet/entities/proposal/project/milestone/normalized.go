package milestone

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project/feature"
)

type normalizedMilestone struct {
	ID          string             `json:"id"`
	Project     project.Normalized `json:"project"`
	Feature     feature.Normalized `json:"feature"`
	Wallet      wallet.Normalized  `json:"wallet"`
	Shares      int                `json:"shares"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Details     string             `json:"details"`
}

func createNormalizedMilestone(ins Milestone) (*normalizedMilestone, error) {
	proj, projErr := project.SDKFunc.CreateMetaData().Normalize()(ins.Project())
	if projErr != nil {
		return nil, projErr
	}

	wal, walErr := wallet.SDKFunc.CreateMetaData().Normalize()(ins.Wallet())
	if walErr != nil {
		return nil, walErr
	}

	var feat feature.Normalized
	if ins.HasFeature() {
		featNorm, featNormErr := feature.SDKFunc.CreateMetaData().Normalize()(ins.Feature())
		if featNormErr != nil {
			return nil, featNormErr
		}

		feat = featNorm
	}

	out := normalizedMilestone{
		ID:          ins.ID().String(),
		Project:     proj,
		Feature:     feat,
		Wallet:      wal,
		Shares:      ins.Shares(),
		Title:       ins.Title(),
		Description: ins.Description(),
		Details:     ins.Details(),
	}

	return &out, nil
}
