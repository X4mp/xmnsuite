package crypto

import (
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"

	"github.com/dedis/kyber"
)

type signature struct {
	r kyber.Point
	s kyber.Scalar
}

func createSignature(r kyber.Point, s kyber.Scalar) Signature {
	out := signature{r: r, s: s}
	return &out
}

func createSignatureFromString(sigAsString string) (Signature, error) {
	pattern, patternErr := regexp.Compile("([0-9a-f]+)")
	if patternErr != nil {
		return nil, patternErr
	}

	strs := pattern.FindAllString(sigAsString, -1)
	if len(strs) != 2 {
		str := fmt.Sprintf("the given string (%s) us not a valid signature", sigAsString)
		return nil, errors.New(str)
	}

	decodedR, decodedRErr := hex.DecodeString(strs[0])
	if decodedRErr != nil {
		return nil, decodedRErr
	}

	decodedS, decodedSErr := hex.DecodeString(strs[1])
	if decodedSErr != nil {
		return nil, decodedSErr
	}

	r := curve.Point()
	rErr := r.UnmarshalBinary(decodedR)
	if rErr != nil {
		return nil, rErr
	}

	s := curve.Scalar()
	sErr := s.UnmarshalBinary(decodedS)
	if sErr != nil {
		return nil, sErr
	}

	return createSignature(r, s), nil
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
func (app *signature) Verify(msg string) bool {

	// retrieve pubKey:
	p := app.PublicKey(msg)

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

// String returns the string representation of the signature
func (app *signature) String() string {
	rAsString := app.r.String()
	sAsString := app.s.String()
	return fmt.Sprintf("%s-%s", rAsString, sAsString)
}
