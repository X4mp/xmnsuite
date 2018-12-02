package transfer

import (
	"bytes"
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/withdrawal"
	"github.com/xmnservices/xmnsuite/datastore"
)

// Transfer represents a transfer of token that can be claimed
type Transfer interface {
	ID() *uuid.UUID
	Withdrawal() withdrawal.Withdrawal
	Deposit() deposit.Deposit
}

// Normalized represents the normalized transfer
type Normalized interface {
}

// Repository represents the transfer reposiotry
type Repository interface {
	RetrieveByID(id *uuid.UUID) (Transfer, error)
}

// CreateParams represents the Create params
type CreateParams struct {
	ID         *uuid.UUID
	Withdrawal withdrawal.Withdrawal
	Deposit    deposit.Deposit
}

// SDKFunc represents the Transfer SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Transfer
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
}{
	Create: func(params CreateParams) Transfer {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out := createTransfer(params.ID, params.Withdrawal, params.Deposit)
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
				if trans, ok := ins.(Transfer); ok {
					out := createStorableTransfer(trans)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Transfer instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				return []string{
					retrieveAllTransfersKeyname(),
				}, nil
			},
			Sync: func(ds datastore.DataStore, ins entity.Entity) error {
				// create the repository and service:
				repository := entity.SDKFunc.CreateRepository(ds)
				service := entity.SDKFunc.CreateService(ds)

				// create the representations:
				withdrawalRepresentation := withdrawal.SDKFunc.CreateRepresentation()
				depositRepresentation := deposit.SDKFunc.CreateRepresentation()

				if trsf, ok := ins.(Transfer); ok {
					// variables:
					with := trsf.Withdrawal()
					dep := trsf.Deposit()

					// make sure the withdrawal wallet is not the same as the deposit wallet:
					if bytes.Compare(with.From().ID().Bytes(), dep.To().ID().Bytes()) == 0 {
						str := fmt.Sprintf("the wallet of the from withdrawal (ID: %s) cannot be the same as the deposit wallet (ID: %s)", with.From().ID().String(), dep.To().ID().String())
						return errors.New(str)
					}

					// make sure the token of the withdrawal matches the deposit token:
					if bytes.Compare(with.Token().ID().Bytes(), dep.Token().ID().Bytes()) != 0 {
						str := fmt.Sprintf("the withdrawal token (ID: %s) does not match the deposit token (ID: %s)", with.Token().ID().String(), dep.Token().ID().String())
						return errors.New(str)
					}

					// make sure the withdrawal amount matches the deposit amount:
					if with.Amount() != dep.Amount() {
						str := fmt.Sprintf("the withdrawal amount (%d) does not match the deposit amount (%d)", with.Amount(), dep.Amount())
						return errors.New(str)
					}

					// try to retrieve the withdrawal, send an error if it exists:
					_, retWithErr := repository.RetrieveByID(withdrawalRepresentation.MetaData(), with.ID())
					if retWithErr == nil {
						str := fmt.Sprintf("the Transfer instance (ID: %s) contains a Withdrawal instance that already exists", with.ID().String())
						return errors.New(str)
					}

					// save the withdrawal:
					saveErr := service.Save(with, withdrawalRepresentation)
					if saveErr != nil {
						return saveErr
					}

					// try to retrieve the deposit, send an error if it exists:
					_, retDepErr := repository.RetrieveByID(withdrawalRepresentation.MetaData(), dep.ID())
					if retDepErr == nil {
						str := fmt.Sprintf("the Transfer instance (ID: %s) contains a Withdrawal instance that already exists", dep.ID().String())
						return errors.New(str)
					}

					// save the deposit:
					saveDepErr := service.Save(dep, depositRepresentation)
					if saveDepErr != nil {
						return saveDepErr
					}

					return nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Transfer instance", ins.ID().String())
				return errors.New(str)
			},
		})
	},
}
