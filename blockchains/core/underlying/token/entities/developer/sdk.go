package developer

import (
	"bytes"
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/entities/pledge"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/datastore"
)

// Developer represents a developer
type Developer interface {
	ID() *uuid.UUID
	User() user.User
	Pledge() pledge.Pledge
	Name() string
	Resume() string
}

// Normalized represents a normalized developer
type Normalized interface {
}

// Repository represents a developer repository
type Repository interface {
	RetrieveByUser(usr user.User) (Developer, error)
	RetrieveByPledge(pldge pledge.Pledge) (Developer, error)
}

// CreateParams represents the create params
type CreateParams struct {
	ID     *uuid.UUID
	Pledge pledge.Pledge
	User   user.User
	Name   string
	Resume string
}

// SDKFunc represents the Developer SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Developer
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	CreateRepository     func(store datastore.DataStore) Repository
}{
	Create: func(params CreateParams) Developer {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out := createDeveloper(params.ID, params.User, params.Pledge, params.Name, params.Resume)
		return out
	},
	CreateMetaData: func() entity.MetaData {
		out := createMetaData()
		return out
	},
	CreateRepresentation: func() entity.Representation {
		return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
			Met: createMetaData(),
			ToStorable: func(ins entity.Entity) (interface{}, error) {
				if dev, ok := ins.(Developer); ok {
					out := createStorableDeveloper(dev)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Developer instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				if dev, ok := ins.(Developer); ok {
					return []string{
						retrieveAllDevelopersKeyname(),
						retrieveDevelopersByUserIDKeyname(dev.User().ID()),
						retrieveDevelopersByPledgeIDKeyname(dev.Pledge().ID()),
					}, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Developer instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Sync: func(ds datastore.DataStore, ins entity.Entity) error {
				// create the metadata:
				metaData := createMetaData()
				pledgeMetaData := pledge.SDKFunc.CreateMetaData()
				userMetaData := user.SDKFunc.CreateMetaData()

				// create the repository and service:
				repository := entity.SDKFunc.CreateRepository(ds)
				devRepository := createRepository(repository, metaData)

				if dev, ok := ins.(Developer); ok {
					// if the developer already exists, return an error:
					_, retDevErr := repository.RetrieveByID(metaData, dev.ID())
					if retDevErr == nil {
						str := fmt.Sprintf("the Developer (ID: %s) already exists", dev.ID().String())
						return errors.New(str)
					}

					// if the user does not exists, return an error:
					usr := dev.User()
					_, retUsrErr := repository.RetrieveByID(userMetaData, usr.ID())
					if retUsrErr != nil {
						str := fmt.Sprintf("the User (ID: %s) in the Developer instance (ID: %s) does not exists", usr.ID().String(), dev.ID().String())
						return errors.New(str)
					}

					// if the pledge does not exists, return an error:
					pldge := dev.Pledge()
					_, retPledgeErr := repository.RetrieveByID(pledgeMetaData, pldge.ID())
					if retPledgeErr != nil {
						str := fmt.Sprintf("the Pledge (ID: %s) in the Developer instance (ID: %s) does not exists", pldge.ID().String(), dev.ID().String())
						return errors.New(str)
					}

					// if the pledge is already attached to a developer, return an error:
					retByPldgeID, retByPldgeIDErr := devRepository.RetrieveByPledge(pldge)
					if retByPldgeIDErr == nil {
						str := fmt.Sprintf("the Pledge (ID: %s) already exists in another Developer instance (ID: %s)", pldge.ID().String(), retByPldgeID.ID().String())
						return errors.New(str)
					}

					// if the user is already attached to a developer, return an error:
					retDevByUserID, retDevByUserIDErr := devRepository.RetrieveByUser(usr)
					if retDevByUserIDErr == nil {
						str := fmt.Sprintf("the User (ID: %s) already exists in another Developer instance (ID: %s)", usr.ID().String(), retDevByUserID.ID().String())
						return errors.New(str)
					}

					// if the wallet of the pledge is different from the one of the user, return an error:
					if bytes.Compare(pldge.From().From().ID().Bytes(), usr.Wallet().ID().Bytes()) != 0 {
						str := fmt.Sprintf("the Pledge (ID: %s) is from a Wallet (ID: %s) that is different from the User (ID: %s) Wallet (ID: %s), in the Developer instance (ID: %s)", pldge.ID().String(), pldge.From().From().ID().String(), usr.ID().String(), usr.Wallet().ID().String(), dev.ID().String())
						return errors.New(str)
					}

					// everything is alright:
					return nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Developer instance", ins.ID().String())
				return errors.New(str)
			},
		})
	},
	CreateRepository: func(store datastore.DataStore) Repository {
		met := createMetaData()
		entityRepository := entity.SDKFunc.CreateRepository(store)
		rep := createRepository(entityRepository, met)
		return rep
	},
}
