package keys

import (
	"encoding/gob"

	"github.com/xmnservices/xmnsuite/hashtree"
)

func init() {
	RegisterGob()
}

// RegisterGob registers the hashtree for gob
func RegisterGob() {
	hashtree.RegisterGob()
	gob.Register(&storedInstance{})
	gob.Register(&concreteKeys{})
}
