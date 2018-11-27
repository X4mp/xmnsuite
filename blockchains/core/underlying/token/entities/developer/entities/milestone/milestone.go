package milestone

import (
	"errors"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/entities/developer/entities/project"
)

type milestone struct {
	UUID *uuid.UUID      `json:"id"`
	Proj project.Project `json:"project"`
	Titl string          `json:"title"`
	Desc string          `json:"description"`
	CrOn time.Time       `json:"created_on"`
	DuOn time.Time       `json:"due_on"`
}

func createMilestone(id *uuid.UUID, project project.Project, title string, description string, createdOn time.Time, dueOn time.Time) Milestone {
	out := milestone{
		UUID: id,
		Proj: project,
		Titl: title,
		Desc: description,
		CrOn: createdOn,
		DuOn: dueOn,
	}

	return &out
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

	if proj, ok := projIns.(project.Project); ok {
		out := createMilestone(&id, proj, normalized.Title, normalized.Description, normalized.CreatedOn, normalized.DueOn)
		return out, nil
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

// Title returns the title
func (obj *milestone) Title() string {
	return obj.Titl
}

// Description returns the description
func (obj *milestone) Description() string {
	return obj.Desc
}

// CreatedOn returns the creation time
func (obj *milestone) CreatedOn() time.Time {
	return obj.CrOn
}

// DueOn returns the due time
func (obj *milestone) DueOn() time.Time {
	return obj.DuOn
}
