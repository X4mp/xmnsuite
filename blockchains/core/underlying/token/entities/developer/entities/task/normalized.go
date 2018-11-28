package task

import (
	"time"

	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/entities/developer"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/entities/developer/entities/milestone"
)

type normalizedTask struct {
	ID          string               `json:"id"`
	Milestone   milestone.Normalized `json:"milestone"`
	Creator     developer.Normalized `json:"creator"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	CreatedOn   time.Time            `json:"created_on"`
	DueOn       time.Time            `json:"due_on"`
}

func createNormalizedTask(ins Task) (*normalizedTask, error) {
	mstone, mstoneErr := milestone.SDKFunc.CreateMetaData().Normalize()(ins.Milestone())
	if mstoneErr != nil {
		return nil, mstoneErr
	}

	creator, creatorErr := developer.SDKFunc.CreateMetaData().Normalize()(ins.Creator())
	if creatorErr != nil {
		return nil, creatorErr
	}

	out := normalizedTask{
		ID:          ins.ID().String(),
		Milestone:   mstone,
		Creator:     creator,
		Title:       ins.Title(),
		Description: ins.Description(),
		CreatedOn:   ins.CreatedOn(),
		DueOn:       ins.DueOn(),
	}

	return &out, nil
}
