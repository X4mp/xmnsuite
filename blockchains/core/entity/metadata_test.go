package entity

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	uuid "github.com/satori/go.uuid"
)

func TestMetaData_Success(t *testing.T) {
	// variables:
	name := "MyEntity"
	empStorable := new(testEntity)

	toEntity := func(rep Repository, data interface{}) (Entity, error) {
		if casted, ok := data.(*storableTestEntity); ok {
			id, idErr := uuid.FromString(casted.ID)
			if idErr != nil {
				return nil, idErr
			}

			return createTestEntity(&id, casted.Name), nil
		}

		return nil, errors.New("invalid")
	}

	normalize := func(ins Entity) (interface{}, error) {
		obj := ins.(*testEntity)
		return &storableTestEntity{
			ID:   obj.ID().String(),
			Name: obj.Name(),
		}, nil
	}

	denormalize := func(ins interface{}) (Entity, error) {
		if casted, ok := ins.(*storableTestEntity); ok {
			id, idErr := uuid.FromString(casted.ID)
			if idErr != nil {
				return nil, idErr
			}

			return createTestEntity(&id, casted.Name), nil
		}

		return nil, errors.New("invalid")
	}

	// execute, name too small, returns error:
	_, validMetErr := createMetaData("s", toEntity, normalize, denormalize, empStorable, empStorable)
	if validMetErr == nil {
		t.Errorf("the returned error was expected to be valid, nil returned")
		return
	}

	// execute:
	met, metErr := createMetaData(name, toEntity, normalize, denormalize, empStorable, empStorable)
	if metErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", metErr.Error())
		return
	}

	// compare:
	if met.Name() != name {
		t.Errorf("the returned name is invalid.  Expected: %s, returned: %s", name, met.Name())
		return
	}

	// copy the storable:
	cpy := met.CopyStorable()
	if cpy == empStorable {
		t.Errorf("the empty storable was expected to be copied, same instance returned")
		return
	}

	if cpy == nil {
		t.Errorf("the empty storable was not expected to be nil")
		return
	}

	// convert an instance baxck and forth:
	ins := createTestEntityForTests()

	// data:
	data := &storableTestEntity{
		ID:   ins.ID().String(),
		Name: ins.(*testEntity).Name(),
	}

	// to entity:
	retIns, retInsErr := met.ToEntity()(nil, data)
	if retInsErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retInsErr.Error())
		return
	}

	// compare:
	castedIns := ins.(*testEntity)
	castedRetIns := retIns.(*testEntity)

	if !reflect.DeepEqual(castedIns, castedRetIns) {
		t.Errorf("the entity instances must not be the same instance")
		return
	}

	if bytes.Compare(castedIns.ID().Bytes(), castedRetIns.ID().Bytes()) != 0 {
		t.Errorf("the returned ID is invalid.  Expected: %s, Returned: %s", castedIns.ID().String(), castedRetIns.ID().String())
		return
	}

	if castedIns.Name() != castedRetIns.Name() {
		t.Errorf("the returned name is invalid")
		return
	}

}
