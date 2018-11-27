package developer

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/datastore"
)

// Developer represents a developer
type Developer interface {
	ID() *uuid.UUID
	User() user.User
	Name() string
	Resume() string
}

// Normalized represents a normalized developer
type Normalized interface {
}

// CreateParams represents the create params
type CreateParams struct {
	ID     *uuid.UUID
	User   user.User
	Name   string
	Resume string
}

// SDKFunc represents the Transfer SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Developer
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
}{
	Create: func(params CreateParams) Developer {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out := createDeveloper(params.ID, params.User, params.Name, params.Resume)
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
					}, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Developer instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Sync: func(ds datastore.DataStore, ins entity.Entity) error {
				// create the repository and service:
				repository := entity.SDKFunc.CreateRepository(ds)

				// create the metadata:
				metaData := createMetaData()
				userMetaData := user.SDKFunc.CreateMetaData()

				if dev, ok := ins.(Developer); ok {
					// if the developer already exists, return an error:
					_, retDevErr := repository.RetrieveByID(metaData, dev.ID())
					if retDevErr == nil {
						str := fmt.Sprintf("the Developer (ID: %s) already exists", dev.ID().String())
						return errors.New(str)
					}

					// if the user does not exists, return an error:
					usr := dev.User()
					_, retUserErr := repository.RetrieveByID(userMetaData, usr.ID())
					if retUserErr != nil {
						str := fmt.Sprintf("the User (ID: %s) in the Developer instance (ID: %s) does not exists", usr.ID().String(), dev.ID().String())
						return errors.New(str)
					}

					// if the user is already attached to a developer, return an error:
					retDevByUserID, retDevByUserIDErr := repository.RetrieveByIntersectKeynames(metaData, []string{retrieveDevelopersByUserIDKeyname(usr.ID())})
					if retDevByUserIDErr == nil {
						str := fmt.Sprintf("the User (ID: %s) already exists in another Developer instance (ID: %s)", usr.ID().String(), retDevByUserID.ID().String())
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
}
