package request

import (
	"errors"
	"fmt"
	"log"

	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
)

type registry struct {
	metadatas map[string]entity.MetaData
}

func createRegistry() *registry {
	out := registry{
		metadatas: map[string]entity.MetaData{},
	}
	return &out
}

func (app *registry) register(metadata entity.MetaData) error {
	keyname := metadata.Keyname()
	if _, ok := app.metadatas[keyname]; ok {
		log.Printf("the given metadata (entity name: %s) is already registered ... skipping", keyname)
		return nil
	}

	app.metadatas[keyname] = metadata
	return nil
}

func (app *registry) fromJSONToEntity(js []byte, name string) (entity.Entity, error) {

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

	if met, ok := app.metadatas[name]; ok {
		// try the normalized:
		nIns, nInsErr := denormalize(js, met.CopyNormalized(), met.Denormalize())
		if nInsErr == nil {
			return nIns, nil
		}

		// try the storable:
		sIns, sInsErr := denormalize(js, met.CopyStorable(), met.Denormalize())
		if sInsErr == nil {
			return sIns, nil
		}

		str := fmt.Sprintf("the given name (%s) does not have a metadata that matches the given entity for a request", name)
		return nil, errors.New(str)
	}

	str := fmt.Sprintf("the given name (%s) is not a valid registered entity name for a request", name)
	return nil, errors.New(str)
}

func (app *registry) fromEntityToJSON(ins entity.Entity, name string) ([]byte, error) {
	if met, ok := app.metadatas[name]; ok {
		normalized, normalizedErr := met.Normalize()(ins)
		if normalizedErr != nil {
			return nil, normalizedErr
		}

		js, jsErr := cdc.MarshalJSON(normalized)
		if jsErr != nil {
			return nil, jsErr
		}

		return js, nil
	}

	str := fmt.Sprintf("the given name (%s) is not a valid registered entity name for a request", name)
	return nil, errors.New(str)
}
