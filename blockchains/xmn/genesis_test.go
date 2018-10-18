package xmn

import (
	"math/rand"
	"testing"

	"github.com/xmnservices/xmnsuite/datastore"
	"github.com/xmnservices/xmnsuite/tests"
)

func createGenesisForTests() Genesis {
	gazPricePerKb := rand.Int()
	maxAmountOfValidators := rand.Intn(200)
	dep := createInitialDepositForTests()
	tok := createTokenForTests()
	out := createGenesis(gazPricePerKb, maxAmountOfValidators, dep, tok)
	return out
}

func compareGenesisForTests(t *testing.T, first Genesis, second Genesis) {
	if first.GazPricePerKb() != second.GazPricePerKb() {
		t.Errorf("the returned gaz price is invalid.  Expected: %d, Returned: %d", first.GazPricePerKb(), second.GazPricePerKb())
		return
	}

	if first.MaxAmountOfValidators() != second.MaxAmountOfValidators() {
		t.Errorf("the returned maximum amount of validatoirs is invalid.  Expected: %d, Returned: %d", first.MaxAmountOfValidators(), second.MaxAmountOfValidators())
		return
	}

	compareInitialDepositForTests(t, first.Deposit(), second.Deposit())
	compareTokensForTests(t, first.Token(), second.Token())
}

func TestGenesis_Success(t *testing.T) {
	gen := createGenesisForTests()

	// create services:
	store := datastore.SDKFunc.Create()
	walService := createWalletService(store)
	tokenService := createTokenService(store)
	initialDepService := createInitialDepositService(store, walService)
	genesisService := createGenesisService(store, walService, tokenService, initialDepService)

	// save the genesis:
	saveErr := genesisService.Save(gen)
	if saveErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveErr.Error())
		return
	}

	// sagain again, returns error:
	saveAgainErr := genesisService.Save(gen)
	if saveAgainErr == nil {
		t.Errorf("the returned was expected to be valid, nil returned")
		return
	}

	// retrieve genesis:
	retGen, retGenErr := genesisService.Retrieve()
	if retGenErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retGenErr.Error())
		return
	}

	// compare:
	compareGenesisForTests(t, gen, retGen)

	empty := new(genesis)
	tests.ConvertToBinary(t, gen, empty, cdc)

	anotherEmpty := new(genesis)
	tests.ConvertToJSON(t, gen, anotherEmpty, cdc)
}
