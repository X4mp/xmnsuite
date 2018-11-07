package user

import (
	"math/rand"
	"reflect"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/wallet"
	"github.com/xmnservices/xmnsuite/crypto"
)

// CreateUserForTests creates a User for tests
func CreateUserForTests() User {
	id := uuid.NewV4()
	pubKey := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{}).PublicKey()
	shares := rand.Int()
	wal := wallet.CreateWalletForTests()
	out := createUser(&id, pubKey, shares, wal)
	return out
}

// CreateUserWithSharesAndConcensusNeededForTests creates a User with shares and concensusNeeded for tests
func CreateUserWithSharesAndConcensusNeededForTests(concensusNeeded int, shares int) User {
	id := uuid.NewV4()
	pubKey := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{}).PublicKey()
	wal := wallet.CreateWalletWithConcensusNeededForTests(concensusNeeded)
	out := createUser(&id, pubKey, shares, wal)
	return out
}

// createUserWithSharesForTests creates a User with shares for tests
func createUserWithSharesForTests(shares int) User {
	id := uuid.NewV4()
	pubKey := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{}).PublicKey()
	wal := wallet.CreateWalletForTests()
	out := createUser(&id, pubKey, shares, wal)
	return out
}

// CreateUserWithWalletForTests creates a User with Wallet for tests
func CreateUserWithWalletForTests(wal wallet.Wallet) User {
	id := uuid.NewV4()
	pubKey := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{}).PublicKey()
	shares := rand.Int()
	out := createUser(&id, pubKey, shares, wal)
	return out
}

// CreateUserWithPublicKeyForTests creates a User with PubKey for tests
func CreateUserWithPublicKeyForTests(pubKey crypto.PublicKey) User {
	id := uuid.NewV4()
	shares := rand.Int()
	wal := wallet.CreateWalletForTests()
	out := createUser(&id, pubKey, shares, wal)
	return out
}

// CreateUserWithWalletAndPublicKeyForTests creates a User with Wallet and PubKey for tests
func CreateUserWithWalletAndPublicKeyForTests(wal wallet.Wallet, pubKey crypto.PublicKey) User {
	id := uuid.NewV4()
	shares := rand.Int()
	out := createUser(&id, pubKey, shares, wal)
	return out
}

// CompareUserForTests compare User instances for tests
func CompareUserForTests(t *testing.T, firstUser User, secondUser User) {
	if !reflect.DeepEqual(firstUser.ID(), secondUser.ID()) {
		t.Errorf("the first ID (%s) does not match the second ID (%s)", firstUser.ID().String(), secondUser.ID().String())
		return
	}

	if !firstUser.PubKey().Equals(secondUser.PubKey()) {
		t.Errorf("the public keys do not match")
		return
	}

	if firstUser.Shares() != secondUser.Shares() {
		t.Errorf("the first shares (%d) do not match the second shares (%d)", firstUser.Shares(), secondUser.Shares())
		return
	}

	wallet.CompareWalletsForTests(t, firstUser.Wallet(), secondUser.Wallet())
}