package milestone

import (
	"time"

	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/entities/developer/entities/project"
)

type normalizedMilestone struct {
	ID          string             `json:"id"`
	Project     project.Normalized `json:"project"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	CreatedOn   time.Time          `json:"created_on"`
	DueOn       time.Time          `json:"due_on"`
}

func createNormalizedMilestone(ins Milestone) (*normalizedMilestone, error) {
	proj, projErr := project.SDKFunc.CreateMetaData().Normalize()(ins.Project())
	if projErr != nil {
		return nil, projErr
	}

	out := normalizedMilestone{
		ID:          ins.ID().String(),
		Project:     proj,
		Title:       ins.Title(),
		Description: ins.Description(),
		CreatedOn:   ins.CreatedOn(),
		DueOn:       ins.DueOn(),
	}

	return &out, nil
}
