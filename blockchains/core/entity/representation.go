package entity

import "errors"

type representation struct {
	met      MetaData
	keynames Keynames
	toData   ToStorable
	sync     Sync
}

func createRepresentation(met MetaData, toData ToStorable, keynames Keynames, sync Sync) (Representation, error) {

	if met == nil {
		return nil, errors.New("the metadata is mandatory in order to create a representation instance")
	}

	if toData == nil {
		return nil, errors.New("the toData is mandatory in order to create a representation instance")
	}

	out := representation{
		met:      met,
		keynames: keynames,
		toData:   toData,
		sync:     sync,
	}

	return &out, nil
}

// MetaData returns the MetaData instance
func (obj *representation) MetaData() MetaData {
	return obj.met
}

// ToStorable returns the instance to data
func (obj *representation) ToStorable() ToStorable {
	return obj.toData
}

// HasKeynames returns true if there is keynames, false otherwise
func (obj *representation) HasKeynames() bool {
	return obj.keynames != nil
}

// Keynames returns the keynames if any
func (obj *representation) Keynames() Keynames {
	return obj.keynames
}

// HasSync returns true if there is a Sync func, false otherwise
func (obj *representation) HasSync() bool {
	return obj.sync != nil
}

// Sync returns the sync func
func (obj *representation) Sync() Sync {
	return obj.sync
}
