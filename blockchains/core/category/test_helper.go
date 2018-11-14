package category

import (
	"bytes"
	"testing"

	uuid "github.com/satori/go.uuid"
)

// CreateCategoryForTests creates a category instance for tests
func CreateCategoryForTests() Category {
	id := uuid.NewV4()
	name := "Category Name"
	description := "This is the description of the category"
	out := createCategory(&id, name, description)
	return out
}

// CompareCategoriesForTests compare category instances for tests
func CompareCategoriesForTests(t *testing.T, first Category, second Category) {
	if bytes.Compare(first.ID().Bytes(), second.ID().Bytes()) != 0 {
		t.Errorf("the IDs are different, expected: %s, received: %s", first.ID().String(), second.ID().String())
		return
	}

	if first.Name() != second.Name() {
		t.Errorf("the names are different, expected: %s, received: %s", first.Name(), second.Name())
		return
	}

	if first.Description() != second.Description() {
		t.Errorf("the descriptions are different, expected: %s, received: %s", first.Description(), second.Description())
		return
	}
}
