package milestone

import (
	"math/rand"
	"reflect"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project"
)

// CreateMilestoneWithProjectForTests creates a milestone with project for tests
func CreateMilestoneWithProjectForTests(proj project.Project) Milestone {
	id := uuid.NewV4()
	share := rand.Int()%20 + 1
	wal := wallet.CreateWalletForTests()
	out, outErr := createMilestone(&id, proj, wal, share, "This is the milestone title", "This is the milestone description", "this is the milestone details")
	if outErr != nil {
		panic(outErr)
	}

	return out
}

// CompareMilestonesForTests compare Milestone instances for tests
func CompareMilestonesForTests(t *testing.T, first Milestone, second Milestone) {
	if !reflect.DeepEqual(first.ID(), second.ID()) {
		t.Errorf("the IDs are different.  Expected: %s, Returned: %s", first.ID().String(), second.ID().String())
		return
	}

	if first.Title() != second.Title() {
		t.Errorf("the title is different.  Expected: %s, Returned: %s", first.Title(), second.Title())
		return
	}

	if first.Description() != second.Description() {
		t.Errorf("the description is different.  Expected: %s, Returned: %s", first.Description(), second.Description())
		return
	}

	if first.Details() != second.Details() {
		t.Errorf("the details is different.  Expected: %s, Returned: %s", first.Details(), second.Details())
		return
	}

	// compare project and wallet:
	project.CompareProjectsForTests(t, first.Project(), second.Project())
	wallet.CompareWalletsForTests(t, first.Wallet(), second.Wallet())
}
