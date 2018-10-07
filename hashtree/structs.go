package hashtree

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
)

/*
 * Hash
 */

type hash struct {
	H []byte `json:"hash"`
}

func createHashFromString(str string) (Hash, error) {
	dec, decErr := hex.DecodeString(str)
	if decErr != nil {
		return nil, decErr
	}

	out := hash{
		H: dec,
	}

	return &out, nil
}

func createHashFromData(data []byte) Hash {
	sha := sha256.New()
	sha.Write(data)

	out := hash{
		H: sha.Sum(nil),
	}

	return &out
}

// String returns a string representation of the hash
func (obj *hash) String() string {
	return hex.EncodeToString(obj.H)
}

// Get returns the hash as byte
func (obj *hash) Get() []byte {
	return obj.H
}

// Compare compares the hashes.  If equal, returns true, otherwise false
func (obj *hash) Compare(h Hash) bool {
	return bytes.Compare(obj.H, h.Get()) == 0
}

/*
* Block
 */

type block struct {
	List []Hash `json:"list"`
}

func createBlockFromData(data [][]byte) (Block, error) {

	if len(data) <= 1 {
		data = append(data, []byte(""))
	}

	hashes := []Hash{}
	for _, oneData := range data {
		oneHash := createHashFromData(oneData)
		hashes = append(hashes, oneHash)
	}

	blk := block{
		List: hashes,
	}

	return blk.resize(), nil
}

func (obj *block) resize() Block {
	//need to make sure the elements are always a power of 2:
	isPowerOfTwo := obj.isLengthPowerForTwo()
	if !isPowerOfTwo {
		obj.resizeToNextPowerOfTwo()
	}

	return obj
}

func (obj *block) isLengthPowerForTwo() bool {
	length := len(obj.List)
	return (length != 0) && ((length & (length - 1)) == 0)
}

func (obj *block) resizeToNextPowerOfTwo() Block {
	lengthAsFloat := float64(len(obj.List))
	next := uint(math.Pow(2, math.Ceil(math.Log(lengthAsFloat)/math.Log(2))))
	remaining := int(next) - int(lengthAsFloat)
	for i := 0; i < remaining; i++ {
		single := createHashFromData(nil)
		obj.List = append(obj.List, single)
	}

	return obj
}

// Leaves returns the leaves of the block
func (obj *block) Leaves() Leaves {
	leaves := []Leaf{}
	for _, oneBlockHash := range obj.List {
		oneLeaf := createLeaf(oneBlockHash)
		leaves = append(leaves, oneLeaf)
	}

	return createLeaves(leaves)
}

// HashTree returns the HashTree
func (obj *block) HashTree() HashTree {
	leaves := obj.Leaves()
	tree := leaves.HashTree()
	return tree
}

// MarshalJSON converts the instance to JSON
func (obj *block) MarshalJSON() ([]byte, error) {
	list := []string{}
	for _, oneHash := range obj.List {
		list = append(list, oneHash.String())
	}

	js, jsErr := json.Marshal(list)
	if jsErr != nil {
		return nil, jsErr
	}

	return js, nil
}

// UnmarshalJSON converts the JSON to an instance
func (obj *block) UnmarshalJSON(data []byte) error {
	hashes := new([]string)
	jsErr := json.Unmarshal(data, hashes)
	if jsErr != nil {
		return jsErr
	}

	out := []Hash{}
	for _, oneHashAsString := range *hashes {
		oneHash, oneHashErr := createHashFromString(oneHashAsString)
		if oneHashErr != nil {
			return oneHashErr
		}

		out = append(out, oneHash)
	}

	obj.List = out
	return nil
}

/*
 * ParentLeaf
 */

type parentLeaf struct {
	Lft Leaf `json:"left"`
	Rgt Leaf `json:"right"`
}

func createParentLeaf(left Leaf, right Leaf) ParentLeaf {
	out := parentLeaf{
		Lft: left,
		Rgt: right,
	}

	return &out
}

// HashTree returns the hashtree
func (obj *parentLeaf) HashTree() HashTree {
	data := bytes.Join([][]byte{
		obj.Left().Head().Get(),
		obj.Right().Head().Get(),
	}, []byte{})

	hash := createHashFromData(data)
	out := createHashTree(hash, obj)
	return out
}

// BlockLeaves returns the block leaves
func (obj *parentLeaf) BlockLeaves() Leaves {
	left := obj.Left()
	right := obj.Right()
	leftLeaves := left.Leaves()
	rightLeaves := right.Leaves()
	return leftLeaves.Merge(rightLeaves)
}

// Left returns the left leaf
func (obj *parentLeaf) Left() Leaf {
	return obj.Lft
}

// Right returns the right leaf
func (obj *parentLeaf) Right() Leaf {
	return obj.Rgt
}

/*
* Leaf
 */

type leaf struct {
	Hd Hash       `json:"head"`
	Pt ParentLeaf `json:"parent"`
}

func createLeaf(head Hash) Leaf {
	out := leaf{
		Hd: head,
		Pt: nil,
	}

	return &out
}

func createLeafWithParent(head Hash, parent ParentLeaf) Leaf {
	out := leaf{
		Hd: head,
		Pt: parent,
	}

	return &out
}

func createChildLeaf(left Leaf, right Leaf) Leaf {
	data := bytes.Join([][]byte{
		left.Head().Get(),
		right.Head().Get(),
	}, []byte{})

	h := createHashFromData(data)
	out := createLeaf(h)
	return out
}

