package crypto

import (
	"github.com/dedis/kyber"
)

type ringSignature struct {
	ring []kyber.Point
	s    []kyber.Scalar
	e    kyber.Scalar
}

func createRingSignature(ring []kyber.Point, s []kyber.Scalar, e kyber.Scalar) RingSignature {
	out := ringSignature{
		ring: ring,
		s:    s,
		e:    e,
	}

	return &out
}

// Verify verifies if the message has been signed by at least 1 shared signature
func (app *ringSignature) Verify(msg string) bool {
	// random base:
	g := curve.Point().Base()

	// first e:
	e := app.e

	//e = H(m || s[i] * G + e * P[i]);
	amount := len(app.ring)
	for i := 0; i < amount; i++ {
		sg := curve.Point().Mul(app.s[i], g)
		ep := curve.Point().Mul(e, app.ring[i])
		added := curve.Point().Add(sg, ep)
		e = hash(msg + added.String())
	}

	return app.e.Equal(e)
}
