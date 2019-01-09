package affiliates

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/datastore"
)

func retrieveAllAffiliatesKeyname() string {
	return "affiliates"
}

func retrieveAffiliatesByWalletKeyname(wal wallet.Wallet) string {
	base := retrieveAllAffiliatesKeyname()
	return fmt.Sprintf("%s:by_wallet_id:%s", base, wal.ID().String())
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Affiliate",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			if storable, ok := data.(*storableAffiliate); ok {
				return createAffiliateFromStorable(storable, rep)
			}

			ptr := new(normalizedAffiliate)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return createAffiliateFromNormalized(ptr)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if aff, ok := ins.(Affiliate); ok {
				out, outErr := createNormalizedAffiliate(aff)
				if outErr != nil {
					return nil, outErr
				}

				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Affiliate instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*normalizedAffiliate); ok {
				return createAffiliateFromNormalized(normalized)
			}

			return nil, errors.New("the given instance is not a valid normalized Affiliate instance")
		},
		EmptyNormalized: new(normalizedAffiliate),
		EmptyStorable:   new(storableAffiliate),
	})
}

func createRepresentation() entity.Representation {
	return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
		Met: createMetaData(),
		ToStorable: func(ins entity.Entity) (interface{}, error) {
			if aff, ok := ins.(Affiliate); ok {
				out := createStorableAffiliate(aff)
				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Affiliate instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Keynames: func(ins entity.Entity) ([]string, error) {
			if aff, ok := ins.(Affiliate); ok {
				return []string{
					retrieveAllAffiliatesKeyname(),
					retrieveAffiliatesByWalletKeyname(aff.Owner()),
				}, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Affiliate instance", ins.ID().String())
			return nil, errors.New(str)
		},
		OnSave: func(ds datastore.DataStore, ins entity.Entity) error {
			// create the entity repository and service:
			entityRepository := entity.SDKFunc.CreateRepository(ds)
			walletReposiotry := wallet.SDKFunc.CreateRepository(wallet.CreateRepositoryParams{
				EntityRepository: entityRepository,
			})

			if aff, ok := ins.(Affiliate); ok {
				// make sure the wallet exists:
				_, walErr := walletReposiotry.RetrieveByID(aff.Owner().ID())
				if walErr != nil {
					str := fmt.Sprintf("the Wallet (ID: %s) in the Affiliate instance does not exists", aff.Owner().ID().String())
					return errors.New(str)
				}

				return nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Affiliate instance", ins.ID().String())
			return errors.New(str)
		},
	})
}
