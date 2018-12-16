package entity

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type retrievalMetaData struct {
	name          string
	keyname       string
	toEntity      ToEntity
	normalize     Normalize
	denormalize   Denormalize
	empStorable   interface{}
	empNormalized interface{}
}

func createMetaData(name string, toEntity ToEntity, normalize Normalize, denormalize Denormalize, empStorable interface{}, empNormalized interface{}) (MetaData, error) {

	if len(name) < 3 {
		str := fmt.Sprintf("the minimum length for the name is 3 characters: %d given", len(name))
		return nil, errors.New(str)
	}

	if toEntity == nil {
		str := fmt.Sprintf("the toEntity param is mandatory in order to create a MetaData instance (name: %s)", name)
		return nil, errors.New(str)
	}

	if normalize == nil {
		str := fmt.Sprintf("the normalize param is mandatory in order to create a MetaData instance (name: %s)", name)
		return nil, errors.New(str)
	}

	if denormalize == nil {
		str := fmt.Sprintf("the denormalize param is mandatory in order to create a MetaData instance (name: %s)", name)
		return nil, errors.New(str)
	}

	if empStorable == nil {
		str := fmt.Sprintf("the emptyStorable param is mandatory in order to create a MetaData instance (name: %s)", name)
		return nil, errors.New(str)
	}

	if empNormalized == nil {
		str := fmt.Sprintf("the emptyNormalized param is mandatory in order to create a MetaData instance (name: %s)", name)
		return nil, errors.New(str)
	}

	out := retrievalMetaData{
		name:          name,
		keyname:       strings.ToLower(name),
		toEntity:      toEntity,
		normalize:     normalize,
		denormalize:   denormalize,
		empStorable:   empStorable,
		empNormalized: empNormalized,
	}

	return &out, nil
}

// Name returns the name
func (obj *retrievalMetaData) Name() string {
	return obj.name
}

// Keyname returns the keyname
func (obj *retrievalMetaData) Keyname() string {
	return obj.keyname
}

// ToEntity returns the ToEntity func
func (obj *retrievalMetaData) ToEntity() ToEntity {
	return obj.toEntity
}

// Normalize returns the Normalization func
func (obj *retrievalMetaData) Normalize() Normalize {
	return obj.normalize
}

// Denormalize returns the Denormalization func
func (obj *retrievalMetaData) Denormalize() Denormalize {
	return obj.denormalize
}

// CopyStorable copies the empty storable instance and returns it
func (obj *retrievalMetaData) CopyStorable() interface{} {
	return reflect.New(reflect.ValueOf(obj.empStorable).Elem().Type()).Interface().(interface{})
}

// CopyNormalized copies the empty normalized instance and returns it
func (obj *retrievalMetaData) CopyNormalized() interface{} {
	return reflect.New(reflect.ValueOf(obj.empNormalized).Elem().Type()).Interface().(interface{})
}
