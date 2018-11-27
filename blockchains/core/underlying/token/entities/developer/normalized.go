package developer

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/entities/user"
)

type normalizedDeveloper struct {
	ID     string          `json:"id"`
	User   user.Normalized `json:"user"`
	Name   string          `json:"name"`
	Resume string          `json:"resume"`
}

func createNormalizedDeveloper(dev Developer) (*normalizedDeveloper, error) {
	usr, usrErr := user.SDKFunc.CreateMetaData().Normalize()(dev.User())
	if usrErr != nil {
		return nil, usrErr
	}

	out := normalizedDeveloper{
		ID:     dev.ID().String(),
		User:   usr,
		Name:   dev.Name(),
		Resume: dev.Resume(),
	}

	return &out, nil

}
