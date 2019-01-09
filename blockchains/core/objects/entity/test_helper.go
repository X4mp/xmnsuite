package entity

import (
	"errors"
	"reflect"
	"testing"

	uuid "github.com/satori/go.uuid"
)

// CreateRepresentationForTests creates a Representation instance for tests
func CreateRepresentationForTests() Representation {
	met := CreateMetaDataForTests()
	keynames := func(ins Entity) ([]string, error) {
		return []string{
			"first",
			"second",
		}, nil
	}

	toData := func(ins Entity) (interface{}, error) {
		obj := ins.(*testEntity)
		return &storableTestEntity{
			ID:   obj.ID().String(),
			Name: obj.Name(),
		}, nil
	}

	out, _ := createRepresentation(met, toData, keynames, nil, nil)
	return out
}

// CreateMetaDataForTests creates a MetaData instance for tests
func CreateMetaDataForTests() MetaData {
	name := "MyEntity"
	empStorable := new(storableTestEntity)
	empNormalized := new(storableTestEntity)

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

	// execute:
	met, _ := createMetaData(name, toEntity, normalize, denormalize, empStorable, empNormalized)
	return met
}

// CompareEntityPartialSetForTests compares EntityPartialSet instances for tests
func CompareEntityPartialSetForTests(t *testing.T, ps PartialSet, set []Entity, index int, totalAmount int) {
	if ps.Index() != index {
		t.Errorf("the index is invalid.  Expected: %d, returned: %d", index, ps.Index())
		return
	}

	amount := len(set)
	if ps.Amount() != amount {
		t.Errorf("the amount is invalid.  Expected: %d, returned: %d", amount, ps.Amount())
		return
	}

	if ps.TotalAmount() != totalAmount {
		t.Errorf("the totalAmount is invalid.  Expected: %d, returned: %d", totalAmount, ps.TotalAmount())
		return
	}

	if !reflect.DeepEqual(ps.Instances(), set) {
		t.Errorf("the instances do not match")
		return
	}
}
