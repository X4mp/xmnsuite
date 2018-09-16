package crypto

import (
	"testing"
)

func TestSignature_Success(t *testing.T) {
	// variables:
	msg := "this is a message to sign"
	pk := createPrivateKey()

	// create the signature:
	sig := pk.Sign(msg)

	// derive a PublicKey from the signature:
	derivedPubKey := sig.PublicKey(msg)

	// make sure the original PublicKey and the derived PublicKey are the same:
	if !pk.PublicKey().Equal(derivedPubKey) {
		t.Errorf("the original PublicKey was expected to be the same as the derived PublicKey")
		return
	}

	// verify the signature:
	if !sig.Verify(msg) {
		t.Errorf("the signature was expected to be verified using this message and PublicKey")
		return
	}

	// convert back and forth to string:
	sigAsString := sig.String()
	newSig, newSigErr := createSignatureFromString(sigAsString)
	if newSigErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", newSigErr.Error())
		return
	}

	if sigAsString != newSig.String() {
		t.Errorf("the signatures were expected to be the same.  Expected: %s, Actual: %s", sigAsString, newSig.String())
		return
	}

}
