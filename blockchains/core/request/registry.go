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

// FromJSONToEntity converts JSON data to an entity instande
func (app *registry) FromJSONToEntity(js []byte) (entity.Entity, error) {
	for _, oneMetaData := range app.metadatas {
		ptr := oneMetaData.CopyNormalized()
		jsErr := cdc.UnmarshalJSON(js, ptr)
		if jsErr == nil {
			ins, insErr := oneMetaData.Denormalize()(ptr)
			if insErr != nil {
				return nil, insErr
			}

			return ins, nil
		}
	}

	return nil, errors.New("the given JSON data does not match any registered entity")
}
