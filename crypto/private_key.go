package crypto

import (
	"errors"
	"fmt"

	"github.com/dedis/kyber"
)

type privateKey struct {
	x kyber.Scalar
}

type jsonPK struct {
	PkAsString string `json:"pk"`
}

func createPrivateKey() PrivateKey {
	x := curve.Scalar().Pick(curve.RandomStream())
	out := privateKey{
		x: x,
	}

	return &out
}

func createPrivateKeyFromString(str string) (PrivateKey, error) {
	x, xErr := fromStringToScalar(str)
	if xErr != nil {
		return nil, xErr
	}

	out := privateKey{
		x: x,
	}

	return &out, nil
}

// PublicKey returns the public key
func (app *privateKey) PublicKey() PublicKey {
	g := curve.Point().Base()
	p := curve.Point().Mul(app.x, g)
	return createPublicKey(p)
}

// RingSign signs a ring signature on the given message, in the given ring pubKey
func (app *privateKey) RingSign(msg string, ringPubKeys []PublicKey) (RingSignature, error) {

	retrieveSignerIndexFn := func(ringPubKeys []PublicKey, pk PrivateKey) int {
		pubKey := pk.PublicKey()
		for index, oneRingPubKey := range ringPubKeys {
			if oneRingPubKey.Equals(pubKey) {
				return index
			}
		}

		return -1
	}

	// retrieve our signerIndex:
	signerIndex := retrieveSignerIndexFn(ringPubKeys, app)
	if signerIndex == -1 {
		return nil, errors.New("the signer PublicKey is not in the ring")
	}

	// generate k:
	k := genK(app.x, msg)

	// random base:
	g := curve.Point().Base()

	// length:
	r := len(ringPubKeys)

	// initialize:
	es := make([]kyber.Scalar, r)
	ss := make([]kyber.Scalar, r)
	beginIndex := (signerIndex + 1) % r

	// ei = H(m || k * G)
	es[beginIndex] = hash(msg + curve.Point().Mul(k, g).String())

	// loop:
	for i := beginIndex; i != signerIndex; i = (i + 1) % r {
		// si = random value
		ss[i] = genK(app.x, fmt.Sprintf("%s%d", msg, i))

		//eiPlus1ModR = H(m || si * G + ei * Pi)
		sig := curve.Point().Mul(ss[i], g)
		eipi := curve.Point().Mul(es[i], ringPubKeys[i].Point())
		es[(i+1)%r] = hash(msg + curve.Point().Add(sig, eipi).String())

	}

	// close the ring:
	esx := curve.Scalar().Mul(es[signerIndex], app.x)
	ss[signerIndex] = curve.Scalar().Sub(k, esx)
	out := createRingSignature(ringPubKeys, ss, es[0])
	return out, nil
}

// Sign signs a message
func (app *privateKey) Sign(msg string) Signature {
	// generate k:
	k := genK(app.x, msg)

	// random base:
	g := curve.Point().Base()

	// r = k * G (a.k.a the same operation as r = g^k)
	r := curve.Point().Mul(k, g)

	// hash(m || r)
	e := hash(msg + r.String())

	// s = k - e * x
	s := curve.Scalar().Sub(k, curve.Scalar().Mul(e, app.x))

	// create signature:
	pubKey := createPublicKey(r)
	out := createSignature(pubKey, s)
	return out
}

// String returns the string representation of the PrivateKey
func (app *privateKey) String() string {
	return app.x.String()
}

// MarshalJSON returns a JSON representation of the object
func (app *privateKey) MarshalJSON() ([]byte, error) {
	return cdc.MarshalJSON(jsonPK{
		PkAsString: app.String(),
	})
}

// UnmarshalJSON returns an object based on the JSON data
func (app *privateKey) UnmarshalJSON(data []byte) error {
	ptr := new(jsonPK)
	err := cdc.UnmarshalJSON(data, ptr)
	if err != nil {
		return err
	}

	x, xErr := fromStringToScalar(ptr.PkAsString)
	if xErr != nil {
		return xErr
	}

	app.x = x
	return nil
}
