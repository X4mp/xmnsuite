package crypto

import (
	"encoding/hex"
	"testing"
)

func TestSDK_createPK_Success(t *testing.T) {
	// variables:
	pk := SDKFunc.GenPK()
	pkAsString := pk.String()
	unmarshalPK := SDKFunc.CreatePK(CreatePKParams{
		PKAsString: pkAsString,
	})

	if pk.String() != unmarshalPK.String() {
		t.Errorf("the PrivateKeys were expected to be the same.  Expected: %s, Returned: %s", pk.String(), unmarshalPK.String())
		return
	}
}

func TestSDK_createPK_withoutPK_Success(t *testing.T) {
	// variables:
	pk := SDKFunc.CreatePK(CreatePKParams{})
	pkAsString := pk.String()
	unmarshalPK := SDKFunc.CreatePK(CreatePKParams{
		PKAsString: pkAsString,
	})

	if pk.String() != unmarshalPK.String() {
		t.Errorf("the PrivateKeys were expected to be the same.  Expected: %s, Returned: %s", pk.String(), unmarshalPK.String())
		return
	}
}

func TestSDK_createPK_withInvalidHexPK_panic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			return
		}

		t.Errorf("the func was expected to panic")
	}()

	// variables:
	SDKFunc.CreatePK(CreatePKParams{
		PKAsString: "this is not a PK",
	})
}

func TestSDK_createPK_withValidHex_invalidPK_panic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			return
		}

		t.Errorf("the func was expected to panic")
	}()

	// variables:
	SDKFunc.CreatePK(CreatePKParams{
		PKAsString: hex.EncodeToString([]byte("this is not a PK")),
	})
}

func TestSDK_createPubKey_Success(t *testing.T) {
	//variables:
	p := curve.Point().Base()
	pubKey := createPublicKey(p)
	pubKeyAsString := pubKey.String()

	// execute:
	samePubKey := SDKFunc.CreatePubKey(CreatePubKeyParams{
		PubKeyAsString: pubKeyAsString,
	})

	if !pubKey.Equals(samePubKey) {
		t.Errorf("the public keys should be equal")
		return
	}
}

func TestSDK_createPubKey_withValidHex_invalidPubKey_panic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			return
		}

		t.Errorf("the func was expected to panic")
	}()

	// variables:
	SDKFunc.CreatePubKey(CreatePubKeyParams{
		PubKeyAsString: hex.EncodeToString([]byte("this is an invalid pubkey")),
	})
}
