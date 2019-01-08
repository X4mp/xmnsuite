package information

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

/*
 * Information
 */

type information struct {
	UUID                 *uuid.UUID `json:"id"`
	ConNeeded            int        `json:"concensus_needed"`
	GzPricePerKb         int        `json:"gaz_price_per_kb"`
	MxAmountOfValidators int        `json:"max_amount_of_validators"`
	NetShare             int        `json:"network_share"`
	ValShare             int        `json:"validator_share"`
	AffShare             int        `json:"affiliate_share"`
}

func createInformation(
	id *uuid.UUID,
	concensusNeeded int,
	gazPricePerKb int,
	maxAmountOfValidators int,
	netShare int,
	valShare int,
	affShare int,
) (Information, error) {

	sum := netShare + valShare + affShare
	if sum != 100 {
		str := fmt.Sprintf("the sum of the shares (network: %d, validators: %d, affiliate: %d) was expected to be 100, %d given", netShare, valShare, affShare, sum)
		return nil, errors.New(str)
	}

	if gazPricePerKb == 0 {
		return nil, errors.New("the gazPricePerKb cannot be 0")
	}

	if maxAmountOfValidators == 0 {
		return nil, errors.New("the maxAmountOfValidators cannot be 0")
	}

	if concensusNeeded == 0 {
		return nil, errors.New("the concensusNeeded cannot be 0")
	}

	out := information{
		UUID:                 id,
		ConNeeded:            concensusNeeded,
		GzPricePerKb:         gazPricePerKb,
		MxAmountOfValidators: maxAmountOfValidators,
		NetShare:             netShare,
		ValShare:             valShare,
		AffShare:             affShare,
	}

	return &out, nil
}

func createInformationFromNormalized(ins *normalizedInformation) (Information, error) {
	id, idErr := uuid.FromString(ins.ID)
	if idErr != nil {
		return nil, idErr
	}

	return createInformation(&id, ins.ConcensusNeeded, ins.GzPricePerKb, ins.MxAmountOfValidators, ins.NetworkShare, ins.ValidatorsShare, ins.AffiliateShare)
}

// ID returns the ID
func (app *information) ID() *uuid.UUID {
	return app.UUID
}

// GazPricePerKb returns the gazPricePerKb
func (app *information) GazPricePerKb() int {
	return app.GzPricePerKb
}

// ConcensusNeeded returns the concensusNeeded
func (app *information) ConcensusNeeded() int {
	return app.ConNeeded
}

// MaxAmountOfValidators returns the maxAmountOfValidators
func (app *information) MaxAmountOfValidators() int {
	return app.MxAmountOfValidators
}

// NetworkShare returns the network fees share
func (app *information) NetworkShare() int {
	return app.NetShare
}

// ValidatorsShare returns the validators fees share
func (app *information) ValidatorsShare() int {
	return app.ValShare
}

// AffiliateShare returns the affiliate fees share
func (app *information) AffiliateShare() int {
	return app.AffShare
}
