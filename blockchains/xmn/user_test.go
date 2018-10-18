package xmn

import (
	"math/rand"
	"testing"

	"github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/tests"
)

func createUserForTests() User {
	pubKey := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{}).PublicKey()
	shares := rand.Int()
	wal := createWalletForTests()
	out := createUser(pubKey, shares, wal)
	return out
}

func TestUser_Success(t *testing.T) {
	usr := createUserForTests()

	empty := new(user)
	tests.ConvertToBinary(t, usr, empty, cdc)

	anotherEmpty := new(user)
	tests.ConvertToJSON(t, usr, anotherEmpty, cdc)
}
