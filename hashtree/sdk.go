package hashtree

import "errors"

//Errors:
var errBlockIsMandatory = errors.New("the blocks are mandatory")

// CreateHashTreeParams represents the CreateHashTree params
type CreateHashTreeParams struct {
	Blocks [][]byte
	JS     []byte
}

// SDKFunc represents the public func of the hashtree
var SDKFunc = struct {
	CreateHashTree func(params CreateHashTreeParams) HashTree
}{
	CreateHashTree: func(params CreateHashTreeParams) HashTree {
		if params.JS != nil {
			ptr := new(hashTree)
			jsErr := cdc.UnmarshalJSON(params.JS, ptr)
			if jsErr != nil {
				panic(jsErr)
			}

			return ptr
		}
		if params.Blocks == nil {
			panic(errBlockIsMandatory)
		}

		if len(params.Blocks) <= 0 {
			panic(errBlockIsMandatory)
		}

		out, outErr := createHashTreeFromBlocks(params.Blocks)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
}

// Hash represents a single hash
type Hash interface {
	String() string
	Get() []byte
	Compare(h Hash) bool
}

// Block represents a block of hashes
type Block interface {
	Leaves() Leaves
	HashTree() HashTree
}

// ParentLeaf represents an hashtree parent leaf
type ParentLeaf interface {
	Left() Leaf
	Right() Leaf
	BlockLeaves() Leaves
	HashTree() HashTree
}

// Leaf represents an hashtree leaf
type Leaf interface {
	Head() Hash
	HasParent() bool
	Parent() ParentLeaf
	Leaves() Leaves
	Height() int
}

// Leaves represents a list of Leaf instances
type Leaves interface {
	Leaves() []Leaf
	Merge(lves Leaves) Leaves
	HashTree() HashTree
}

// Compact represents a compact hashtree
type Compact interface {
	Head() Hash
	Leaves() Leaves
	Length() int
}

// HashTree represents an hashtree
type HashTree interface {
	Height() int
	Length() int
	Head() Hash
	Parent() ParentLeaf
	Compact() Compact
	Order(data [][]byte) ([][]byte, error)
}
