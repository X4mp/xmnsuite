package objects

import (
	uuid "github.com/satori/go.uuid"
	amino "github.com/tendermint/go-amino"
)

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	codec.RegisterConcrete(objFortests{}, "objects.objFortests", nil)
}

type objFortests struct {
	ID          *uuid.UUID `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
}

func createObjForTests() *objFortests {
	id, _ := uuid.NewV4()
	out := objFortests{
		ID:          &id,
		Name:        "My name",
		Description: "This is a simple description",
	}

	return &out
}
