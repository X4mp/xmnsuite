package task

import (
	"errors"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project/milestone"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
)

type task struct {
	UUID        *uuid.UUID          `json:"id"`
	Mils        milestone.Milestone `json:"milestone"`
	CrBy        user.User           `json:"created_by"`
	Titl        string              `json:"title"`
	Det         string              `json:"details"`
	DeadLn      time.Time           `json:"deadline"`
	Rewrd       int                 `json:"reward"`
	PldgeNeeded int                 `json:"pledge_needed"`
}

func createTask(
	id *uuid.UUID,
	mils milestone.Milestone,
	crBy user.User,
	title string,
	details string,
	deadline time.Time,
	reward int,
	pledgeNeeded int,
) (Task, error) {
	out := task{
		UUID:        id,
		Mils:        mils,
		CrBy:        crBy,
		Titl:        title,
		Det:         details,
		DeadLn:      deadline,
		Rewrd:       reward,
		PldgeNeeded: pledgeNeeded,
	}

	return &out, nil
}

func createTaskFromNormalized(normalized *normalizedTask) (Task, error) {
	id, idErr := uuid.FromString(normalized.ID)
	if idErr != nil {
		return nil, idErr
	}

	milsIns, milsInsErr := milestone.SDKFunc.CreateMetaData().Denormalize()(normalized.Milestone)
	if milsInsErr != nil {
		return nil, milsInsErr
	}

	crByIns, crByInsErr := user.SDKFunc.CreateMetaData().Denormalize()(normalized.CreatedBy)
	if crByInsErr != nil {
		return nil, crByInsErr
	}

	if mils, ok := milsIns.(milestone.Milestone); ok {
		if crBy, ok := crByIns.(user.User); ok {
			return createTask(&id, mils, crBy, normalized.Title, normalized.Details, normalized.Deadline, normalized.Reward, normalized.PledgeNeeded)
		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid User instance", crByIns.ID().String())
		return nil, errors.New(str)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Milestone instance", milsIns.ID().String())
	return nil, errors.New(str)
}

func createTaskFromStorable(storable *storableTask, rep entity.Repository) (Task, error) {
	id, idErr := uuid.FromString(storable.ID)
	if idErr != nil {
		return nil, idErr
	}

	milsID, milsIDErr := uuid.FromString(storable.MilestoneID)
	if milsIDErr != nil {
		return nil, milsIDErr
	}

	crByID, crByIDErr := uuid.FromString(storable.CreatedByUserID)
	if crByIDErr != nil {
		return nil, crByIDErr
	}

	milsIns, milsInsErr := rep.RetrieveByID(milestone.SDKFunc.CreateMetaData(), &milsID)
	if milsInsErr != nil {
		return nil, milsInsErr
	}

	crByIns, crByInsErr := rep.RetrieveByID(user.SDKFunc.CreateMetaData(), &crByID)
	if crByInsErr != nil {
		return nil, crByInsErr
	}

	if mils, ok := milsIns.(milestone.Milestone); ok {
		if crBy, ok := crByIns.(user.User); ok {
			return createTask(&id, mils, crBy, storable.Title, storable.Details, storable.Deadline, storable.Reward, storable.PledgeNeeded)
		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid User instance", crByIns.ID().String())
		return nil, errors.New(str)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Milestone instance", milsIns.ID().String())
	return nil, errors.New(str)
}

// ID returns the ID
func (obj *task) ID() *uuid.UUID {
	return obj.UUID
}

// Milestone returns the milestone
func (obj *task) Milestone() milestone.Milestone {
	return obj.Mils
}

// CreatedBy returns the createdBy user
func (obj *task) CreatedBy() user.User {
	return obj.CrBy
}

// Title returns the title
func (obj *task) Title() string {
	return obj.Titl
}

// Details returns the details
func (obj *task) Details() string {
	return obj.Det
}

// Deadline returns the deadline
func (obj *task) Deadline() time.Time {
	return obj.DeadLn
}

// Reward returns the reward
func (obj *task) Reward() int {
	return obj.Rewrd
}

// PledgeNeeded returns the pledge needed
func (obj *task) PledgeNeeded() int {
	return obj.PldgeNeeded
}
