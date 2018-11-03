package token

import (
	"reflect"
	"testing"

	uuid "github.com/satori/go.uuid"
)

// CreateTokenForTests creates a Token instance for tests
func CreateTokenForTests() Token {
	id := uuid.NewV4()
	symbol := "XMN"
	name := "XMN Token"
	desc := "This is the XMN token"
	out := createToken(&id, symbol, name, desc)
	return out
}

// CompareTokensForTests compare Token instances for tests
func CompareTokensForTests(t *testing.T, first Token, second Token) {
	if !reflect.DeepEqual(first.ID(), second.ID()) {
		t.Errorf("the IDs are different.  Expected: %s, Returned: %s", first.ID().String(), second.ID().String())
		return
	}

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
