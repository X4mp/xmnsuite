package crypto

/*
 * H'(m, s, e) = H(m || s * G + e * P)
 * P = x * G
 * e = H(m || k * G)
 * k = s + e * x
 * s = k – e * x
 * k = H(m || x) -> to generate a new k, since nobody but us knows x
 * where ...
 * 1. H is a hash function, for instance SHA256.
 * 2. s and e are 2 numbers forming the ring signature
 * 3. s and r are a pubKey and a number forming a signature
 * 4. m is the message we want to sign
 * 5. P is the public key.
 * 6. G is the random base
 * 7. k is a number chosen randomly.  A new one every time we sign must be generated
 * 8. x is the private key
 */

import (
	"bytes"
	"encoding/hex"

	"github.com/dedis/kyber"
	"github.com/dedis/kyber/group/edwards25519"
)

var curve = edwards25519.NewBlakeSHA256Ed25519()

func hash(msg string) kyber.Scalar {
	sha256 := curve.Hash()
	sha256.Reset()
	sha256.Write([]byte(msg))

	return curve.Scalar().SetBytes(sha256.Sum(nil))
}

func genK(x kyber.Scalar, msg string) kyber.Scalar {
	return hash(msg + x.String())
}

func fromStringToScalar(str string) (kyber.Scalar, error) {
	decoded, decodedErr := hex.DecodeString(str)
	if decodedErr != nil {
		return nil, decodedErr
	}

	x := curve.Scalar()
	reader := bytes.NewReader(decoded)
	_, err := x.UnmarshalFrom(reader)
	if err != nil {
		return nil, err
	}

	return x, nil
}

func fromStringToPoint(str string) (kyber.Point, error) {
	decoded, decodedErr := hex.DecodeString(str)
	if decodedErr != nil {
		return nil, decodedErr
	}

	p := curve.Point()
	err := p.UnmarshalBinary(decoded)
	if err != nil {
		return nil, err
	}

	return p, nil
}
