package hashtree

import "encoding/gob"

func init() {
	RegisterGob()
}

// RegisterGob registers the hashtree for gob
func RegisterGob() {
	gob.Register(&hash{})
	gob.Register(&parentLeaf{})
	gob.Register(&leaf{})
	gob.Register(&leaves{})
	gob.Register(&compact{})
	gob.Register(&hashTree{})
}
