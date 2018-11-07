package wallet

import (
	"math/rand"
	"reflect"
	"testing"

	uuid "github.com/satori/go.uuid"
)

// CreateWalletForTests creates a Wallet instance, for tests
func CreateWalletForTests() Wallet {
	concensusNeeded := rand.Int()
	return CreateWalletWithConcensusNeededForTests(concensusNeeded)
}

// CreateWalletWithConcensusNeededForTests creates a Wallet instance with a required concensus, for tests
func CreateWalletWithConcensusNeededForTests(concensusNeeded int) Wallet {
	id := uuid.NewV4()
	out := createWallet(&id, concensusNeeded)
	return out
}

// CompareWalletsForTests compares 2 Wallet instances for tests
func CompareWalletsForTests(t *testing.T, first Wallet, second Wallet) {
	if !reflect.DeepEqual(first.ID(), second.ID()) {
		t.Errorf("the ID is invalid")
		return
	}

	if !reflect.DeepEqual(first.ConcensusNeeded(), second.ConcensusNeeded()) {
		t.Errorf("the concensusNeeded is invalid.  Expected: %d, Returned: %d", first.ConcensusNeeded(), second.ConcensusNeeded())
		return
	}
}
