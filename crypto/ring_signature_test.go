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
	invalidPK := createPrivateKey()
	ringPubKeys := []kyber.Point{
		pk.PublicKey(),
		secondPK.PublicKey(),
	}

	firstRing := pk.RingSign(msg, ringPubKeys, 0)
	secondRing := secondPK.RingSign(msg, ringPubKeys, 1)
	invalidRing := invalidPK.RingSign(msg, ringPubKeys, 0)

	if !firstRing.Verify(msg) {
		t.Errorf("the first ring was expected to be verified")
		return
	}

	if !secondRing.Verify(msg) {
		t.Errorf("the second ring was expected to be verified")
		return
	}

	if invalidRing.Verify(msg) {
		t.Errorf("the invalid ring was NOT expected to be verified")
		return
	}
}
