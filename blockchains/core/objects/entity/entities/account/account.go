package account

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/work"
)

type account struct {
	Usr user.User `json:"user"`
	Wrk work.Work `json:"Work"`
}

func createAccount(usr user.User, wrk work.Work) (Account, error) {
	out := account{
		Usr: usr,
		Wrk: wrk,
	}

	return &out, nil
}

func fromNormalizedToAccount(ins Normalized) (Account, error) {
	if normalized, ok := ins.(*normalizedAccount); ok {
		usrIns, usrInsErr := user.SDKFunc.CreateMetaData().Denormalize()(normalized.User)
		if usrInsErr != nil {
			return nil, usrInsErr
		}

		if usr, ok := usrIns.(user.User); ok {
			wrk := work.SDKFunc.Denormalize(normalized.Work)
			out, outErr := createAccount(usr, wrk)
			if outErr != nil {
				return nil, outErr
			}

			return out, nil
		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid User instance", usrIns.ID().String())
		return nil, errors.New(str)
	}

	str := fmt.Sprintf("the given normalized instance is not a valid Account")
	return nil, errors.New(str)
}

// User returns the user
func (obj *account) User() user.User {
	return obj.Usr
}

// Work returns the work
func (obj *account) Work() work.Work {
	return obj.Wrk
}
