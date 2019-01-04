package validator

import (
	"bytes"
	"encoding/hex"
	"math/rand"
	"net"
	"reflect"
	"testing"

	uuid "github.com/satori/go.uuid"
	ed25519 "github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/pledge"
)

// CreateValidatorForTests creates a Validator instance for tests
func CreateValidatorForTests() Validator {
	id := uuid.NewV4()
	pkey := ed25519.GenPrivKey().PubKey()
	pldge := pledge.CreatePledgeForTests()
	ip := net.ParseIP("127.0.0.1")
	port := rand.Int()%9000 + 1000
	out := createValidator(&id, ip, port, pkey, pldge)
	return out
}

// CompareValidatorsForTests compares 2 Validator instances for tests
func CompareValidatorsForTests(t *testing.T, first Validator, second Validator) {
	if !reflect.DeepEqual(first.ID(), second.ID()) {
		t.Errorf("the ID is invalid")
		return
	}

	if bytes.Compare(first.PubKey().Bytes(), second.PubKey().Bytes()) != 0 {
		t.Errorf("the pubKey instances are invalid.  Expected: %s, returned: %s", hex.EncodeToString(first.PubKey().Bytes()), hex.EncodeToString(second.PubKey().Bytes()))
		return
	}

	pledge.ComparePledgesForTests(t, first.Pledge(), second.Pledge())
}
