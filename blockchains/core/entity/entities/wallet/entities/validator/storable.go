package validator

import "encoding/hex"

type storableValidator struct {
	ID       string `json:"id"`
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	PubKey   string `json:"pubkey"`
	PledgeID string `json:"pledge_id"`
}

func createStorableValidator(ins Validator) *storableValidator {
	out := storableValidator{
		ID:       ins.ID().String(),
		IP:       ins.IP().String(),
		PubKey:   hex.EncodeToString(ins.PubKey().Bytes()),
		PledgeID: ins.Pledge().ID().String(),
	}

	return &out
}
