package project

import (
	"reflect"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal"
)

// CreateProjectWithProposalForTests creates a project with proposal for tests
func CreateProjectWithProposalForTests(prop proposal.Proposal) Project {
	id := uuid.NewV4()
	out, outErr := createProject(&id, prop)
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

	// compare proposal:
	proposal.CompareProposalsForTests(t, first.Proposal(), second.Proposal())
}
