package crypto

import (
	"testing"
)

func TestPrivateKey_Success(t *testing.T) {
	// variables:
	pk := createPrivateKey()
	pkAsString := pk.String()
	unmarshalPK, unmarshalPKErr := createPrivateKeyFromString(pkAsString)
	if unmarshalPKErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", unmarshalPKErr.Error())
		return
	}

	if pk.String() != unmarshalPK.String() {
		t.Errorf("the PrivateKeys were expected to be the same.  Expected: %s, Returned: %s", pk.String(), unmarshalPK.String())
		return
	}
}
