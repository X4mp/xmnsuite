package xmn

import (
	"math/rand"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/tests"
)

func createUserForTests() User {
	pubKey := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{}).PublicKey()
	shares := rand.Int()
	out := createUser(pubKey, shares)
	return out
}

func createWalletForTests() Wallet {
	id := uuid.NewV4()
	concensusNeeded := rand.Float64()
	usrs := []User{
		createUserForTests(),
		createUserForTests(),
	}

	out := createWallet(&id, concensusNeeded, usrs)
	return out
}

func TestUser_Success(t *testing.T) {
	usr := createUserForTests()

	empty := new(user)
	tests.ConvertToBinary(t, usr, empty, cdc)

	anotherEmpty := new(user)
	tests.ConvertToJSON(t, usr, anotherEmpty, cdc)
}

func TestWallet_Success(t *testing.T) {
	wal := createWalletForTests()

	empty := new(wallet)
	tests.ConvertToBinary(t, wal, empty, cdc)

	anotherEmpty := new(wallet)
	tests.ConvertToJSON(t, wal, anotherEmpty, cdc)
}
