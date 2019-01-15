package genesis

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/information"
)

/*
 * Genesis
 */

type genesis struct {
	UUID *uuid.UUID              `json:"id"`
	Inf  information.Information `json:"information"`
	Usr  user.User               `json:"user"`
	Dep  deposit.Deposit         `json:"deposit"`
	Tok  token.Token             `json:"token"`
}

func createGenesis(
	id *uuid.UUID,
	inf information.Information,
	dep deposit.Deposit,
	usr user.User,
	tok token.Token,
) (Genesis, error) {

	out := genesis{
		UUID: id,
		Inf:  inf,
		Usr:  usr,
		Dep:  dep,
		Tok:  tok,
	}

	return &out, nil
}

func createGenesisFromNormalized(ins *normalizedGenesis) (Genesis, error) {
	id, idErr := uuid.FromString(ins.ID)
	if idErr != nil {
		return nil, idErr
	}

	infIns, infInsErr := information.SDKFunc.CreateMetaData().Denormalize()(ins.Info)
	if infInsErr != nil {
		return nil, infInsErr
	}

	depIns, depInsErr := deposit.SDKFunc.CreateMetaData().Denormalize()(ins.Deposit)
	if depInsErr != nil {
		return nil, depInsErr
	}

	usrIns, usrInsErr := user.SDKFunc.CreateMetaData().Denormalize()(ins.User)
	if usrInsErr != nil {
		return nil, usrInsErr
	}

	tokIns, tokInsErr := token.SDKFunc.CreateMetaData().Denormalize()(ins.Token)
	if tokInsErr != nil {
		return nil, tokInsErr
	}

	if dep, ok := depIns.(deposit.Deposit); ok {
		if usr, ok := usrIns.(user.User); ok {
			if inf, ok := infIns.(information.Information); ok {
				if tok, ok := tokIns.(token.Token); ok {
					return createGenesis(&id, inf, dep, usr, tok)
				}

				str := fmt.Sprintf("the entity (ID: %s) is not a valid Token instance", tokIns.ID().String())
				return nil, errors.New(str)
			}

			str := fmt.Sprintf("the entity (ID: %s) is not a valid Information instance", infIns.ID().String())
			return nil, errors.New(str)
		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid User instance", usrIns.ID().String())
		return nil, errors.New(str)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Deposit instance", depIns.ID().String())
	return nil, errors.New(str)
}

// ID returns the ID
func (app *genesis) ID() *uuid.UUID {
	return app.UUID
}

// Info returns the information
func (app *genesis) Info() information.Information {
	return app.Inf
}

// User returns the user
func (app *genesis) User() user.User {
	return app.Usr
}

// Deposit returns the initial deposit
func (app *genesis) Deposit() deposit.Deposit {
	return app.Dep
}

// Token returns the token
func (app *genesis) Token() token.Token {
	return app.Tok
}
