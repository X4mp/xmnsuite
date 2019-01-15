package transfer

import (
	"bytes"
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/withdrawal"
	"github.com/xmnservices/xmnsuite/datastore"
)

func retrieveAllTransfersKeyname() string {
	return "transfers"
}

func retrieveTransfersByDepositKeyname(dep deposit.Deposit) string {
	base := retrieveAllTransfersKeyname()
	return fmt.Sprintf("%s:by_deposit_id:%s", base, dep.ID().String())
}

func retrieveTransfersByWithdrawalKeyname(with withdrawal.Withdrawal) string {
	base := retrieveAllTransfersKeyname()
	return fmt.Sprintf("%s:by_withdrawal_id:%s", base, with.ID().String())
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Transfer",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			fromStorableToEntity := func(storable *storableTransfer) (entity.Entity, error) {
				id, idErr := uuid.FromString(storable.ID)
				if idErr != nil {
					return nil, idErr
				}

				withdrawalID, withdrawalIDErr := uuid.FromString(storable.WithdrawalID)
				if withdrawalIDErr != nil {
					return nil, withdrawalIDErr
				}

				depositID, depositIDErr := uuid.FromString(storable.DepositID)
				if depositIDErr != nil {
					return nil, depositIDErr
				}

				// retrieve the withdrawal:
				withdrawalMetaData := withdrawal.SDKFunc.CreateMetaData()
				withdrawalIns, withdrawalInsErr := rep.RetrieveByID(withdrawalMetaData, &withdrawalID)
				if withdrawalInsErr != nil {
					return nil, withdrawalInsErr
				}

				// retrieve the deposit:
				depositMetaData := deposit.SDKFunc.CreateMetaData()
				depositIns, depositInsErr := rep.RetrieveByID(depositMetaData, &depositID)
				if depositInsErr != nil {
					return nil, depositInsErr
				}

				if withdrawl, ok := withdrawalIns.(withdrawal.Withdrawal); ok {
					if dep, ok := depositIns.(deposit.Deposit); ok {
						out := createTransfer(&id, withdrawl, dep)
						return out, nil
					}

					str := fmt.Sprintf("the entity (ID: %s) is not a valid Deposit instance", depositID.String())
					return nil, errors.New(str)
				}

				str := fmt.Sprintf("the entity (ID: %s) is not a valid Withdrawal instance", withdrawalID.String())
				return nil, errors.New(str)
			}

			if storable, ok := data.(*storableTransfer); ok {
				return fromStorableToEntity(storable)
			}

			ptr := new(normalizedTransfer)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return createTransferFromNormalized(ptr)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if trsf, ok := ins.(Transfer); ok {
				return createNormalizedTransfer(trsf)
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Transfer instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*normalizedTransfer); ok {
				return createTransferFromNormalized(normalized)
			}

			return nil, errors.New("the given instance is not a valid normalized Transfer instance")
		},
		EmptyStorable:   new(storableTransfer),
		EmptyNormalized: new(normalizedTransfer),
	})
}

func createRepresentation() entity.Representation {
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
			if trsf, ok := ins.(Transfer); ok {
				return []string{
					retrieveAllTransfersKeyname(),
					retrieveTransfersByDepositKeyname(trsf.Deposit()),
					retrieveTransfersByWithdrawalKeyname(trsf.Withdrawal()),
				}, nil
			}

			str := fmt.Sprintf("the entity (ID: %s) is not a valid Transfer instance", ins.ID().String())
			return nil, errors.New(str)
		},
		OnSave: func(ds datastore.DataStore, ins entity.Entity) error {
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
}
