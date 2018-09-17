package crypto

import (
	"testing"

	tests "github.com/xmnservices/xmnsuite/tests"
)

func TestPublicKey_Success(t *testing.T) {
	//variables:
	p := curve.Point().Base()

	// execute:
	pKey := createPublicKey(p)
	pubKeyAsString := pKey.String()
	samePubKey, samePubKeyErr := createPublicKeyFromString(pubKeyAsString)

	if samePubKeyErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", samePubKeyErr.Error())
		return
	}

	if !pKey.Equals(samePubKey) {
		t.Errorf("the public keys should be equal")
		return
	}

	// convert to json back and forth:
	empty := new(pubKey)
	tests.ConvertToJSON(t, pKey, empty, cdc)
}
