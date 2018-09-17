package crypto

import (
	"github.com/dedis/kyber"
)

type pubKey struct {
	p kyber.Point
}

type jsonPubKey struct {
	PubKeyAsString string `json:"pubkey"`
}

func createPublicKey(p kyber.Point) PublicKey {
	out := pubKey{
		p: p,
	}

	return &out
}

func createPublicKeyFromString(str string) (PublicKey, error) {
	p, pErr := fromStringToPoint(str)
	if pErr != nil {
		return nil, pErr
	}

	return createPublicKey(p), nil
}

// Point returns the point
func (obj *pubKey) Point() kyber.Point {
	return obj.p
}

// Equals returns true if the given PublicKey equals the current one
func (obj *pubKey) Equals(pubKey PublicKey) bool {
	return obj.p.Equal(pubKey.Point())
}

// String returns the string representation of the public key
func (obj *pubKey) String() string {
	return obj.p.String()
}

// MarshalJSON returns a JSON representation of the object
func (obj *pubKey) MarshalJSON() ([]byte, error) {
	return cdc.MarshalJSON(jsonPubKey{
		PubKeyAsString: obj.String(),
	})
}

// UnmarshalJSON returns an object based on the JSON data
func (obj *pubKey) UnmarshalJSON(data []byte) error {
	ptr := new(jsonPubKey)
	err := cdc.UnmarshalJSON(data, ptr)
	if err != nil {
		return err
	}

	p, pErr := fromStringToPoint(ptr.PubKeyAsString)
	if pErr != nil {
		return pErr
	}

	obj.p = p
	return nil
}
