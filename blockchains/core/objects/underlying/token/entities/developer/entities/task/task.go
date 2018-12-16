package task

import (
	"errors"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/developer"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/developer/entities/milestone"
)

type task struct {
	UUID   *uuid.UUID          `json:"id"`
	Mstone milestone.Milestone `json:"milestone"`
	Crea   developer.Developer `json:"creator"`
	Titl   string              `json:"title"`
	Desc   string              `json:"description"`
	CrOn   time.Time           `json:"created_on"`
	DOn    time.Time           `json:"due_on"`
}

func createTask(id *uuid.UUID, mstone milestone.Milestone, createdBy developer.Developer, title string, description string, createdOn time.Time, dueOn time.Time) Task {
	out := task{
		UUID:   id,
		Mstone: mstone,
		Crea:   createdBy,
		Titl:   title,
		Desc:   description,
		CrOn:   createdOn,
		DOn:    dueOn,
	}

	return &out
}

func createTaskFromNormalized(normalized *normalizedTask) (Task, error) {
	id, idErr := uuid.FromString(normalized.ID)
	if idErr != nil {
		return nil, idErr
	}

	mstoneIns, mstoneInsErr := milestone.SDKFunc.CreateMetaData().Denormalize()(normalized.Milestone)
	if mstoneInsErr != nil {
		return nil, mstoneInsErr
	}

	creatorIns, creatorInsErr := developer.SDKFunc.CreateMetaData().Denormalize()(normalized.Creator)
	if creatorInsErr != nil {
		return nil, creatorInsErr
	}

	if mstone, ok := mstoneIns.(milestone.Milestone); ok {
		if creator, ok := creatorIns.(developer.Developer); ok {
			out := createTask(&id, mstone, creator, normalized.Title, normalized.Description, normalized.CreatedOn, normalized.DueOn)
			return out, nil
		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid Developer instance", creatorIns.ID().String())
		return nil, errors.New(str)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Milestone instance", mstoneIns.ID().String())
	return nil, errors.New(str)

}

// ID returns the ID
func (obj *task) ID() *uuid.UUID {
	return obj.UUID
}

// Milestone returns the milestone
func (obj *task) Milestone() milestone.Milestone {
	return obj.Mstone
}

// Creator returns the creator
func (obj *task) Creator() developer.Developer {
	return obj.Crea
}

// Title returns the title
func (obj *task) Title() string {
	return obj.Titl
}

// Description returns the description
func (obj *task) Description() string {
	return obj.Desc
}

// CreatedOn returns the creation time
func (obj *task) CreatedOn() time.Time {
	return obj.CrOn
}

// DueOn returns the due-on time
func (obj *task) DueOn() time.Time {
	return obj.DOn
}
