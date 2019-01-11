package pledge

import (
	"reflect"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/pledge"
	mils_task "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project/milestone/task"
)

// CreateTaskWithMilestoneTaskAndPledge creates a task with milestone and user for tests
func CreateTaskWithMilestoneTaskAndPledge(tsk mils_task.Task, pldge pledge.Pledge) Task {
	id := uuid.NewV4()
	out, outErr := createTask(&id, tsk, pldge)
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

	// compare instances:
	mils_task.CompareTasksForTests(t, first.Task(), second.Task())
	pledge.ComparePledgesForTests(t, first.Pledge(), second.Pledge())
}
