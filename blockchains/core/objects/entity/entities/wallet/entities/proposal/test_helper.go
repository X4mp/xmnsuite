package proposal

import (
	"reflect"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/category"
)

// CreateProposalWithCategoryForTests creates a proposal with category for tests
func CreateProposalWithCategoryForTests(cat category.Category) Proposal {
	id := uuid.NewV4()
	out, outErr := createProposal(&id, "This is the proposal title", "This is the proposal description", "this is the proposal details", cat, 1, 2)
	if outErr != nil {
		panic(outErr)
	}

	return out
}

// CompareProposalsForTests compare Proposal instances for tests
func CompareProposalsForTests(t *testing.T, first Proposal, second Proposal) {
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

	// compare categories:
	category.CompareCategoriesForTests(t, first.Category(), second.Category())
}
