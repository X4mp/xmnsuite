package crypto

import (
	"testing"

	"github.com/xmnservices/xmnsuite/tests"
)

func TestRingSignature_Success(t *testing.T) {
	// variables:
	msg := "this is a message to sign"
	pk := createPrivateKey()
	secondPK := createPrivateKey()
	ringPubKeys := []PublicKey{
		pk.PublicKey(),
		secondPK.PublicKey(),
	}

	firstRing, firstRingErr := pk.RingSign(msg, ringPubKeys)
	if firstRingErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned; %s", firstRingErr.Error())
		return
	}

	secondRing, secondRingErr := secondPK.RingSign(msg, ringPubKeys)
	if secondRingErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned; %s", secondRingErr.Error())
		return
	}

	if !firstRing.Verify(msg) {
		t.Errorf("the first ring was expected to be verified")
		return
	}

	if !secondRing.Verify(msg) {
		t.Errorf("the second ring was expected to be verified")
		return
	}

	// encode to striong, back and forth:
	firstRingAsString := firstRing.String()
	newRing, newRingErr := createRingSignatureFromString(firstRingAsString)
	if newRingErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", newRingErr.Error())
		return
	}

	if firstRingAsString != newRing.String() {
		t.Errorf("the rings were expected to be the same.  Expected: %s, Actual: %s", firstRingAsString, newRing.String())
		return
	}

	// convert to json back and forth:
	empty := new(ringSignature)
	tests.ConvertToJSON(t, firstRing, empty, cdc)
}

func TestRingSignature_PubKeyIsNotInTheRing_returnsError(t *testing.T) {
	// variables:
	msg := "this is a message to sign"
	pk := createPrivateKey()
	secondPK := createPrivateKey()
	invalidPK := createPrivateKey()
	ringPubKeys := []PublicKey{
		pk.PublicKey(),
		secondPK.PublicKey(),
	}

	_, err := invalidPK.RingSign(msg, ringPubKeys)
	if err == nil {
		t.Errorf("the returned error was expected to be valid, nil returned")
		return
	}
}
