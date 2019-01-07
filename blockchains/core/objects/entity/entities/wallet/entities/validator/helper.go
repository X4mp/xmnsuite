package validator

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"net"

	uuid "github.com/satori/go.uuid"
	"github.com/tendermint/tendermint/crypto"
	ed25519 "github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/pledge"
)

func retrieveAllValidatorsKeyname() string {
	return "validators"
}

func retrieveValidatorsByPledgeKeyname(pldge pledge.Pledge) string {
	base := retrieveAllValidatorsKeyname()
	return fmt.Sprintf("%s:by_pledge_id:%s", base, pldge.ID().String())
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Validator",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			fromStorableToEntity := func(storable *storableValidator) (entity.Entity, error) {
				id, idErr := uuid.FromString(storable.ID)
				if idErr != nil {
					return nil, idErr
				}

				pledgeID, pledgeIDErr := uuid.FromString(storable.PledgeID)
				if pledgeIDErr != nil {
					return nil, pledgeIDErr
				}

				pubkey, pubKeyErr := fromEncodedStringToPubKey(storable.PubKey)
				if pubKeyErr != nil {
					str := fmt.Sprintf("the storable pubKey (%s) is invalid: %s", storable.PubKey, pubKeyErr.Error())
					return nil, errors.New(str)
				}

				// retrieve the pledge:
				pledgeMetaData := pledge.SDKFunc.CreateMetaData()
				pledgeIns, pledgeInsErr := rep.RetrieveByID(pledgeMetaData, &pledgeID)
				if pledgeInsErr != nil {
					return nil, pledgeInsErr
				}

				if pldge, ok := pledgeIns.(pledge.Pledge); ok {
					ip := net.ParseIP(storable.IP)
					out := createValidator(&id, ip, storable.Port, pubkey, pldge)
					return out, nil
				}

				str := fmt.Sprintf("the entity (ID: %s) is not a valid Pledge instance", pledgeID.String())
				return nil, errors.New(str)
			}

			if storable, ok := data.(*storableValidator); ok {
				return fromStorableToEntity(storable)
			}

			ptr := new(normalizedValidator)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return createValidatorFromNormalized(ptr)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if val, ok := ins.(Validator); ok {
				return createNormalizedValidator(val)
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Validator instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*normalizedValidator); ok {
				return createValidatorFromNormalized(normalized)
			}

			return nil, errors.New("the given instance is not a valid normalized Validator instance")
		},
		EmptyStorable:   new(storableValidator),
		EmptyNormalized: new(normalizedValidator),
	})
}

func fromEncodedStringToPubKey(str string) (crypto.PubKey, error) {
	pubKeyAsBytes, pubKeyAsBytesErr := hex.DecodeString(str)
	if pubKeyAsBytesErr != nil {
		return nil, pubKeyAsBytesErr
	}

	pubKey := new(ed25519.PubKeyEd25519)
	pubKeyErr := cdc.UnmarshalBinaryBare(pubKeyAsBytes, pubKey)
	if pubKeyErr != nil {
		str := fmt.Sprintf("the public key []byte is invalid: %s", pubKeyErr.Error())
		return nil, errors.New(str)
	}

	return pubKey, nil
}

func orderValPSByPledge(valIns []Validator, index int, amount int) ([]Validator, error) {

	getSmallestValidator := func(combinedVals map[int][]Validator) (int, []Validator) {
		smallest := math.MaxInt64 - 1
		for amount := range combinedVals {
			if smallest >= amount {
				smallest = amount
			}
		}

		return smallest, combinedVals[smallest]
	}

	combinedVals := map[int][]Validator{}
	for _, val := range valIns {
		isIn := false
		pldgeAmount := val.Pledge().From().Amount()
		for am := range combinedVals {
			if am == pldgeAmount {
				combinedVals[am] = append(combinedVals[am], val)
				isIn = true
				break
			}
		}

		if !isIn {
			combinedVals[pldgeAmount] = []Validator{
				val,
			}
		}
	}

	length := len(combinedVals)
	orderedCombinedVals := [][]Validator{}
	for i := 0; i < length; i++ {
		idx, vals := getSmallestValidator(combinedVals)
		orderedCombinedVals = append(orderedCombinedVals, vals)
		delete(combinedVals, idx)
	}

	if index < 0 {
		index = 0
	}

	if amount < 0 {
		amount = 0
	}

	valLength := len(valIns)
	if index >= valLength {
		index = valLength
	}

	reqLength := index + amount
	if reqLength >= valLength {
		amount = valLength - index
	}

	out := []Validator{}
	for _, vals := range orderedCombinedVals {
		out = append(out, vals...)
	}

	return out[index : index+amount], nil
}
