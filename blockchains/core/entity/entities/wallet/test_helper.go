package wallet

import (
	"math/rand"
	"reflect"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/crypto"
)

// CreateWalletWithPublicKeyForTests creates a Wallet instance, for tests
func CreateWalletWithPublicKeyForTests(pubKey crypto.PublicKey) Wallet {
	concensusNeeded := rand.Int() % 100
	return CreateWalletWithPubKeyAndConcensusNeededForTests(pubKey, concensusNeeded)
}

// CreateWalletWithPubKeyAndConcensusNeededForTests creates a Wallet instance with a required concensus, for tests
func CreateWalletWithPubKeyAndConcensusNeededForTests(pubKey crypto.PublicKey, concensusNeeded int) Wallet {
	id := uuid.NewV4()
	out := createWallet(&id, pubKey, concensusNeeded)
	return out
}

// CompareWalletsForTests compares 2 Wallet instances for tests
func CompareWalletsForTests(t *testing.T, first Wallet, second Wallet) {
	if !reflect.DeepEqual(first.ID(), second.ID()) {
		t.Errorf("the ID is invalid, expected: %s, returned: %s", first.ID().String(), second.ID().String())
		return
	}

	if !reflect.DeepEqual(first.ConcensusNeeded(), second.ConcensusNeeded()) {
		t.Errorf("the concensusNeeded is invalid.  Expected: %d, Returned: %d", first.ConcensusNeeded(), second.ConcensusNeeded())
		return
	}
}
