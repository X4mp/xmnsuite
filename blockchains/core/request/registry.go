package request

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
)

type registry struct {
	metadatas map[string]entity.MetaData
}

func createRegistry() Registry {
	out := registry{
		metadatas: map[string]entity.MetaData{},
	}
	return &out
}

// Register registers an entity
func (app *registry) Register(metadata entity.MetaData) error {
	keyname := metadata.Keyname()
	if _, ok := app.metadatas[keyname]; ok {
		str := fmt.Sprintf("the given metadata (entity name: %s) is already registered", metadata.Name())
		return errors.New(str)
	}

	app.metadatas[keyname] = metadata
	return nil
}

// FromJSONToEntity converts JSON data to an entity instance
func (app *registry) FromJSONToEntity(js []byte) (entity.Entity, error) {

	denormalize := func(js []byte, ptr interface{}, fn entity.Denormalize) (entity.Entity, error) {
		jsErr := cdc.UnmarshalJSON(js, ptr)
		if jsErr != nil {
			return nil, jsErr
		}

		ins, insErr := fn(ptr)
		if insErr != nil {
			return nil, insErr
		}

		return ins, nil
	}

	for _, oneMetaData := range app.metadatas {
		// try the normalized:
		nIns, nInsErr := denormalize(js, oneMetaData.CopyNormalized(), oneMetaData.Denormalize())
		if nInsErr == nil {
			return nIns, nil
		}

		// try the storable:
		sIns, sInsErr := denormalize(js, oneMetaData.CopyStorable(), oneMetaData.Denormalize())
		if sInsErr == nil {
			return sIns, nil
		}
	}

	return nil, errors.New("the given JSON data does not match any registered entity")
}

// FromEntityToJSON converts an entity instance to JSON
func (app *registry) FromEntityToJSON(ins entity.Entity) ([]byte, error) {
	for _, oneMetaData := range app.metadatas {
		normalized, normalizedErr := oneMetaData.Normalize()(ins)
		if normalizedErr == nil {
			js, jsErr := cdc.MarshalJSON(normalized)
			if jsErr != nil {
				return nil, jsErr
			}

			return js, nil
		}
	}

	return nil, errors.New("the given entity data does not match any registered entity")
}
