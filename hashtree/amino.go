package hashtree

import (
	amino "github.com/tendermint/go-amino"
)

var cdc = amino.NewCodec()

const (

	// XMNSuiteHashTreeHash represents the xmnsuite Hash resource
	XMNSuiteHashTreeHash = "xmnsuite/Hash"

	// XMNSuiteHashTreeBlock represents the xmnsuite Block resource
	XMNSuiteHashTreeBlock = "xmnsuite/Block"

	// XMNSuiteHashTreeParentLeaf represents the xmnsuite ParentLeaf resource
	XMNSuiteHashTreeParentLeaf = "xmnsuite/ParentLeaf"

	// XMNSuiteHashTreeLeaf represents the xmnsuite Leaf resource
	XMNSuiteHashTreeLeaf = "xmnsuite/Leaf"

	// XMNSuiteHashTreeLeaves represents the xmnsuite Leaves resource
	XMNSuiteHashTreeLeaves = "xmnsuite/Leaves"

	// XMNSuiteHashTreeCompact represents the xmnsuite Compact resource
	XMNSuiteHashTreeCompact = "xmnsuite/Compact"

	// XMNSuiteHashTreeHashTree represents the xmnsuite HashTree resource
	XMNSuiteHashTreeHashTree = "xmnsuite/HashTree"
)

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {

	// Hash
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Hash)(nil), nil)
		codec.RegisterConcrete(&hash{}, XMNSuiteHashTreeHash, nil)
	}()

	// Block
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Block)(nil), nil)
		codec.RegisterConcrete(&block{}, XMNSuiteHashTreeBlock, nil)
	}()

	// ParentLeaf
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*ParentLeaf)(nil), nil)
		codec.RegisterConcrete(&parentLeaf{}, XMNSuiteHashTreeParentLeaf, nil)
	}()

	// Leaf
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Leaf)(nil), nil)
		codec.RegisterConcrete(&leaf{}, XMNSuiteHashTreeLeaf, nil)
	}()

	// Leaves
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Leaves)(nil), nil)
		codec.RegisterConcrete(&leaves{}, XMNSuiteHashTreeLeaves, nil)
	}()

	// Compact
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Compact)(nil), nil)
		codec.RegisterConcrete(&compact{}, XMNSuiteHashTreeCompact, nil)
	}()

	// HashTree
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*HashTree)(nil), nil)
		codec.RegisterConcrete(&hashTree{}, XMNSuiteHashTreeHashTree, nil)
	}()
}