// Head returns the head hash
func (obj *leaf) Head() Hash {
	return obj.Hd
}

// HasParent returns true if there is a parent, false otherwise
func (obj *leaf) HasParent() bool {
	return obj.Pt != nil
}

// Parent returns the parent, if any
func (obj *leaf) Parent() ParentLeaf {
	return obj.Pt
}

// Leaves returns the leaves
func (obj *leaf) Leaves() Leaves {
	if obj.HasParent() {
		return obj.Parent().BlockLeaves()
	}

	leaves := []Leaf{
		obj,
	}

	output := createLeaves(leaves)
	return output
}

// Height returns the leaf height
func (obj *leaf) Height() int {
	cpt := 0
	var oneLeaf Leaf
	for {

		if oneLeaf == nil {
			oneLeaf = obj
		}

		if !oneLeaf.HasParent() {
			return cpt
		}

		cpt++
		oneLeaf = oneLeaf.Parent().Left()
	}
}

/*
* Leaves
 */

type leaves struct {
	Lst []Leaf `json:"leaves"`
}

func createLeaves(list []Leaf) Leaves {
	out := leaves{
		Lst: list,
	}

	return &out
}

// Leaves returns the leaves
func (obj *leaves) Leaves() []Leaf {
	return obj.Lst
}

// Merge merge Leaves instances
func (obj *leaves) Merge(lves Leaves) Leaves {
	for _, oneLeaf := range lves.Leaves() {
		obj.Lst = append(obj.Lst, oneLeaf)
	}

	return obj
}

// HashTree returns the hashtree
func (obj *leaves) HashTree() HashTree {
	length := len(obj.Lst)
	if length == 2 {
		left := obj.Lst[0]
		right := obj.Lst[1]
		parent := createParentLeaf(left, right)
		tree := parent.HashTree()
		return tree
	}

	childrenLeaves := obj.createChildrenLeaves()
	tree := childrenLeaves.HashTree()
	return tree
}

func (obj *leaves) createChildrenLeaves() Leaves {
	var childrenLeaves []Leaf
	for index, oneLeaf := range obj.Lst {

		if index%2 != 0 {
			continue
		}

		left := oneLeaf
		right := obj.Lst[index+1]
		child := createChildLeaf(left, right)
		parent := createParentLeaf(left, right)
		childWithParent := createLeafWithParent(child.Head(), parent)
		childrenLeaves = append(childrenLeaves, childWithParent)
	}

	return createLeaves(childrenLeaves)
}

/*
* Compact
 */

type compact struct {
	Hd   Hash   `json:"head"`
	Lves Leaves `json:"leaves"`
}

func createCompact(head Hash, leaves Leaves) Compact {
	out := compact{
		Hd:   head,
		Lves: leaves,
	}

	return &out
}

// Head returns the head hash
func (obj *compact) Head() Hash {
	return obj.Hd
}

// Leaves returns the leaves
func (obj *compact) Leaves() Leaves {
	return obj.Lves
}

// Length returns the length of the compact hashtree
func (obj *compact) Length() int {
	return len(obj.Lves.Leaves())
}

/*
* HashTree
 */

// hashTree represents a concrete HashTree implementation
type hashTree struct {
	Hd Hash       `json:"head"`
	Pt ParentLeaf `json:"parent"`
}

func createHashTree(head Hash, parent ParentLeaf) HashTree {
	out := hashTree{
		Hd: head,
		Pt: parent,
	}

	return &out
}

func createHashTreeFromBlocks(blocks [][]byte) (HashTree, error) {
	blockHashes, blockHashesErr := createBlockFromData(blocks)
	if blockHashesErr != nil {
		return nil, blockHashesErr
	}

	tree := blockHashes.HashTree()
	return tree, nil
}

// Height returns the hashtree height
func (obj *hashTree) Height() int {
	left := obj.Pt.Left()
	return left.Height() + 2
}

// Length returns the hashtree length
func (obj *hashTree) Length() int {
	blockLeaves := obj.Pt.BlockLeaves()
	return len(blockLeaves.Leaves())
}

// Head returns the head hash
func (obj *hashTree) Head() Hash {
	return obj.Hd
}

// Parent returns the parent leaf
func (obj *hashTree) Parent() ParentLeaf {
	return obj.Pt
}

// Compact returns the compact version of the hashtree
func (obj *hashTree) Compact() Compact {
	blockLeaves := obj.Pt.BlockLeaves()
	return createCompact(obj.Hd, blockLeaves)
}

// Order orders data that matches the leafs of the HashTree
func (obj *hashTree) Order(data [][]byte) ([][]byte, error) {
	hashed := map[string][]byte{}
	for _, oneData := range data {
		sha := sha256.New()
		sha.Write(oneData)
		hashAsString := hex.EncodeToString(sha.Sum(nil))
		hashed[hashAsString] = oneData
	}

	out := [][]byte{}
	leaves := obj.Pt.BlockLeaves().Leaves()
	for _, oneLeaf := range leaves {
		LeafHashAsString := oneLeaf.Head().String()
		if oneData, ok := hashed[LeafHashAsString]; ok {
			out = append(out, oneData)
			continue
		}

		//must be a filling Leaf, so continue:
		continue
	}

	if len(out) != len(data) {
		str := fmt.Sprintf("the length of the input data (%d) does not match the length of the output (%d), therefore, some data blocks could not be found in the hash leaves", len(data), len(out))
		return nil, errors.New(str)
	}

	return out, nil
}
