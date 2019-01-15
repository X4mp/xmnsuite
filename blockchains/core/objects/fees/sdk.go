package fees

import (
	"errors"
	"fmt"
	"unsafe"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/genesis"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/affiliates"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/validator"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/information"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/withdrawal"
	"github.com/xmnservices/xmnsuite/datastore"
)

// Fee represents a transaction fee
type Fee interface {
	ID() *uuid.UUID
	Client() withdrawal.Withdrawal
	Network() deposit.Deposit
	Validators() []deposit.Deposit
	HasAffiliate() bool
	Affiliate() deposit.Deposit
}

// Normalized represents a normalized fee
type Normalized interface {
}

// CreateParams represents the create params
type CreateParams struct {
	ID         *uuid.UUID
	Gen        genesis.Genesis
	StoredData []byte
	Client     user.User
	Affiliate  affiliates.Affiliate
	Validators []validator.Validator
}

// SDKFunc represents the fees SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Fee
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
}{
	Create: func(params CreateParams) Fee {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		createPrices := func(aff affiliates.Affiliate, info information.Information, price int) (int, int, int) {
			if aff == nil {
				affShares := info.AffiliateShare()
				splitAffShares := int(affShares / 2)
				networkShares := (info.NetworkShare() + splitAffShares)
				validatorsShare := (info.ValidatorsShare() + splitAffShares)
				networkPrice := int((price * networkShares) / 100)
				validatorsPrice := int((price * validatorsShare) / 100)
				return networkPrice, validatorsPrice, 0
			}

			affiliatePrice := int((price * info.AffiliateShare()) / 100)
			networkPrice := int((price * info.NetworkShare()) / 100)
			validatorsPrice := int((price * info.ValidatorsShare()) / 100)
			return networkPrice, validatorsPrice, affiliatePrice
		}

		createValidatorsDeposits := func(vals []validator.Validator, totalPrice int) ([]deposit.Deposit, int) {
			totalPledge := 0
			for _, oneVal := range vals {
				totalPledge += oneVal.Pledge().From().Amount()
			}

			totalPaid := 0
			deps := []deposit.Deposit{}
			for _, oneVal := range vals {
				amount := int((oneVal.Pledge().From().Amount() / totalPledge) * totalPrice)
				deps = append(deps, deposit.SDKFunc.Create(deposit.CreateParams{
					To:     oneVal.Pledge().To(),
					Amount: amount,
				}))

				totalPaid += amount
			}

			return deps, totalPaid
		}

		// create the price to pay:
		price := int(unsafe.Sizeof(params.StoredData)) * params.Gen.Info().GazPricePerKb()
		networkPrice, validatorsPrice, affiliatePrice := createPrices(params.Affiliate, params.Gen.Info(), price)

		// create the validators deposits:
		validatorsDeps, totalPaidToVals := createValidatorsDeposits(params.Validators, validatorsPrice)

		// create the network deposit:
		networkDeposit := deposit.SDKFunc.Create(deposit.CreateParams{
			To:     params.Gen.Deposit().To(),
			Amount: networkPrice,
		})

		// create the client withdrawal:
		client := withdrawal.SDKFunc.Create(withdrawal.CreateParams{
			From:   params.Client.Wallet(),
			Amount: networkPrice + totalPaidToVals + affiliatePrice,
		})

		// if there is no affiliate:
		if params.Affiliate == nil {
			out, outErr := createFee(params.ID, client, networkDeposit, validatorsDeps)
			if outErr != nil {
				panic(outErr)
			}

			return out
		}

		// crate the affiliate deposit:
		affDeposit := deposit.SDKFunc.Create(deposit.CreateParams{
			To:     params.Affiliate.Owner(),
			Amount: affiliatePrice,
		})

		out, outErr := createFeeWithAffiliate(params.ID, client, networkDeposit, validatorsDeps, affDeposit)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	CreateMetaData: func() entity.MetaData {
		return createMetaData()
	},
	CreateRepresentation: func() entity.Representation {
		return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
			Met: createMetaData(),
			ToStorable: func(ins entity.Entity) (interface{}, error) {
				if fe, ok := ins.(Fee); ok {
					out := createStorableFee(fe)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Fee instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				if _, ok := ins.(Fee); ok {
					return []string{
						retrieveAllFeesKeyname(),
					}, nil
				}

				return nil, errors.New("the given entity is not a valid Fee instance")
			},
			OnSave: func(ds datastore.DataStore, ins entity.Entity) error {
				if fe, ok := ins.(Fee); ok {
					// metadata:
					metaData := createMetaData()
					withdrawalMetaData := withdrawal.SDKFunc.CreateMetaData()
					depositMetaData := deposit.SDKFunc.CreateMetaData()

					// create the repository and service:
					repository := entity.SDKFunc.CreateRepository(ds)

					// make sure the fee doesnt already exists:
					_, retFee := repository.RetrieveByID(metaData, fe.ID())
					if retFee == nil {
						str := fmt.Sprintf("the fee (ID: %s) already exists", fe.ID().String())
						return errors.New(str)
					}

					// make sure the client withdrawal does not exists:
					_, retClientErr := repository.RetrieveByID(withdrawalMetaData, fe.Client().ID())
					if retClientErr == nil {
						str := fmt.Sprintf("the client withdrawal (ID: %s) inside the fee (ID: %s) already exists", fe.Client().ID().String(), fe.ID().String())
						return errors.New(str)
					}

					// make sure the network deposit does not exists:
					_, retNetworkErr := repository.RetrieveByID(depositMetaData, fe.Network().ID())
					if retNetworkErr == nil {
						str := fmt.Sprintf("the network deposit (ID: %s) inside the fee (ID: %s) already exists", fe.Network().ID().String(), fe.ID().String())
						return errors.New(str)
					}

					if fe.HasAffiliate() {
						// make sure the affiliate deposit does not exists:
						_, retAffErr := repository.RetrieveByID(depositMetaData, fe.Affiliate().ID())
						if retAffErr == nil {
							str := fmt.Sprintf("the affiliate deposit (ID: %s) inside the fee (ID: %s) already exists", fe.Affiliate().ID().String(), fe.ID().String())
							return errors.New(str)
						}
					}

					validators := fe.Validators()
					for _, oneValidator := range validators {
						// make sure the validator deposit does not exists:
						_, retValidatorErr := repository.RetrieveByID(depositMetaData, oneValidator.ID())
						if retValidatorErr == nil {
							str := fmt.Sprintf("the validator deposit (ID: %s) inside the fee (ID: %s) already exists", oneValidator.ID().String(), fe.ID().String())
							return errors.New(str)
						}
					}

					return nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Fee instance", ins.ID().String())
				return errors.New(str)
			},
		})
	},
}
