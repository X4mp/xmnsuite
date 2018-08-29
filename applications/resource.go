package applications

import (
	crypto "github.com/tendermint/tendermint/crypto"
)

/*
 * ResourcePointer
 */

type resourcePointer struct {
	Frm crypto.PubKey `json:"from"`
	Pth string        `json:"path"`
}

func createResourcePointer(from crypto.PubKey, path string) ResourcePointer {
	out := resourcePointer{
		Frm: from,
		Pth: path,
	}

	return &out
}

// From returns the requester's public key
func (obj *resourcePointer) From() crypto.PubKey {
	return obj.Frm
}

// Path represents the resource path
func (obj *resourcePointer) Path() string {
	return obj.Pth
}

// Hash represents the resource hash
func (obj *resourcePointer) Hash() []byte {
	return createResourceHash(obj)
}

/*
 * Resource
 */

type resource struct {
	Ptr ResourcePointer `json:"pointer"`
	Dat []byte          `json:"data"`
}

func createResource(ptr ResourcePointer, data []byte) Resource {
	out := resource{
		Ptr: ptr,
		Dat: data,
	}

	return &out
}

// Pointer returns the resource pointer
func (obj *resource) Pointer() ResourcePointer {
	return obj.Ptr
}

// Data returns the resource data
func (obj *resource) Data() []byte {
	return obj.Dat
}

// Hash returns the hash
func (obj *resource) Hash() []byte {
	return createResourceHash(obj)
}
