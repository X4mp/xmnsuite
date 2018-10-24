package xmn

import (
	"math/rand"
	"reflect"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/tests"
)

func createUserForTests() User {
	id := uuid.NewV4()
	pubKey := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{}).PublicKey()
	shares := rand.Int()
	wal := createWalletForTests()
	out := createUser(&id, pubKey, shares, wal)
	return out
}

func createUserWithSharesAndConcensusNeededForTests(concensusNeeded int, shares int) User {
	id := uuid.NewV4()
	pubKey := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{}).PublicKey()
	wal := createWalletWithConcensusNeededForTests(concensusNeeded)
	out := createUser(&id, pubKey, shares, wal)
	return out
}

func createUserWithSharesForTests(shares int) User {
	id := uuid.NewV4()
	pubKey := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{}).PublicKey()
	wal := createWalletForTests()
	out := createUser(&id, pubKey, shares, wal)
	return out
}

func createUserWithWalletForTests(wal Wallet) User {
	id := uuid.NewV4()
	pubKey := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{}).PublicKey()
	shares := rand.Int()
	out := createUser(&id, pubKey, shares, wal)
	return out
}

func createUserWithPublicKeyForTests(pubKey crypto.PublicKey) User {
	id := uuid.NewV4()
	shares := rand.Int()
	wal := createWalletForTests()
	out := createUser(&id, pubKey, shares, wal)
	return out
}

func createUserWithWalletAndPublicKeyForTests(wal Wallet, pubKey crypto.PublicKey) User {
	id := uuid.NewV4()
	shares := rand.Int()
	out := createUser(&id, pubKey, shares, wal)
	return out
}

func compareUserForTests(t *testing.T, firstUser User, secondUser User) {
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

	compareWalletsForTests(t, firstUser.Wallet(), secondUser.Wallet())
}

func TestUser_Success(t *testing.T) {
	usr := createUserForTests()

	empty := new(user)
	tests.ConvertToBinary(t, usr, empty, cdc)

	anotherEmpty := new(user)
	tests.ConvertToJSON(t, usr, anotherEmpty, cdc)
}
