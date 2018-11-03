package xmn

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/crypto"
)

type transferClaim struct {
	UUID *uuid.UUID           `json:"id"`
	To   Wallet               `json:"to"`
	Sig  crypto.RingSignature `json:"signed_content"`
	Am   int                  `json:"amount"`
}

func createTransferClaim(id *uuid.UUID, to Wallet, sig crypto.RingSignature, amount int) TransferClaim {
	out := transferClaim{
		UUID: id,
		To:   to,
		Sig:  sig,
		Am:   amount,
	}

	return &out
}

// ID returns the ID
func (obj *transferClaim) ID() *uuid.UUID {
	return obj.UUID
}

// DepositTo returns the Wallet to deposit the claim to
func (obj *transferClaim) DepositTo() Wallet {
	return obj.To
}

// SignedContent represents the signed content signature
func (obj *transferClaim) SignedContent() crypto.RingSignature {
	return obj.Sig
}

// Amount represents the amount
func (obj *transferClaim) Amount() int {
	return obj.Am
}
