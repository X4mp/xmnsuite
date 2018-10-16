package xmn

import (
	"math/rand"
	"testing"

	"github.com/xmnservices/xmnsuite/tests"
)

func createInitialDepositForTests() InitialDeposit {
	wallet := createWalletForTests()
	amount := rand.Int()
	out := createInitialDeposit(wallet, amount)
	return out
}

func createTokenForTests() Token {
	symbol := "XMN"
	name := "XMN Token"
	desc := "This is the XMN token"
	out := createToken(symbol, name, desc)
	return out
}

func createGenesisForTests() Genesis {
	gazPricePerKb := rand.Int()
	maxAmountOfValidators := rand.Intn(200)
	devs := createWalletForTests()
	dep := createInitialDepositForTests()
	tok := createTokenForTests()
	out := createGenesis(gazPricePerKb, maxAmountOfValidators, devs, dep, tok)
	return out
}

func TestInitialDeposit_Success(t *testing.T) {
	initialDep := createInitialDepositForTests()

	empty := new(initialDeposit)
	tests.ConvertToBinary(t, initialDep, empty, cdc)

	anotherEmpty := new(initialDeposit)
	tests.ConvertToJSON(t, initialDep, anotherEmpty, cdc)
}

func TestToken_Success(t *testing.T) {
	tok := createTokenForTests()

	empty := new(token)
	tests.ConvertToBinary(t, tok, empty, cdc)

	anotherEmpty := new(token)
	tests.ConvertToJSON(t, tok, anotherEmpty, cdc)
}

func TestGenesis_Success(t *testing.T) {
	gen := createGenesisForTests()

	empty := new(genesis)
	tests.ConvertToBinary(t, gen, empty, cdc)

	anotherEmpty := new(genesis)
	tests.ConvertToJSON(t, gen, anotherEmpty, cdc)
}
