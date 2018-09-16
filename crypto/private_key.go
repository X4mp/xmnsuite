package crypto

import (
	"fmt"

	"github.com/dedis/kyber"
)

type privateKey struct {
	x kyber.Scalar
}

func createPrivateKey() *privateKey {
	x := curve.Scalar().Pick(curve.RandomStream())
	out := privateKey{
		x: x,
	}

	return &out
}

// PublicKey returns the public key
func (app *privateKey) PublicKey() kyber.Point {
	g := curve.Point().Base()
	return curve.Point().Mul(app.x, g)
}

// RingSign signs a ring signature on the given message, in the given ring pubKey
func (app *privateKey) RingSign(msg string, ringPubKeys []kyber.Point, signerIndex int) *ringSignature {
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
		eipi := curve.Point().Mul(es[i], ringPubKeys[i])
		es[(i+1)%r] = hash(msg + curve.Point().Add(sig, eipi).String())

	}

	// close the ring:
	esx := curve.Scalar().Mul(es[signerIndex], app.x)
	ss[signerIndex] = curve.Scalar().Sub(k, esx)
	out := ringSignature{ring: ringPubKeys, e: es[0], s: ss}
	return &out
}

// Sign signs a message
func (app *privateKey) Sign(msg string) *signature {
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
	out := createSignature(r, s)
	return out
}
