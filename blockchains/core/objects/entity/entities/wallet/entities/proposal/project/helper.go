package project

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/pledge"
	approved_project "github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/project"
	"github.com/xmnservices/xmnsuite/datastore"
)

func retrieveAllProjectKeyname() string {
	return "projects"
}

func retrieveProjectByProjectKeyname(proj approved_project.Project) string {
	base := retrieveAllProjectKeyname()
	return fmt.Sprintf("%s:by_project_id:%s", base, proj.ID().String())
}

func retrieveProjectByOwnerWalletKeyname(owner wallet.Wallet) string {
	base := retrieveAllProjectKeyname()
	return fmt.Sprintf("%s:by_owner_wallet_id:%s", base, owner.ID().String())
}

func retrieveProjectByManagerWalletKeyname(mgr wallet.Wallet) string {
	base := retrieveAllProjectKeyname()
	return fmt.Sprintf("%s:by_manager_wallet_id:%s", base, mgr.ID().String())
}

func retrieveProjectByLinkerWalletKeyname(linker wallet.Wallet) string {
	base := retrieveAllProjectKeyname()
	return fmt.Sprintf("%s:by_linker_wallet_id:%s", base, linker.ID().String())
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Project",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			if storable, ok := data.(*storableProject); ok {
				return createProjectFromStorable(storable, rep)
			}

			if dataAsBytes, ok := data.([]byte); ok {
				ptr := new(normalizedProject)
				jsErr := cdc.UnmarshalJSON(dataAsBytes, ptr)
				if jsErr != nil {
					return nil, jsErr
				}

				return createProjectFromNormalized(ptr)
			}

			str := fmt.Sprintf("the given data does not represent a Project instance: %s", data)
			return nil, errors.New(str)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if proj, ok := ins.(Project); ok {
				out, outErr := createNormalizedProject(proj)
				if outErr != nil {
					return nil, outErr
				}

				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Project instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*normalizedProject); ok {
				return createProjectFromNormalized(normalized)
			}

			return nil, errors.New("the given normalized instance cannot be converted to a Project instance")
		},
		EmptyStorable:   new(storableProject),
		EmptyNormalized: new(normalizedProject),
	})
}

