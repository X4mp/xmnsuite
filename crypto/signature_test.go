package crypto

import (
	"testing"
)

func TestSignature_Success(t *testing.T) {
	// variables:
	msg := "this is a message to sign"
	pk := createPrivateKey()
	invalidPK := createPrivateKey()

	// create the signature:
	sig := pk.Sign(msg)

	// derive a PublicKey from the signature:
	derivedPubKey := sig.PublicKey(msg)
	invalidDerivedPubKey := sig.PublicKey("invalid msg")

	// make sure the original PublicKey and the derived PublicKey are the same:
	if !pk.PublicKey().Equal(derivedPubKey) {
		t.Errorf("the original PublicKey was expected to be the same as the derived PublicKey")
		return
	}

	// verify the signature:
	if !sig.Verify(msg, derivedPubKey) {
		t.Errorf("the signature was expected to be verified using this message and PublicKey")
		return
	}

	// verify the signature with an invalid derived PublicKey:
	if sig.Verify(msg, invalidDerivedPubKey) {
		t.Errorf("the signature was expected to be verified using this message and invalid derived PublicKey")
		return
	}

	// verify the signature with an invalid pubKey:
	if sig.Verify(msg, invalidPK.PublicKey()) {
		t.Errorf("the signature was NOT expected to be verified with an invalid PublicKey")
		return
	}

}
