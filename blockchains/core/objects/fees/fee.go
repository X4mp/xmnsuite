package fees

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/withdrawal"
)

type fee struct {
	UUID  *uuid.UUID            `json:"id"`
	Cl    withdrawal.Withdrawal `json:"client"`
	Nwork deposit.Deposit       `json:"network"`
	Vals  []deposit.Deposit     `json:"validators"`
	Aff   deposit.Deposit       `json:"affiliate"`
}

func createFee(id *uuid.UUID, cl withdrawal.Withdrawal, net deposit.Deposit, vals []deposit.Deposit) (Fee, error) {
	return createFeeWithAffiliate(id, cl, net, vals, nil)
}

func createFeeWithAffiliate(
	id *uuid.UUID,
	cl withdrawal.Withdrawal,
	net deposit.Deposit,
	vals []deposit.Deposit,
	aff deposit.Deposit,
) (Fee, error) {
	out := fee{
		UUID:  id,
		Cl:    cl,
		Nwork: net,
		Vals:  vals,
		Aff:   aff,
	}

	return &out, nil
}

func createFeeFromNormalized(normalized *normalizedFee) (Fee, error) {
	id, idErr := uuid.FromString(normalized.ID)
	if idErr != nil {
		return nil, idErr
	}

	clientIns, clientInsErr := withdrawal.SDKFunc.CreateMetaData().Denormalize()(normalized.Client)
	if clientInsErr != nil {
		return nil, clientInsErr
	}

	depositDenFunc := deposit.SDKFunc.CreateMetaData().Denormalize()
	networkIns, networkInsErr := depositDenFunc(normalized.Network)
	if networkInsErr != nil {
		return nil, networkInsErr
	}

	validators := []deposit.Deposit{}
	for _, oneValidator := range normalized.Validators {
		denValidator, denValidatorErr := depositDenFunc(oneValidator)
		if denValidatorErr != nil {
			return nil, denValidatorErr
		}

		if dep, ok := denValidator.(deposit.Deposit); ok {
			validators = append(validators, dep)
			continue
		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid Deposit instance", denValidator.ID().String())
		return nil, errors.New(str)

	}

	var affIns entity.Entity
	if normalized.Affiliate != nil {
		affDenIns, affDenInsErr := depositDenFunc(normalized.Affiliate)
		if affDenInsErr != nil {
			return nil, affDenInsErr
		}

		affIns = affDenIns
	}

	if client, ok := clientIns.(withdrawal.Withdrawal); ok {
		if network, ok := networkIns.(deposit.Deposit); ok {
			if affIns == nil {
				return createFee(&id, client, network, validators)
			}

			if aff, ok := affIns.(deposit.Deposit); ok {
				return createFeeWithAffiliate(&id, client, network, validators, aff)
			}

			str := fmt.Sprintf("the entity (ID: %s) is not a valid Deposit instance", affIns.ID().String())
			return nil, errors.New(str)
		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid Deposit instance", networkIns.ID().String())
		return nil, errors.New(str)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Withdrawal instance", clientIns.ID().String())
	return nil, errors.New(str)

}

func createFeeFromStorable(storable *storableFee, rep entity.Repository) (Fee, error) {
	// metadata:
	withMet := withdrawal.SDKFunc.CreateMetaData()
	depMet := deposit.SDKFunc.CreateMetaData()

	id, idErr := uuid.FromString(storable.ID)
	if idErr != nil {
		return nil, idErr
	}

	clientID, clientIDErr := uuid.FromString(storable.ClientID)
	if clientIDErr != nil {
		return nil, clientIDErr
	}

	networkID, networkIDErr := uuid.FromString(storable.NetworkID)
	if networkIDErr != nil {
		return nil, networkIDErr
	}

	validators := []deposit.Deposit{}
	for _, oneValidatorIDAsString := range storable.ValidatorsIDs {
		validatorID, validatorIDErr := uuid.FromString(oneValidatorIDAsString)
		if validatorIDErr != nil {
			return nil, validatorIDErr
		}

		validatorIns, validatorInsErr := rep.RetrieveByID(depMet, &validatorID)
		if validatorInsErr != nil {
			return nil, validatorInsErr
		}

		if validator, ok := validatorIns.(deposit.Deposit); ok {
			validators = append(validators, validator)
			continue
		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid Deposit instance", validatorIns.ID().String())
		return nil, errors.New(str)
	}

	var affIns entity.Entity
	if storable.AffiliateID != "" {
		affID, affIDErr := uuid.FromString(storable.AffiliateID)
		if affIDErr != nil {
			return nil, affIDErr
		}

		retAffIns, retAffInsErr := rep.RetrieveByID(depMet, &affID)
		if retAffInsErr != nil {
			return nil, retAffInsErr
		}

		affIns = retAffIns
	}

	clientIns, clientInsErr := rep.RetrieveByID(withMet, &clientID)
	if clientInsErr != nil {
		return nil, clientInsErr
	}

	networkIns, networkInsErr := rep.RetrieveByID(depMet, &networkID)
	if networkInsErr != nil {
		return nil, networkInsErr
	}

	if client, ok := clientIns.(withdrawal.Withdrawal); ok {
		if network, ok := networkIns.(deposit.Deposit); ok {
			if affIns == nil {
				return createFee(&id, client, network, validators)
			}

			if aff, ok := affIns.(deposit.Deposit); ok {
				return createFeeWithAffiliate(&id, client, network, validators, aff)
			}

			str := fmt.Sprintf("the entity (ID: %s) is not a valid Deposit instance", affIns.ID().String())
			return nil, errors.New(str)
		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid Deposit instance", networkIns.ID().String())
		return nil, errors.New(str)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Withdrawal instance", clientIns.ID().String())
	return nil, errors.New(str)
}

// ID returns the ID
func (obj *fee) ID() *uuid.UUID {
	return obj.UUID
}

// Client returns the client
func (obj *fee) Client() withdrawal.Withdrawal {
	return obj.Cl
}

// Network returns the network
func (obj *fee) Network() deposit.Deposit {
	return obj.Nwork
}

// Validators returns the validators
func (obj *fee) Validators() []deposit.Deposit {
	return obj.Vals
}

// HasAffiliate returns true if there is an affiliate fee, false otherwise
func (obj *fee) HasAffiliate() bool {
	return obj.Aff != nil
}

// Affiliate returns the affiliate deposit, if any
func (obj *fee) Affiliate() deposit.Deposit {
	return obj.Aff
}
