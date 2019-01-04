package information

import (
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
}

func createInformation(
	id *uuid.UUID,
	concensusNeeded int,
	gazPricePerKb int,
	maxAmountOfValidators int,
) (Information, error) {

	out := information{
		UUID:                 id,
		ConNeeded:            concensusNeeded,
		GzPricePerKb:         gazPricePerKb,
		MxAmountOfValidators: maxAmountOfValidators,
	}

	return &out, nil
}

func createInformationFromNormalized(ins *normalizedInformation) (Information, error) {
	id, idErr := uuid.FromString(ins.ID)
	if idErr != nil {
		return nil, idErr
	}

	return createInformation(&id, ins.ConcensusNeeded, ins.GzPricePerKb, ins.MxAmountOfValidators)
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
