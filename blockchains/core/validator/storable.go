package validator

type storableValidator struct {
	ID       string `json:"id"`
	PubKey   string `json:"pubkey"`
	PledgeID string `json:"pledge_id"`
}

func createStorableValidator(ins Validator) *storableValidator {
	out := storableValidator{
		ID:       ins.ID().String(),
		PubKey:   string(ins.PubKey().Bytes()),
		PledgeID: ins.Pledge().ID().String(),
	}

	return &out
}
