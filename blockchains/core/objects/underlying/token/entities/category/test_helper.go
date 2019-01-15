package category

import (
	"reflect"
	"testing"

	uuid "github.com/satori/go.uuid"
)

// CreateCategoryForTests creates a category for tests
func CreateCategoryForTests() Category {
	id := uuid.NewV4()
	out, outErr := createCategory(&id, "My Category", "This is the category description")
	if outErr != nil {
		panic(outErr)
	}

	return out
}

// CreateCategoryWithParentForTests creates a category with parent for tests
func CreateCategoryWithParentForTests(par Category) Category {
	id := uuid.NewV4()
	out, outErr := createCategoryWithParent(&id, "My Category", "This is the category description", par)
	if outErr != nil {
		panic(outErr)
	}

	return out
}

// CompareCategoriesForTests compare Category instances for tests
func CompareCategoriesForTests(t *testing.T, first Category, second Category) {
	if !reflect.DeepEqual(first.ID(), second.ID()) {
		t.Errorf("the IDs are different.  Expected: %s, Returned: %s", first.ID().String(), second.ID().String())
		return
	}

	if first.Title() != second.Title() {
		t.Errorf("the title is different.  Expected: %s, Returned: %s", first.Title(), second.Title())
		return
	}

	if first.Description() != second.Description() {
		t.Errorf("the description is different.  Expected: %s, Returned: %s", first.Description(), second.Description())
		return
	}

	if first.HasParent() != second.HasParent() {
		t.Errorf("one of the category instance have a parent, the other not")
		return
	}

	if first.HasParent() {
		CompareCategoriesForTests(t, first.Parent(), second.Parent())
	}
}