func representation() entity.Representation {
	return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
		Met: createMetaData(),
		ToStorable: func(ins entity.Entity) (interface{}, error) {
			if proj, ok := ins.(Project); ok {
				out := createStorableProject(proj)
				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Project instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Keynames: func(ins entity.Entity) ([]string, error) {
			if proj, ok := ins.(Project); ok {
				return []string{
					retrieveAllProjectKeyname(),
					retrieveProjectByProjectKeyname(proj.Project()),
					retrieveProjectByOwnerWalletKeyname(proj.Owner()),
					retrieveProjectByManagerWalletKeyname(proj.Manager()),
					retrieveProjectByLinkerWalletKeyname(proj.Linker()),
				}, nil
			}

			str := fmt.Sprintf("the entity (ID: %s) is not a valid Project instance", ins.ID().String())
			return nil, errors.New(str)
		},
		OnSave: func(ds datastore.DataStore, ins entity.Entity) error {
			if proj, ok := ins.(Project); ok {
				// crate metadata and representation:
				rep := representation()
				metaData := rep.MetaData()
				walletMetaData := wallet.SDKFunc.CreateMetaData()
				projectMetaData := approved_project.SDKFunc.CreateMetaData()

				// create the repository and service:
				entityRepository := entity.SDKFunc.CreateRepository(ds)
				repository := createRepository(metaData, entityRepository)
				pledgeRepository := pledge.SDKFunc.CreateRepository(pledge.CreateRepositoryParams{
					EntityRepository: entityRepository,
				})

				// make sure the project exists:
				_, retProjErr := entityRepository.RetrieveByID(projectMetaData, proj.Project().ID())
				if retProjErr != nil {
					str := fmt.Sprintf("the given project (ID: %s) contains a project (ID: %s) that does not exists", proj.ID().String(), proj.Project().ID().String())
					return errors.New(str)
				}

				// make sure the wallet owner exists:
				_, retOwnerWalErr := entityRepository.RetrieveByID(walletMetaData, proj.Owner().ID())
				if retOwnerWalErr != nil {
					str := fmt.Sprintf("the given project (ID: %s) contains a wallet owner (ID: %s) that does not exists", proj.ID().String(), proj.Owner().ID().String())
					return errors.New(str)
				}

				// make sure the wallet manager exists:
				_, retMgrWalErr := entityRepository.RetrieveByID(walletMetaData, proj.Manager().ID())
				if retMgrWalErr != nil {
					str := fmt.Sprintf("the given project (ID: %s) contains a wallet manager (ID: %s) that does not exists", proj.ID().String(), proj.Manager().ID().String())
					return errors.New(str)
				}

				// make sure the wallet linker exists:
				_, retLinkerWalErr := entityRepository.RetrieveByID(walletMetaData, proj.Linker().ID())
				if retLinkerWalErr != nil {
					str := fmt.Sprintf("the given project (ID: %s) contains a wallet linker (ID: %s) that does not exists", proj.ID().String(), proj.Linker().ID().String())
					return errors.New(str)
				}

				// make sure the project is not linked to any other project:
				retBindedProj, retBindedProjErr := repository.RetrieveByProject(proj.Project())
				if retBindedProjErr == nil {
					str := fmt.Sprintf("the given project (ID: %s) contains a project (ID: %s) that is already binded to another project (ID: %s)", proj.ID().String(), proj.Project().ID().String(), retBindedProj.ID().String())
					return errors.New(str)
				}

				// make sure the wallet owner is not linked to any other project:
				retBindedOwner, retBindedOwnerErr := repository.RetrieveByOwner(proj.Owner())
				if retBindedOwnerErr == nil {
					str := fmt.Sprintf("the given project (ID: %s) contains a wallet owner (ID: %s) that is already binded to another project (ID: %s)", proj.ID().String(), proj.Owner().ID().String(), retBindedOwner.ID().String())
					return errors.New(str)
				}

				// make sure the wallet manager is not linked to any other project:
				retBindedMgr, retBindedMgrErr := repository.RetrieveByManager(proj.Manager())
				if retBindedMgrErr == nil {
					str := fmt.Sprintf("the given project (ID: %s) contains a wallet manager (ID: %s) that is already binded to another project (ID: %s)", proj.ID().String(), proj.Manager().ID().String(), retBindedMgr.ID().String())
					return errors.New(str)
				}

				// make sure the wallet linker is not linked to any other project:
				retBindedLinker, retBindedLinkerErr := repository.RetrieveByLinker(proj.Linker())
				if retBindedLinkerErr == nil {
					str := fmt.Sprintf("the given project (ID: %s) contains a wallet linker (ID: %s) that is already binded to another project (ID: %s)", proj.ID().String(), proj.Linker().ID().String(), retBindedLinker.ID().String())
					return errors.New(str)
				}

				// make sure the manager made his pledge:
				managerPledge, managerPledgeErr := pledgeRepository.RetrieveByFromAndToWallet(proj.Manager(), proj.Owner())
				if managerPledgeErr != nil {
					str := fmt.Sprintf("there was an error while retrieving the manager pledge: %s", managerPledgeErr.Error())
					return errors.New(str)
				}

				// make sure the manager pledge is of the right amount:
				proposal := proj.Project().Proposal()
				if managerPledge.From().Amount() != proposal.ManagerPledgeNeeded() {
					str := fmt.Sprintf("the manager (ID: %s) pledge %d tokens, but the proposal (ID: %s) required %d tokens", managerPledge.ID().String(), managerPledge.From().Amount(), proposal.ID().String(), proposal.ManagerPledgeNeeded())
					return errors.New(str)
				}

				// make sure the linker made his pledge:
				linkerPledge, linkerPledgeErr := pledgeRepository.RetrieveByFromAndToWallet(proj.Linker(), proj.Owner())
				if linkerPledgeErr != nil {
					str := fmt.Sprintf("there was an error while retrieving the linker pledge: %s", linkerPledgeErr.Error())
					return errors.New(str)
				}

				// make sure the linker pledge is of the right amount:
				if linkerPledge.From().Amount() != proposal.LinkerPledgeNeeded() {
					str := fmt.Sprintf("the linker (ID: %s) pledge %d tokens, but the proposal (ID: %s) required %d tokens", linkerPledge.ID().String(), linkerPledge.From().Amount(), proposal.ID().String(), proposal.ManagerPledgeNeeded())
					return errors.New(str)
				}

				// everything is alright:
				return nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Project instance", ins.ID().String())
			return errors.New(str)
		},
	})
}
