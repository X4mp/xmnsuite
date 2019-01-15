package completed

import (
	"reflect"
	"testing"

	uuid "github.com/satori/go.uuid"
	mils_task "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project/milestone/task"
)

// CreateTaskWithMilestoneTask creates a task with milestone and user for tests
func CreateTaskWithMilestoneTask(tsk mils_task.Task) Task {
	id := uuid.NewV4()
	out, outErr := createTask(&id, tsk, "this is some details")
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

	if first.Details() != second.Details() {
		t.Errorf("the details are different.  Expected: %s, Returned: %s", first.Details(), second.Details())
		return
	}

	// compare instances:
	mils_task.CompareTasksForTests(t, first.Task(), second.Task())
}
