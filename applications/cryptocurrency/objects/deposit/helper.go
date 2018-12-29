package deposit

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/objects/address"
	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/objects/offer"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/datastore"
)

func retrieveAllDepositsKeyname() string {
	return "cryptodeposits"
}

func retrieveDepositByOfferKeyname(off offer.Offer) string {
	base := retrieveAllDepositsKeyname()
	return fmt.Sprintf("%s:by_offer_id:%s", base, off.ID().String())
}

func retrieveDepositByFromAddressKeyname(frmAddress address.Address) string {
	base := retrieveAllDepositsKeyname()
	return fmt.Sprintf("%s:by_from_address_id:%s", base, frmAddress.ID().String())
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "CryptoDeposit",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			if storable, ok := data.(*storableDeposit); ok {
				return createDepositFromStorable(storable, rep)
			}

			ptr := new(normalizedDeposit)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return createDepositFromNormalized(ptr)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if dep, ok := ins.(Deposit); ok {
				out, outErr := createNormalizedDeposit(dep)
				if outErr != nil {
					return nil, outErr
				}

				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Deposit instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*normalizedDeposit); ok {
				return createDepositFromNormalized(normalized)
			}

			return nil, errors.New("the given instance is not a valid normalized Deposit instance")
		},
		EmptyNormalized: new(normalizedDeposit),
		EmptyStorable:   new(storableDeposit),
	})
}

func createRepresentation() entity.Representation {
	return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
		Met: createMetaData(),
		ToStorable: func(ins entity.Entity) (interface{}, error) {
			if dep, ok := ins.(Deposit); ok {
				out := createStorableDeposit(dep)
				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Deposit instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Keynames: func(ins entity.Entity) ([]string, error) {
			if dep, ok := ins.(Deposit); ok {
				return []string{
					retrieveAllDepositsKeyname(),
					retrieveDepositByOfferKeyname(dep.Offer()),
					retrieveDepositByFromAddressKeyname(dep.From()),
				}, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Deposit instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Sync: func(ds datastore.DataStore, ins entity.Entity) error {

			// create the representations:
			metaData := createMetaData()

			// create the entity repository and service:
			entityRepository := entity.SDKFunc.CreateRepository(ds)
			depositRepository := createRepository(metaData, entityRepository)
			addressRepository := address.SDKFunc.CreateRepository(address.CreateRepositoryParams{
				EntityRepository: entityRepository,
			})

			offerRepository := offer.SDKFunc.CreateRepository(offer.CreateRepositoryParams{
				EntityRepository: entityRepository,
			})

			if dep, ok := ins.(Deposit); ok {
				// make sure the offer exists:
				_, retOffErr := offerRepository.RetrieveByID(dep.Offer().ID())
				if retOffErr != nil {
					str := fmt.Sprintf("the Offer (ID: %s) in the Deposit instance does not exists", dep.Offer().ID().String())
					return errors.New(str)
				}

				// make sure the address exists:
				_, retAddrErr := addressRepository.RetrieveByID(dep.From().ID())
				if retAddrErr != nil {
					str := fmt.Sprintf("the from Address (ID: %s) in the Deposit instance does not exists", dep.From().ID().String())
					return errors.New(str)
				}

				// retrieve all deposits of the offer:
				depPS, depPSErr := depositRepository.RetrieveSetByOffer(dep.Offer(), 0, -1)
				if depPSErr != nil {
					str := fmt.Sprintf("there was an error while retrieving a Deposit partial set by Offer (ID: %s): %s", dep.Offer().ID().String(), depPSErr.Error())
					return errors.New(str)
				}

				// retrieve the amount deposited so far:
				amountDeposited := 0
				depsIns := depPS.Instances()
				for _, oneDepIns := range depsIns {
					if oneDep, ok := oneDepIns.(Deposit); ok {
						amountDeposited += oneDep.Amount()
					}
				}

				// make sure the new deposit is not bigger than the remaining:
				remainingDeposit := dep.Offer().Amount() - amountDeposited
				if remainingDeposit < dep.Amount() {
					str := fmt.Sprintf("the deposit (%d) is bigger than the remaining potential deposit (%d) because there is already %d deposited on a %d offer", dep.Offer().Amount(), remainingDeposit, amountDeposited, dep.Offer().Amount())
					return errors.New(str)
				}

				return nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Offer instance", ins.ID().String())
			return errors.New(str)
		},
	})
}

func toData(dep Deposit) *Data {
	out := Data{
		ID:     dep.ID().String(),
		Offer:  offer.SDKFunc.ToData(dep.Offer()),
		From:   address.SDKFunc.ToData(dep.From()),
		Amount: dep.Amount(),
	}

	return &out
}

func toDataSet(ps entity.PartialSet) (*DataSet, error) {
	ins := ps.Instances()
	deposits := []*Data{}
	for _, oneIns := range ins {
		if dep, ok := oneIns.(Deposit); ok {
			deposits = append(deposits, toData(dep))
			continue
		}

		return nil, errors.New("there is at least one entity that is not a valid Deposit instance")
	}

	out := DataSet{
		Index:       ps.Index(),
		Amount:      ps.Amount(),
		TotalAmount: ps.TotalAmount(),
		IsLast:      ps.IsLast(),
		Deposits:    deposits,
	}

	return &out, nil
}
