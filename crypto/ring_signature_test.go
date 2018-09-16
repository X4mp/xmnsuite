package crypto

import (
	"testing"

	"github.com/dedis/kyber"
)

func TestRingSignature_Success(t *testing.T) {
	// variables:
	msg := "this is a message to sign"
	pk := createPrivateKey()
	secondPK := createPrivateKey()
	ringPubKeys := []kyber.Point{
		pk.PublicKey(),
		secondPK.PublicKey(),
	}

	firstRing := pk.RingSign(msg, ringPubKeys)
	secondRing := secondPK.RingSign(msg, ringPubKeys)

	if !firstRing.Verify(msg) {
		t.Errorf("the first ring was expected to be verified")
		return
	}

	if !secondRing.Verify(msg) {
		t.Errorf("the second ring was expected to be verified")
		return
	}
}

func TestRingSignature_PubKeyIsNotInTheRing_panic_Success(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			return
		}

		t.Errorf("the func was expected to panic")
	}()

	// variables:
	msg := "this is a message to sign"
	pk := createPrivateKey()
	secondPK := createPrivateKey()
	invalidPK := createPrivateKey()
	ringPubKeys := []kyber.Point{
		pk.PublicKey(),
		secondPK.PublicKey(),
	}

	invalidPK.RingSign(msg, ringPubKeys)
}
