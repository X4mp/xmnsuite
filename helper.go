package datamint

import (
	"bytes"
	"encoding/gob"
)

/*
 * Helper func
 */

// GetBytes returns the []byte of any interface{}
func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Marshal marshals the []byte into the pointer
func Marshal(data []byte, ptr interface{}) error {
	buf := bytes.NewReader(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(ptr)
	if err != nil {
		return err
	}

	return nil
}
