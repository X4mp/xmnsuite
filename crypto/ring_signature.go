package crypto

import (
	"encoding/base64"

	"github.com/dedis/kyber"
)

type ringSignature struct {
	ring []PublicKey
	s    []kyber.Scalar
	e    kyber.Scalar
}

type jsonRingSignature struct {
	RingAsStrings []string `json:"ring"`
	SAsStrings    []string `json:"s"`
	EAsString     string   `json:"e"`
}

func createRingSignature(ring []PublicKey, s []kyber.Scalar, e kyber.Scalar) RingSignature {
	out := ringSignature{
		ring: ring,
		s:    s,
		e:    e,
	}

	return &out
}

func createRingSignatureFromString(str string) (RingSignature, error) {
	decoded, decodedErr := base64.StdEncoding.DecodeString(str)
	if decodedErr != nil {
		return nil, decodedErr
	}

	ptr := new(ringSignature)
	err := cdc.UnmarshalJSON(decoded, ptr)
	if err != nil {
		return nil, err
	}

	return ptr, nil
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
		ep := curve.Point().Mul(e, app.ring[i].Point())
		added := curve.Point().Add(sg, ep)
		e = hash(msg + added.String())
	}

	return app.e.Equal(e)
}

// String returns the string representation of the ring signature
func (app *ringSignature) String() string {
	js, jsErr := cdc.MarshalJSON(app)
	if jsErr != nil {
		panic(jsErr)
	}

	return base64.StdEncoding.EncodeToString(js)
}

// MarshalJSON returns a JSON representation of the object
func (app *ringSignature) MarshalJSON() ([]byte, error) {
	rings := []string{}
	for _, oneRing := range app.ring {
		rings = append(rings, oneRing.String())
	}

	s := []string{}
	for _, oneS := range app.s {
		s = append(s, oneS.String())
	}

	return cdc.MarshalJSON(jsonRingSignature{
		RingAsStrings: rings,
		SAsStrings:    s,
		EAsString:     app.e.String(),
	})
}

// UnmarshalJSON returns an object based on the JSON data
func (app *ringSignature) UnmarshalJSON(data []byte) error {
	ptr := new(jsonRingSignature)
	err := cdc.UnmarshalJSON(data, ptr)
	if err != nil {
		return err
	}

	rings := []PublicKey{}
	for _, oneRingAsString := range ptr.RingAsStrings {
		p, pErr := fromStringToPoint(oneRingAsString)
		if pErr != nil {
			return pErr
		}

		rings = append(rings, createPublicKey(p))
	}

	ss := []kyber.Scalar{}
	for _, oneS := range ptr.SAsStrings {
		s, sErr := fromStringToScalar(oneS)
		if sErr != nil {
			return sErr
		}

		ss = append(ss, s)
	}

	e, eErr := fromStringToScalar(ptr.EAsString)
	if eErr != nil {
		return eErr
	}
	app.ring = rings
	app.s = ss
	app.e = e
	return nil
}
