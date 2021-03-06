package helpers

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"fmt"
)

/*
 * Helper func
 */

// GetUniqueValue returns the unique value the map contains.  If the values are not the same, -1 is returned
func GetUniqueValue(elements map[string]int) int {
	lastElement := -1
	for _, oneElement := range elements {
		if lastElement == -1 {
			lastElement = oneElement
			continue
		}

		if lastElement != oneElement {
			return -1
		}
	}

	return lastElement
}

// MakeUnique makes the elements unique
func MakeUnique(elements ...interface{}) []interface{} {
	unique := []interface{}{}
	for _, onElement := range elements {
		isUnique := true
		oneElementAsBytes, oneElementAsBytesErr := GetHash(onElement)
		if oneElementAsBytesErr != nil {
			str := fmt.Sprintf("there was an error while converting an existing element to []byte: %s", oneElementAsBytesErr.Error())
			panic(errors.New(str))
		}

		for _, oneUnique := range unique {
			oneUniqueAsBytes, oneUniqueAsBytesErr := GetHash(oneUnique)
			if oneUniqueAsBytesErr != nil {
				str := fmt.Sprintf("there was an error while converting a unique element to []byte: %s", oneUniqueAsBytesErr.Error())
				panic(errors.New(str))
			}

			if bytes.Compare(oneElementAsBytes, oneUniqueAsBytes) == 0 {
				isUnique = false
				break
			}
		}

		if isUnique {
			unique = append(unique, onElement)
		}
	}

	return unique
}

// GetHash returns the []byte hash of any interface{}
func GetHash(key interface{}) ([]byte, error) {
	data, dataErr := GetBytes(key)
	if dataErr != nil {
		str := fmt.Sprintf("there was an error while converting an interface{} to []byte: %s", dataErr.Error())
		return nil, errors.New(str)
	}

	ha := sha256.New()
	_, err := ha.Write(data)
	if err != nil {
		return nil, err
	}

	return ha.Sum(nil), nil
}

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

// Print prints a value on screen
func Print(value string) {
	fmt.Printf("%s", write(value))
}

func write(str string) string {
	out := fmt.Sprintf("\n************ XMN ************\n")
	out = fmt.Sprintf("%s%s", out, str)
	out = fmt.Sprintf("%s\n********** END XMN **********\n", out)
	return out
}
