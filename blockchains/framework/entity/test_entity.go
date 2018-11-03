package entity

import (
	"crypto/rand"
	"fmt"
	mathrand "math/rand"

	uuid "github.com/satori/go.uuid"
)

type testEntity struct {
	UUID *uuid.UUID `json:"id"`
	Nam  string     `json:"name"`
}

type storableTestEntity struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func createTestEntityForTests() Entity {
	randName := func(n int) string {
		b := make([]byte, n)
		if _, err := rand.Read(b); err != nil {
			panic(err)
		}
		s := fmt.Sprintf("%X", b)
		return s
	}

	id := uuid.NewV4()
	name := randName(mathrand.Int() % 100)
	return createTestEntity(&id, name)
}

func createTestEntity(id *uuid.UUID, name string) Entity {
	out := testEntity{
		UUID: id,
		Nam:  name,
	}

	return &out
}

// ID returns the ID
func (obj *testEntity) ID() *uuid.UUID {
	return obj.UUID
}

// Name returns the name
func (obj *testEntity) Name() string {
	return obj.Nam
}
