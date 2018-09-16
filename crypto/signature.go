package crypto

import (
	"github.com/dedis/kyber"
)

type signature struct {
	r kyber.Point
	s kyber.Scalar
}

func createSignature(r kyber.Point, s kyber.Scalar) *signature {
	out := signature{r: r, s: s}
	return &out
}

// PublicKey returns the public key of the signature
func (app *signature) PublicKey(msg string) kyber.Point {
	// Create a generator.
	g := curve.Point().Base()

	// e = Hash(m || r)
	e := hash(msg + app.r.String())

	// y = (r - s * G) * (1 / e)
	y := curve.Point().Sub(app.r, curve.Point().Mul(app.s, g))
	y = curve.Point().Mul(curve.Scalar().Div(curve.Scalar().One(), e), y)

	return y
}

// Verify verifies if the signature has been made by the given public key, on the message
func (app *signature) Verify(msg string, p kyber.Point) bool {
	// Create a generator.
	g := curve.Point().Base()

	// e = Hash(m || r)
	e := hash(msg + app.r.String())

	// Attempt to reconstruct 's * G' with a provided signature; s * G = r - e * p
	sGv := curve.Point().Sub(app.r, curve.Point().Mul(e, p))

	// Construct the actual 's * G'
	sG := curve.Point().Mul(app.s, g)

	// Equality check; ensure signature and public key outputs to s * G.
	return sG.Equal(sGv)
}
