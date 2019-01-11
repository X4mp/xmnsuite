package task

import (
	"time"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project/milestone"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
)

type normalizedTask struct {
	ID           string               `json:"id"`
	Milestone    milestone.Normalized `json:"milestone"`
	CreatedBy    user.Normalized      `json:"created_by"`
	Title        string               `json:"title"`
	Details      string               `json:"details"`
	Deadline     time.Time            `json:"deadline"`
	Reward       int                  `json:"reward"`
	PledgeNeeded int                  `json:"pledge_needed"`
}

func createNormalizedTask(ins Task) (*normalizedTask, error) {
	mils, milsErr := milestone.SDKFunc.CreateMetaData().Normalize()(ins.Milestone())
	if milsErr != nil {
		return nil, milsErr
	}

	crBy, crByErr := user.SDKFunc.CreateMetaData().Normalize()(ins.CreatedBy())
	if crByErr != nil {
		return nil, crByErr
	}

	out := normalizedTask{
		ID:           ins.ID().String(),
		Milestone:    mils,
		CreatedBy:    crBy,
		Title:        ins.Title(),
		Details:      ins.Details(),
		Deadline:     ins.Deadline(),
		Reward:       ins.Reward(),
		PledgeNeeded: ins.PledgeNeeded(),
	}

	return &out, nil
}
