package proposal

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/category"
	"github.com/xmnservices/xmnsuite/datastore"
)

func retrieveAllProposalKeyname() string {
	return "proposals"
}

func retrieveProposalByCategoryKeyname(cat category.Category) string {
	base := retrieveAllProposalKeyname()
	return fmt.Sprintf("%s:by_category_id:%s", base, cat.ID().String())
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Proposal",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			if storable, ok := data.(*storableProposal); ok {
				return createProposalFromStorable(storable, rep)
			}

			if dataAsBytes, ok := data.([]byte); ok {
				ptr := new(normalizedProposal)
				jsErr := cdc.UnmarshalJSON(dataAsBytes, ptr)
				if jsErr != nil {
					return nil, jsErr
				}

				return createProposalFromNormalized(ptr)
			}

			str := fmt.Sprintf("the given data does not represent a Proposal instance: %s", data)
			return nil, errors.New(str)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if prop, ok := ins.(Proposal); ok {
				out, outErr := createNormalizedProposal(prop)
				if outErr != nil {
					return nil, outErr
				}

				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Proposal instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*normalizedProposal); ok {
				return createProposalFromNormalized(normalized)
			}

			return nil, errors.New("the given normalized instance cannot be converted to a Proposal instance")
		},
		EmptyStorable:   new(storableProposal),
		EmptyNormalized: new(normalizedProposal),
	})
}

func representation() entity.Representation {
	return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
		Met: createMetaData(),
		ToStorable: func(ins entity.Entity) (interface{}, error) {
			if prop, ok := ins.(Proposal); ok {
				out := createStorableProposal(prop)
				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Proposal instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Keynames: func(ins entity.Entity) ([]string, error) {
			if prop, ok := ins.(Proposal); ok {
				return []string{
					retrieveAllProposalKeyname(),
					retrieveProposalByCategoryKeyname(prop.Category()),
				}, nil
			}

			str := fmt.Sprintf("the entity (ID: %s) is not a valid Proposal instance", ins.ID().String())
			return nil, errors.New(str)
		},
		OnSave: func(ds datastore.DataStore, ins entity.Entity) error {
			if prop, ok := ins.(Proposal); ok {
				// crate metadata and representation:
				categoryMetaData := category.SDKFunc.CreateMetaData()

				// create the repository and service:
				entityRepository := entity.SDKFunc.CreateRepository(ds)

				// make sure the category exists:
				_, retCatErr := entityRepository.RetrieveByID(categoryMetaData, prop.Category().ID())
				if retCatErr != nil {
					str := fmt.Sprintf("the given proposal (ID: %s) contains a category (ID: %s) that does not exists", prop.ID().String(), prop.Category().ID().String())
					return errors.New(str)
				}

				// everything is alright:
				return nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Proposal instance", ins.ID().String())
			return errors.New(str)
		},
	})
}
