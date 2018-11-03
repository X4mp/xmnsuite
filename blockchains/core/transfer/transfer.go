package transfer

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/token"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/framework/wallet"
	"github.com/xmnservices/xmnsuite/crypto"
)

type transfer struct {
	UUID *uuid.UUID       `json:"id"`
	Frm  wallet.Wallet    `json:"from"`
	Tok  token.Token      `json:"token"`
	Am   int              `json:"amount"`
	Cnt  string           `json:"content"`
	PKey crypto.PublicKey `json:"public_key"`
}

func createTransfer(id *uuid.UUID, frm wallet.Wallet, tok token.Token, amount int, content string, pubKey crypto.PublicKey) Transfer {
	out := transfer{
		UUID: id,
		Frm:  frm,
		Am:   amount,
		Cnt:  content,
		PKey: pubKey,
	}

	return &out
}

// ID returns the ID
func (obj *transfer) ID() *uuid.UUID {
	return obj.UUID
}

// From returns the from wallet
func (obj *transfer) From() wallet.Wallet {
	return obj.Frm
}

// Token returns the token
func (obj *transfer) Token() token.Token {
	return obj.Tok
}

// Amount returns the amount
func (obj *transfer) Amount() int {
	return obj.Am
}

// Content returns the content
func (obj *transfer) Content() string {
	return obj.Cnt
}

// PubKey returns the public key
func (obj *transfer) PubKey() crypto.PublicKey {
	return obj.PKey
}
