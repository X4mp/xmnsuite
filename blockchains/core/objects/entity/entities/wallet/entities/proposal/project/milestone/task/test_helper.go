package task

import (
	"reflect"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project/milestone"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
)

// CreateTaskWithMilestoneAndUser creates a task with milestone and user for tests
func CreateTaskWithMilestoneAndUser(mils milestone.Milestone, createdBy user.User) Task {
	id := uuid.NewV4()
	out, outErr := createTask(&id, mils, createdBy, "This is the title", "This is the details", time.Now(), 22, 2)
	if outErr != nil {
		panic(outErr)
	}

	return out
}

// CompareTasksForTests compare tasks instances for tests
func CompareTasksForTests(t *testing.T, first Task, second Task) {
	if !reflect.DeepEqual(first.ID(), second.ID()) {
		t.Errorf("the IDs are different.  Expected: %s, Returned: %s", first.ID().String(), second.ID().String())
		return
	}

	if first.Title() != second.Title() {
		t.Errorf("the title is different.  Expected: %s, Returned: %s", first.Title(), second.Title())
		return
	}

	if first.Details() != second.Details() {
		t.Errorf("the details is different.  Expected: %s, Returned: %s", first.Details(), second.Details())
		return
	}

	if first.Deadline() != second.Deadline() {
		t.Errorf("the deadline is different.  Expected: %s, Returned: %s", first.Deadline(), second.Deadline())
		return
	}

	if first.Reward() != second.Reward() {
		t.Errorf("the reward is different.  Expected: %d, Returned: %d", first.Reward(), second.Reward())
		return
	}

	if first.PledgeNeeded() != second.PledgeNeeded() {
		t.Errorf("the pledgeNeeded is different.  Expected: %d, Returned: %d", first.PledgeNeeded(), second.PledgeNeeded())
		return
	}

	// compare instances:
	milestone.CompareMilestonesForTests(t, first.Milestone(), second.Milestone())
	user.CompareUserForTests(t, first.CreatedBy(), second.CreatedBy())
}
