package fees

type storableFee struct {
	ID          string   `json:"id"`
	ClientID    string   `json:"client"`
	NetworkID   string   `json:"network"`
	ValidatorsIDs    []string `json:"validators"`
	AffiliateID string   `json:"affiliate"`
}

func createStorableFee(ins Fee) *storableFee {

	validators := ins.Validators()
	validatorIDs := []string{}
	for _, oneValidator := range validators {
		validatorIDs = append(validatorIDs, oneValidator.ID().String())
	}

	affIDAsString := ""
	if ins.HasAffiliate() {
		affIDAsString = ins.Affiliate().ID().String()
	}

	out := storableFee{
		ID:          ins.ID().String(),
		ClientID:    ins.Client().ID().String(),
		NetworkID:   ins.Network().ID().String(),
		ValidatorsIDs:    validatorIDs,
		AffiliateID: affIDAsString,
	}

	return &out
}
