package restapis

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type token struct {
	Mthod  string              `json:"method"`
	ReqURI string              `json:"request_uri"`
	Dta    map[string][]string `json:"data"`
}

func createToken(method string, reqURI string, data map[string][]string) Token {
	out := token{
		Mthod:  method,
		ReqURI: reqURI,
		Dta:    data,
	}

	return &out
}

func createTokenWithFormData(method string, reqURI string, formData map[string]string) Token {
	data := map[string][]string{}
	for keyname, element := range formData {
		data[keyname] = []string{
			element,
		}
	}

	return createToken(method, reqURI, data)
}

func createTokenWithoutData(method string, reqURI string) Token {
	return createToken(method, reqURI, map[string][]string{})
}

// Method returns the method
func (obj *token) Method() string {
	return obj.Mthod
}

// RequestURI returns the request URI
func (obj *token) RequestURI() string {
	return obj.ReqURI
}

// Data returns the Dta
func (obj *token) Data() map[string][]string {
	return obj.Dta
}

// Hash returns the hash
func (obj *token) Hash() string {
	js, jsErr := json.Marshal(obj)
	if jsErr != nil {
		panic(jsErr)
	}

	hasher := sha256.New()
	hasher.Write(js)
	return hex.EncodeToString(hasher.Sum(nil))
}
