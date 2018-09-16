package crypto

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

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

func createRingSignatureFromString(str string) (RingSignature, error) {
	decoded, decodedErr := base64.StdEncoding.DecodeString(str)
	if decodedErr != nil {
		return nil, decodedErr
	}

	splitted := strings.Split(string(decoded), "|")
	if len(splitted) != 3 {
		return nil, errors.New("the decoded string was expected to have 3 parts")
	}

	ringAsStrings := strings.Split(splitted[0], "-")
	sAsStrings := strings.Split(splitted[1], "-")
	eAsString := splitted[2]

	rings := []kyber.Point{}
	for _, oneRingAsString := range ringAsStrings {
		decodedRing, decodedRingErr := hex.DecodeString(oneRingAsString)
		if decodedRingErr != nil {
			return nil, decodedRingErr
		}

		oneRing := curve.Point()
		decErr := oneRing.UnmarshalBinary(decodedRing)
		if decErr != nil {
			return nil, decErr
		}

		rings = append(rings, oneRing)
	}

	s := []kyber.Scalar{}
	for _, oneS := range sAsStrings {
		decodedS, decodedSErr := hex.DecodeString(oneS)
		if decodedSErr != nil {
			return nil, decodedSErr
		}

		oneS := curve.Scalar()
		decErr := oneS.UnmarshalBinary(decodedS)
		if decErr != nil {
			return nil, decErr
		}

		s = append(s, oneS)
	}

	decodedE, decodedEErr := hex.DecodeString(eAsString)
	if decodedEErr != nil {
		return nil, decodedEErr
	}

	e := curve.Scalar()
	decErr := e.UnmarshalBinary(decodedE)
	if decErr != nil {
		return nil, decErr
	}

	return createRingSignature(rings, s, e), nil
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

// String returns the string representation of the ring signature
func (app *ringSignature) String() string {
	ringAsStrings := []string{}
	for _, oneRing := range app.ring {
		ringAsStrings = append(ringAsStrings, oneRing.String())
	}

	sAsStrings := []string{}
	for _, oneS := range app.s {
		sAsStrings = append(sAsStrings, oneS.String())
	}

	ringAsString := strings.Join(ringAsStrings, "-")
	sAsString := strings.Join(sAsStrings, "-")
	eAsString := app.e.String()

	str := fmt.Sprintf("%s|%s|%s", ringAsString, sAsString, eAsString)
	return base64.StdEncoding.EncodeToString([]byte(str))
}
