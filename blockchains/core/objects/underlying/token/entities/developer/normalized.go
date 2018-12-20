package developer

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/pledge"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/user"
)

type normalizedDeveloper struct {
	ID     string            `json:"id"`
	Pledge pledge.Normalized `json:"pledge"`
	User   user.Normalized   `json:"user"`
	Name   string            `json:"name"`
	Resume string            `json:"resume"`
}

func createNormalizedDeveloper(dev Developer) (*normalizedDeveloper, error) {
	usr, usrErr := user.SDKFunc.CreateMetaData().Normalize()(dev.User())
	if usrErr != nil {
		return nil, usrErr
	}

	pldge, pldgeErr := pledge.SDKFunc.CreateMetaData().Normalize()(dev.Pledge())
	if pldgeErr != nil {
		return nil, pldgeErr
	}

	out := normalizedDeveloper{
		ID:     dev.ID().String(),
		Pledge: pldge,
		User:   usr,
		Name:   dev.Name(),
		Resume: dev.Resume(),
	}

	return &out, nil

}
