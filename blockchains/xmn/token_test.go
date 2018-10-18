package xmn

import (
	"testing"

	"github.com/xmnservices/xmnsuite/datastore"
	"github.com/xmnservices/xmnsuite/tests"
)

func createTokenForTests() Token {
	symbol := "XMN"
	name := "XMN Token"
	desc := "This is the XMN token"
	out := createToken(symbol, name, desc)
	return out
}

func compareTokensForTests(t *testing.T, first Token, second Token) {
	if first.Symbol() != second.Symbol() {
		t.Errorf("the symbol is different.  Expected: %s, Returned: %s", first.Symbol(), second.Symbol())
		return
	}

	if first.Name() != second.Name() {
		t.Errorf("the name is different.  Expected: %s, Returned: %s", first.Name(), second.Name())
		return
	}

	if first.Description() != second.Description() {
		t.Errorf("the description is different.  Expected: %s, Returned: %s", first.Description(), second.Description())
		return
	}

}

func TestToken_Success(t *testing.T) {
	tok := createTokenForTests()

	// create the service:
	store := datastore.SDKFunc.Create()
	serv := createTokenService(store)

	// save:
	saveErr := serv.Save(tok)
	if saveErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveErr.Error())
		return
	}

	// save again, returns error:
	saveAgainErr := serv.Save(createTokenForTests())
	if saveAgainErr == nil {
		t.Errorf("the returned error was expected to be an error, nil returned")
		return
	}

	// retrieve:
	retTok, retTokErr := serv.Retrieve()
	if retTokErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retTokErr.Error())
		return
	}

	// compare:
	compareTokensForTests(t, tok, retTok)

	// convert back and forth:
	empty := new(token)
	tests.ConvertToBinary(t, tok, empty, cdc)

	anotherEmpty := new(token)
	tests.ConvertToJSON(t, tok, anotherEmpty, cdc)
}
