package milestone

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project"
	"github.com/xmnservices/xmnsuite/datastore"
)

func retrieveAllMilestoneKeyname() string {
	return "milestones"
}

func retrieveMilestoneByWalletKeyname(wal wallet.Wallet) string {
	base := retrieveAllMilestoneKeyname()
	return fmt.Sprintf("%s:by_wallet_id:%s", base, wal.ID().String())
}

func retrieveMilestoneByProjectKeyname(proj project.Project) string {
	base := retrieveAllMilestoneKeyname()
	return fmt.Sprintf("%s:by_project_id:%s", base, proj.ID().String())
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Milestone",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			if storable, ok := data.(*storableMilestone); ok {
				return createMilestoneFromStorable(storable, rep)
			}

			if dataAsBytes, ok := data.([]byte); ok {
				ptr := new(normalizedMilestone)
				jsErr := cdc.UnmarshalJSON(dataAsBytes, ptr)
				if jsErr != nil {
					return nil, jsErr
				}

				return createMilestoneFromNormalized(ptr)
			}

			str := fmt.Sprintf("the given data does not represent a Milestone instance: %s", data)
			return nil, errors.New(str)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if mils, ok := ins.(Milestone); ok {
				out, outErr := createNormalizedMilestone(mils)
				if outErr != nil {
					return nil, outErr
				}

				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Milestone instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*normalizedMilestone); ok {
				return createMilestoneFromNormalized(normalized)
			}

			return nil, errors.New("the given normalized instance cannot be converted to a Milestone instance")
		},
		EmptyStorable:   new(storableMilestone),
		EmptyNormalized: new(normalizedMilestone),
	})
}

func representation() entity.Representation {
	return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
		Met: createMetaData(),
		ToStorable: func(ins entity.Entity) (interface{}, error) {
			if mils, ok := ins.(Milestone); ok {
				out := createStorableMilestone(mils)
				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Milestone instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Keynames: func(ins entity.Entity) ([]string, error) {
			if mils, ok := ins.(Milestone); ok {
				return []string{
					retrieveAllMilestoneKeyname(),
					retrieveMilestoneByWalletKeyname(mils.Wallet()),
					retrieveMilestoneByProjectKeyname(mils.Project()),
				}, nil
			}

			str := fmt.Sprintf("the entity (ID: %s) is not a valid Milestone instance", ins.ID().String())
			return nil, errors.New(str)
		},
		OnSave: func(ds datastore.DataStore, ins entity.Entity) error {
			if mils, ok := ins.(Milestone); ok {
				// crate metadata and representation:
				rep := representation()
				walletRepresentation := wallet.SDKFunc.CreateRepresentation()
				metaData := rep.MetaData()
				projectMetaData := project.SDKFunc.CreateMetaData()

				// create the repository and service:
				entityRepository := entity.SDKFunc.CreateRepository(ds)
				entityService := entity.SDKFunc.CreateService(ds)
				repository := createRepository(metaData, entityRepository)

				// make sure the project exists:
				_, retProjErr := entityRepository.RetrieveByID(projectMetaData, mils.Project().ID())
				if retProjErr != nil {
					str := fmt.Sprintf("the given milestone (ID: %s) contains a project (ID: %s) that does not exists", mils.ID().String(), mils.Project().ID().String())
					return errors.New(str)
				}

				// make sure the wallet does not exists:
				_, retBindedWalErr := repository.RetrieveByWallet(mils.Wallet())
				if retBindedWalErr == nil {
					str := fmt.Sprintf("the given milestone (ID: %s) contains a wallet (ID: %s) that already exists, only new wallets must be used in milestones", mils.ID().String(), mils.Wallet().ID().String())
					return errors.New(str)
				}

				// save the wallet:
				saveWalErr := entityService.Save(mils.Wallet(), walletRepresentation)
				if saveWalErr != nil {
					return saveWalErr
				}

				// everything is alright:
				return nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Milestone instance", ins.ID().String())
			return errors.New(str)
		},
	})
}
