package tendermint

import (
	"reflect"
	"testing"
	"time"

	tests "github.com/xmnservices/xmnsuite/tests"
	uuid "github.com/satori/go.uuid"
	ed25519 "github.com/tendermint/tendermint/crypto/ed25519"
)

func TestValidator_Success(t *testing.T) {

	//variables:
	validatorName := "my-validator"
	validatorPower := 10

	//create the private key:
	privKey := ed25519.GenPrivKey()

	//create a validator:
	val := createValidator(validatorName, validatorPower, privKey.PubKey())

	//encode to json:
	valJS, valJSErr := cdc.MarshalJSON(val)
	if valJSErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", valJSErr.Error())
		return
	}

	//convert the json to validator:
	retVal := new(validator)
	jsErr := cdc.UnmarshalJSON(valJS, retVal)
	if jsErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", jsErr.Error())
		return
	}

	if !reflect.DeepEqual(val, retVal) {
		t.Errorf("the returned validator is invalid")
	}
}

func TestGenesis_Success(t *testing.T) {

	//variables:
	namespace := "xsuite"
	name := "users"
	id := uuid.NewV4()
	validatorName := "my-validator"
	validatorPower := 10
	createdOn := time.Now().UTC()

	//create the private key:
	privKey := ed25519.GenPrivKey()

	//create a validator:
	validator := createValidator(validatorName, validatorPower, privKey.PubKey())

	//create the path:
	path := createPath(namespace, name, &id)

	//create the genesis:
	gen, genErr := createGenesis([]byte(""), path, []Validator{
		validator,
	}, createdOn)

	if genErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", genErr.Error())
		return
	}

	//encode to json:
	genJS, genJSErr := cdc.MarshalJSON(gen)
	if genJSErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", genJSErr.Error())
		return
	}

	//convert the json to genesis:
	retGen := new(genesis)
	jsErr := cdc.UnmarshalJSON(genJS, retGen)
	if jsErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", jsErr.Error())
		return
	}

	if !reflect.DeepEqual(gen, retGen) {
		t.Errorf("the returned genesis is invalid")
	}
}

func TestGenerateBlockchain_Success(t *testing.T) {
	//variables:
	namespace := "xsuite"
	name := "users"
	id := uuid.NewV4()

	blkchain, blkchainErr := generateBlockchain(namespace, name, &id)
	if blkchainErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", blkchainErr.Error())
		return
	}

	//convert the blockchain, back and forth:
	empty := new(blockchain)
	tests.ConvertToJSON(t, blkchain, empty, cdc)
}
