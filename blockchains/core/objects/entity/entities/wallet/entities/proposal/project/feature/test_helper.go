package feature

import (
	"reflect"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
)

// CreateFeatureWithProjectAndCreatedByUser creates a feature with project for tests
func CreateFeatureWithProjectAndCreatedByUser(proj project.Project, createdBy user.User) Feature {
	id := uuid.NewV4()
	out, outErr := createFeature(&id, proj, "This is the title", "this is the details", createdBy)
	if outErr != nil {
		panic(outErr)
	}

	return out
}

// CompareFeaturesForTests compare features instances for tests
func CompareFeaturesForTests(t *testing.T, first Feature, second Feature) {
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

	// compare instances:
	project.CompareProjectsForTests(t, first.Project(), second.Project())
	user.CompareUserForTests(t, first.CreatedBy(), second.CreatedBy())
}
