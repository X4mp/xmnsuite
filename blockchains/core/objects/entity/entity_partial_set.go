package entity

import (
	"errors"
	"fmt"
)

type entityPartialSet struct {
	Ins   []Entity `json:"entities"`
	Idx   int      `json:"index"`
	TotAm int      `json:"total_amount"`
}

func createEntityPartialSet(ins []Entity, index int, totalAmount int) (PartialSet, error) {
	if index < 0 {
		str := fmt.Sprintf("the index (%d) cannot be smaller than 0", index)
		return nil, errors.New(str)
	}

	minAmount := (index + len(ins))
	if totalAmount < minAmount {
		str := fmt.Sprintf("the totalAmount (%d) cannot be smaller than the index + the length of the instances (%d)", totalAmount, minAmount)
		return nil, errors.New(str)
	}

	out := entityPartialSet{
		Ins:   ins,
		Idx:   index,
		TotAm: totalAmount,
	}

	return &out, nil
}

func createEntityPartialSetFromNormalized(normalizedPS *normalizedPartialSet, metaData MetaData) (PartialSet, error) {
	denormalizedEntities := []Entity{}
	for _, oneNormalized := range normalizedPS.Ins {
		denormalized, denormalizedErr := metaData.Denormalize()(oneNormalized)
		if denormalizedErr != nil {
			return nil, denormalizedErr
		}

		denormalizedEntities = append(denormalizedEntities, denormalized)
	}

	return createEntityPartialSet(denormalizedEntities, normalizedPS.Idx, normalizedPS.TotAm)
}

// Instances returns the instances
func (obj *entityPartialSet) Instances() []Entity {
	return obj.Ins
}

// Index returns the index
func (obj *entityPartialSet) Index() int {
	return obj.Idx
}

// Amount returns the amount
func (obj *entityPartialSet) Amount() int {
	return len(obj.Ins)
}

// TotalAmount returns the totalAmount
func (obj *entityPartialSet) TotalAmount() int {
	return obj.TotAm
}

// IsLast returns true if this is the last element sof the partial set, false otherwise
func (obj *entityPartialSet) IsLast() bool {
	return (obj.Index() + obj.Amount()) >= obj.TotalAmount()
}
