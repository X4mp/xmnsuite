package restapis

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	crypto "github.com/xmnservices/xmnsuite/crypto"
)

type account struct {
	Nme  string    `json:"name"`
	PK   string    `json:"encrypted_pk"`
	CrOn time.Time `json:"created_on"`
}

func createAccount(name string, pk string, crOn time.Time) (Account, error) {

	pattern, patternErr := regexp.Compile("[a-zA-Z0-9-]+")
	if patternErr != nil {
		return nil, patternErr
	}

	if len(name) < 2 {
		return nil, errors.New("the account name must have at least 2 characters")
	}

	if pattern.FindString(name) != name {
		str := fmt.Sprintf("the account name (%s) must only contain letters, numbers or hyphens (-)", name)
		return nil, errors.New(str)
	}

	out := account{
		Nme:  name,
		PK:   pk,
		CrOn: crOn,
	}

	return &out, nil
}

// Name returns the name
func (obj *account) Name() string {
	return obj.Nme
}

// EncryptedPK returns the encrypted PK
func (obj *account) EncryptedPK() string {
	return obj.PK
}

// CreatedOn returns the creation time
func (obj *account) CreatedOn() time.Time {
	return obj.CrOn
}

// Name returns the name
func (obj *account) DecryptPK(seedWords []string) crypto.PrivateKey {
	return crypto.SDKFunc.CreatePK(crypto.CreatePKParams{
		PKAsString: string(crypto.SDKFunc.Decrypt(crypto.DecryptParams{
			Pass:         []byte(strings.Join(seedWords, "|")),
			EncryptedMsg: obj.PK,
		})),
	})
}
