package transfer

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/withdrawal"
	"github.com/xmnservices/xmnsuite/crypto"
)

type transfer struct {
	UUID   *uuid.UUID            `json:"id"`
	Withdr withdrawal.Withdrawal `json:"withdrawal"`
	Cnt    string                `json:"content"`
	PKey   crypto.PublicKey      `json:"public_key"`
}

func createTransfer(id *uuid.UUID, withdrawal withdrawal.Withdrawal, content string, pubKey crypto.PublicKey) Transfer {
	out := transfer{
		UUID:   id,
		Withdr: withdrawal,
		Cnt:    content,
		PKey:   pubKey,
	}

	return &out
}

// ID returns the ID
func (obj *transfer) ID() *uuid.UUID {
	return obj.UUID
}

// Withdrawal returns the from withdrawal
func (obj *transfer) Withdrawal() withdrawal.Withdrawal {
	return obj.Withdr
}

// Content returns the content
func (obj *transfer) Content() string {
	return obj.Cnt
}

// PubKey returns the public key
func (obj *transfer) PubKey() crypto.PublicKey {
	return obj.PKey
}
