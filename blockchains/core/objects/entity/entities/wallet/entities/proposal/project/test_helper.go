package project

import (
	"math/rand"
	"reflect"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	community_project "github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/project"
)

// CreateProjectWithCommunityProjectAndWallets creates a project with a community project for tests
func CreateProjectWithCommunityProjectAndWallets(commProject community_project.Project, owner wallet.Wallet, manager wallet.Wallet, linker wallet.Wallet) Project {
	id := uuid.NewV4()
	mgrShares := rand.Int()%20 + 1
	lnkShares := rand.Int()%20 + 1
	wrkShares := rand.Int()%20 + 1
	out, outErr := createProject(&id, commProject, owner, manager, mgrShares, linker, lnkShares, wrkShares)
	if outErr != nil {
		panic(outErr)
	}

	return out
}

// CompareProjectsForTests compare Project instances for tests
func CompareProjectsForTests(t *testing.T, first Project, second Project) {
	if !reflect.DeepEqual(first.ID(), second.ID()) {
		t.Errorf("the IDs are different.  Expected: %s, Returned: %s", first.ID().String(), second.ID().String())
		return
	}

	if first.ManagerShares() != second.ManagerShares() {
		t.Errorf("the managerShares is different.  Expected: %d, Returned: %d", first.ManagerShares(), second.ManagerShares())
		return
	}

	if first.LinkerShares() != second.LinkerShares() {
		t.Errorf("the linkerShares is different.  Expected: %d, Returned: %d", first.LinkerShares(), second.LinkerShares())
		return
	}

	if first.WorkerShares() != second.WorkerShares() {
		t.Errorf("the workerShares is different.  Expected: %d, Returned: %d", first.WorkerShares(), second.WorkerShares())
		return
	}

	// compare instances:
	community_project.CompareProjectsForTests(t, first.Project(), second.Project())
	wallet.CompareWalletsForTests(t, first.Owner(), second.Owner())
	wallet.CompareWalletsForTests(t, first.Manager(), second.Manager())
	wallet.CompareWalletsForTests(t, first.Linker(), second.Linker())
}
