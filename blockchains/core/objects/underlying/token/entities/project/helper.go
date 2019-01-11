package project

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/category"
	"github.com/xmnservices/xmnsuite/datastore"
)

func retrieveAllProjectKeyname() string {
	return "communityprojects"
}

func retrieveProjectByProposalKeyname(prop proposal.Proposal) string {
	base := retrieveAllProjectKeyname()
	return fmt.Sprintf("%s:by_proposal_id:%s", base, prop.ID().String())
}

func retrieveProjectByCategoryWalletKeyname(cat category.Category) string {
	base := retrieveAllProjectKeyname()
	return fmt.Sprintf("%s:by_category_id:%s", base, cat.ID().String())
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "CommunityProject",
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
					retrieveProjectByProposalKeyname(proj.Proposal()),
					retrieveProjectByCategoryWalletKeyname(proj.Proposal().Category()),
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
				proposalMetaData := proposal.SDKFunc.CreateMetaData()

				// create the repository and service:
				entityRepository := entity.SDKFunc.CreateRepository(ds)
				repository := createRepository(metaData, entityRepository)

				// make sure the proposal exists:
				_, retPropErr := entityRepository.RetrieveByID(proposalMetaData, proj.Proposal().ID())
				if retPropErr != nil {
					str := fmt.Sprintf("the given project (ID: %s) contains a proposal (ID: %s) that does not exists", proj.ID().String(), proj.Proposal().ID().String())
					return errors.New(str)
				}

				// make sure the proposal is not linked to any other project:
				retBindedProp, retBindedPropErr := repository.RetrieveByProposal(proj.Proposal())
				if retBindedPropErr == nil {
					str := fmt.Sprintf("the given project (ID: %s) contains a proposal (ID: %s) that is already binded to another project (ID: %s)", proj.ID().String(), proj.Proposal().ID().String(), retBindedProp.ID().String())
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
