package entity

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type retrievalMetaData struct {
	name        string
	keyname     string
	toEntity    ToEntity
	empStorable interface{}
}

func createMetaData(name string, toEntity ToEntity, empStorable interface{}) (MetaData, error) {

	if len(name) < 3 {
		str := fmt.Sprintf("the minimum length for the name is 3 characters: %d given", len(name))
		return nil, errors.New(str)
	}

	out := retrievalMetaData{
		name:        name,
		keyname:     strings.ToLower(name),
		toEntity:    toEntity,
		empStorable: empStorable,
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

// CopyStorable copies the empty storable instance and returns it
func (obj *retrievalMetaData) CopyStorable() interface{} {
	return reflect.New(reflect.ValueOf(obj.empStorable).Elem().Type()).Interface().(interface{})
}
