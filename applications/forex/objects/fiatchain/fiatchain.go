package fiatchain

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/deposit"
)

type fiatChain struct {
	UUID *uuid.UUID        `json:"id"`
	Sds  []string          `json:"seeds"`
	Deps []deposit.Deposit `json:"deposits"`
}

func createFiatChain(id *uuid.UUID, seeds []string, deps []deposit.Deposit) (FiatChain, error) {
	out := fiatChain{
		UUID: id,
		Sds:  seeds,
		Deps: deps,
	}

	return &out, nil
}

func fromNormalizedToFiatChain(normalized *normalizedFiatChain) (FiatChain, error) {
	id, idErr := uuid.FromString(normalized.ID)
	if idErr != nil {
		return nil, idErr
	}

	deps := []deposit.Deposit{}
	depMetaData := deposit.SDKFunc.CreateMetaData()
	for _, oneNormalized := range normalized.Deposits {
		oneDepIns, oneDepInsErr := depMetaData.Denormalize()(oneNormalized)
		if oneDepInsErr != nil {
			return nil, oneDepInsErr
		}

		if oneDep, ok := oneDepIns.(deposit.Deposit); ok {
			deps = append(deps, oneDep)
			continue
		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid Deposit instance", oneDepIns.ID().String())
		return nil, errors.New(str)
	}

	return createFiatChain(&id, normalized.Seeds, deps)
}

func fromStorableToFiatChain(storable *storableFiatChain, rep entity.Repository) (FiatChain, error) {
	id, idErr := uuid.FromString(storable.ID)
	if idErr != nil {
		return nil, idErr
	}

	deps := []deposit.Deposit{}
	depMetaData := deposit.SDKFunc.CreateMetaData()
	for _, oneDepIDAsString := range storable.DepositIDs {
		oneDepID, oneDepIDErr := uuid.FromString(oneDepIDAsString)
		if oneDepIDErr != nil {
			return nil, oneDepIDErr
		}

		oneDepIns, oneDepInsErr := rep.RetrieveByID(depMetaData, &oneDepID)
		if oneDepInsErr != nil {
			return nil, oneDepInsErr
		}

		if oneDep, ok := oneDepIns.(deposit.Deposit); ok {
			deps = append(deps, oneDep)
			continue
		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid Deposit instance", oneDepIns.ID().String())
		return nil, errors.New(str)
	}

	return createFiatChain(&id, storable.Seeds, deps)
}

// ID returns the ID
func (obj *fiatChain) ID() *uuid.UUID {
	return obj.UUID
}

// Seeds returns the seeds
func (obj *fiatChain) Seeds() []string {
	return obj.Sds
}

// Deposits returns the deposits
func (obj *fiatChain) Deposits() []deposit.Deposit {
	return obj.Deps
}
