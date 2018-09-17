package crypto

import (
	"encoding/base64"

	"github.com/dedis/kyber"
)

type signature struct {
	r PublicKey
	s kyber.Scalar
}

type jsonSignature struct {
	RAsString string `json:"r"`
	SAsString string `json:"s"`
}

func createSignature(r PublicKey, s kyber.Scalar) Signature {
	out := signature{r: r, s: s}
	return &out
}

func createSignatureFromString(sigAsString string) (Signature, error) {
	decoded, decodedErr := base64.StdEncoding.DecodeString(sigAsString)
	if decodedErr != nil {
		return nil, decodedErr
	}

	ptr := new(signature)
	err := cdc.UnmarshalJSON(decoded, ptr)
	if err != nil {
		return nil, err
	}

	return ptr, nil
}

// PublicKey returns the public key of the signature
func (app *signature) PublicKey(msg string) PublicKey {
	// Create a generator.
	g := curve.Point().Base()

	// e = Hash(m || r)
	e := hash(msg + app.r.String())

	// y = (r - s * G) * (1 / e)
	y := curve.Point().Sub(app.r.Point(), curve.Point().Mul(app.s, g))
	y = curve.Point().Mul(curve.Scalar().Div(curve.Scalar().One(), e), y)

	return createPublicKey(y)
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
	sGv := curve.Point().Sub(app.r.Point(), curve.Point().Mul(e, p.Point()))

	// Construct the actual 's * G'
	sG := curve.Point().Mul(app.s, g)

	// Equality check; ensure signature and public key outputs to s * G.
	return sG.Equal(sGv)
}

// String returns the string representation of the signature
func (app *signature) String() string {
	js, jsErr := cdc.MarshalJSON(app)
	if jsErr != nil {
		panic(jsErr)
	}

	return base64.StdEncoding.EncodeToString(js)
}

// MarshalJSON returns a JSON representation of the object
func (app *signature) MarshalJSON() ([]byte, error) {
	return cdc.MarshalJSON(jsonSignature{
		RAsString: app.r.String(),
		SAsString: app.s.String(),
	})
}

// UnmarshalJSON returns an object based on the JSON data
func (app *signature) UnmarshalJSON(data []byte) error {
	ptr := new(jsonSignature)
	err := cdc.UnmarshalJSON(data, ptr)
	if err != nil {
		return err
	}

	r, rErr := fromStringToPoint(ptr.RAsString)
	if rErr != nil {
		return rErr
	}

	s, sErr := fromStringToScalar(ptr.SAsString)
	if sErr != nil {
		return sErr
	}

	app.r = createPublicKey(r)
	app.s = s
	return nil
}
