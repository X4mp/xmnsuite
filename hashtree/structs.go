package hashtree

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"

	amino "github.com/tendermint/go-amino"
)

var cdc = amino.NewCodec()

/*
* Hash
 */

type concreteHash struct {
	h []byte
}

func createHashFromString(str string) (Hash, error) {
	dec, decErr := hex.DecodeString(str)
	if decErr != nil {
		return nil, decErr
	}

	out := concreteHash{
		h: dec,
	}

	return &out, nil
}

func createHashFromData(data []byte) Hash {
	sha := sha256.New()
	sha.Write(data)

	out := concreteHash{
		h: sha.Sum(nil),
	}

	return &out
}

// String returns a string representation of the hash
func (obj *concreteHash) String() string {
	return hex.EncodeToString(obj.h)
}

// Get returns the hash as byte
func (obj *concreteHash) Get() []byte {
	return obj.h
}

// Compare compares the hashes.  If equal, returns true, otherwise false
func (obj *concreteHash) Compare(h Hash) bool {
	return bytes.Compare(obj.h, h.Get()) == 0
}

/*
* Block
 */

type concreteBlock struct {
	list []Hash
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

	blk := concreteBlock{
		list: hashes,
	}

	return blk.resize(), nil
}

func (obj *concreteBlock) resize() Block {
	//need to make sure the elements are always a power of 2:
	isPowerOfTwo := obj.isLengthPowerForTwo()
	if !isPowerOfTwo {
		obj.resizeToNextPowerOfTwo()
	}

	return obj
}

func (obj *concreteBlock) isLengthPowerForTwo() bool {
	length := len(obj.list)
	return (length != 0) && ((length & (length - 1)) == 0)
}

func (obj *concreteBlock) resizeToNextPowerOfTwo() Block {
	lengthAsFloat := float64(len(obj.list))
	next := uint(math.Pow(2, math.Ceil(math.Log(lengthAsFloat)/math.Log(2))))
	remaining := int(next) - int(lengthAsFloat)
	for i := 0; i < remaining; i++ {
		single := createHashFromData(nil)
		obj.list = append(obj.list, single)
	}

	return obj
}

// Leaves returns the leaves of the block
func (obj *concreteBlock) Leaves() Leaves {
	leaves := []Leaf{}
	for _, oneBlockHash := range obj.list {
		oneLeaf := createLeaf(oneBlockHash)
		leaves = append(leaves, oneLeaf)
	}

	return createLeaves(leaves)
}

// HashTree returns the HashTree
func (obj *concreteBlock) HashTree() HashTree {
	leaves := obj.Leaves()
	tree := leaves.HashTree()
	return tree
}

// MarshalJSON converts the instance to JSON
func (obj *concreteBlock) MarshalJSON() ([]byte, error) {
	list := []string{}
	for _, oneHash := range obj.list {
		list = append(list, oneHash.String())
	}

	js, jsErr := json.Marshal(list)
	if jsErr != nil {
		return nil, jsErr
	}

	return js, nil
}

// UnmarshalJSON converts the JSON to an instance
func (obj *concreteBlock) UnmarshalJSON(data []byte) error {
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

	obj.list = out
	return nil
}

/*
* ParentLeaf
 */

type jsonParentLeaf struct {
	Left  *jsonLeaf `json:"left"`
	Right *jsonLeaf `json:"right"`
}

func createJSONParentLeaf(parent ParentLeaf) *jsonParentLeaf {
	left := createJSONLeaf(parent.Left())
	right := createJSONLeaf(parent.Right())
	out := jsonParentLeaf{
		Left:  left,
		Right: right,
	}

	return &out
}

type concreteParentLeaf struct {
	left  Leaf
	right Leaf
}

func createParentLeaf(left Leaf, right Leaf) ParentLeaf {
	out := concreteParentLeaf{
		left:  left,
		right: right,
	}

	return &out
}

func createParentLeafFromJSON(js *jsonParentLeaf) (ParentLeaf, error) {
	left, leftErr := createLeafFromJSON(js.Left)
	if leftErr != nil {
		return nil, leftErr
	}

	right, rightErr := createLeafFromJSON(js.Right)
	if rightErr != nil {
		return nil, rightErr
	}

	out := createParentLeaf(left, right)
	return out, nil
}

// HashTree returns the hashtree
func (obj *concreteParentLeaf) HashTree() HashTree {
	data := bytes.Join([][]byte{
		obj.Left().Head().Get(),
		obj.Right().Head().Get(),
	}, []byte{})

	hash := createHashFromData(data)
	out := createHashTree(hash, obj)
	return out
}

// BlockLeaves returns the block leaves
func (obj *concreteParentLeaf) BlockLeaves() Leaves {
	left := obj.Left()
	right := obj.Right()
	leftLeaves := left.Leaves()
	rightLeaves := right.Leaves()
	return leftLeaves.Merge(rightLeaves)
}

// Left returns the left leaf
func (obj *concreteParentLeaf) Left() Leaf {
	return obj.left
}

// Right returns the right leaf
func (obj *concreteParentLeaf) Right() Leaf {
	return obj.right
}

// MarshalJSON converts the instance to JSON
func (obj *concreteParentLeaf) MarshalJSON() ([]byte, error) {
	jsonParent := createJSONParentLeaf(obj)
	js, jsErr := json.Marshal(jsonParent)
	if jsErr != nil {
		return nil, jsErr
	}

	return js, nil
}

// UnmarshalJSON converts the JSON to an instance
func (obj *concreteParentLeaf) UnmarshalJSON(data []byte) error {
	jsonParent := new(jsonParentLeaf)
	jsErr := json.Unmarshal(data, jsonParent)
	if jsErr != nil {
		return jsErr
	}

	parent, parentErr := createParentLeafFromJSON(jsonParent)
	if parentErr != nil {
		return parentErr
	}

	obj.left = parent.Left()
	obj.right = parent.Right()
	return nil
}

/*
* Leaf
 */

type jsonLeaf struct {
	Head   string          `json:"head"`
	Parent *jsonParentLeaf `json:"parent"`
}

func createJSONLeaf(leaf Leaf) *jsonLeaf {
	head := leaf.Head().String()
	if !leaf.HasParent() {
		out := jsonLeaf{
			Head:   head,
			Parent: nil,
		}

		return &out
	}

	parent := createJSONParentLeaf(leaf.Parent())
	out := jsonLeaf{
		Head:   head,
		Parent: parent,
	}

	return &out
}

type concreteLeaf struct {
	head   Hash
	parent ParentLeaf
}

func createLeaf(head Hash) Leaf {
	out := concreteLeaf{
		head:   head,
		parent: nil,
	}

	return &out
}

func createLeafFromJSON(js *jsonLeaf) (Leaf, error) {
	head, headErr := createHashFromString(js.Head)
	if headErr != nil {
		return nil, headErr
	}

	if js.Parent == nil {
		out := createLeaf(head)
		return out, nil
	}

	parent, parentErr := createParentLeafFromJSON(js.Parent)
	if parentErr != nil {
		return nil, parentErr
	}

	out := createLeafWithParent(head, parent)
	return out, nil
}

func createLeafWithParent(head Hash, parent ParentLeaf) Leaf {
	out := concreteLeaf{
		head:   head,
		parent: parent,
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
func (obj *concreteLeaf) Head() Hash {
	return obj.head
}

// HasParent returns true if there is a parent, false otherwise
func (obj *concreteLeaf) HasParent() bool {
	return obj.parent != nil
}

// Parent returns the parent, if any
func (obj *concreteLeaf) Parent() ParentLeaf {
	return obj.parent
}

// Leaves returns the leaves
func (obj *concreteLeaf) Leaves() Leaves {
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
func (obj *concreteLeaf) Height() int {
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

// MarshalJSON converts the instance to JSON
func (obj *concreteLeaf) MarshalJSON() ([]byte, error) {
	jsLeaf := createJSONLeaf(obj)
	js, jsErr := json.Marshal(jsLeaf)
	if jsErr != nil {
		return nil, jsErr
	}

	return js, nil
}

// UnmarshalJSON converts the JSON to an instance
func (obj *concreteLeaf) UnmarshalJSON(data []byte) error {
	jsLeaf := new(jsonLeaf)
	jsErr := json.Unmarshal(data, jsLeaf)
	if jsErr != nil {
		return jsErr
	}

	leaf, leafErr := createLeafFromJSON(jsLeaf)
	if leafErr != nil {
		return leafErr
	}

	obj.head = leaf.Head()
	obj.parent = leaf.Parent()
	return nil
}

/*
* Leaves
 */

type jsonLeaves struct {
	List []*jsonLeaf `json:"leaves"`
}

func createJSONLeaves(leaves Leaves) *jsonLeaves {
	list := []*jsonLeaf{}
	for _, oneLeaf := range leaves.Leaves() {
		jsonLeaf := createJSONLeaf(oneLeaf)
		list = append(list, jsonLeaf)
	}

	out := jsonLeaves{
		List: list,
	}

	return &out
}

type concreteLeaves struct {
	list []Leaf
}

func createLeaves(list []Leaf) Leaves {
	out := concreteLeaves{
		list: list,
	}

	return &out
}

func createLeavesFromJSON(jsLeaves *jsonLeaves) (Leaves, error) {
	list := []Leaf{}
	for _, oneJSLeaf := range jsLeaves.List {
		oneLeaf, oneLeafErr := createLeafFromJSON(oneJSLeaf)
		if oneLeafErr != nil {
			return nil, oneLeafErr
		}

		list = append(list, oneLeaf)
	}

	out := createLeaves(list)
	return out, nil
}

// Leaves returns the leaves
func (obj *concreteLeaves) Leaves() []Leaf {
	return obj.list
}

// Merge merge Leaves instances
func (obj *concreteLeaves) Merge(lves Leaves) Leaves {
	for _, oneLeaf := range lves.Leaves() {
		obj.list = append(obj.list, oneLeaf)
	}

	return obj
}

// HashTree returns the hashtree
func (obj *concreteLeaves) HashTree() HashTree {
	length := len(obj.list)
	if length == 2 {
		left := obj.list[0]
		right := obj.list[1]
		parent := createParentLeaf(left, right)
		tree := parent.HashTree()
		return tree
	}

	childrenLeaves := obj.createChildrenLeaves()
	tree := childrenLeaves.HashTree()
	return tree
}

func (obj *concreteLeaves) createChildrenLeaves() Leaves {
	var childrenLeaves []Leaf
	for index, oneLeaf := range obj.list {

		if index%2 != 0 {
			continue
		}

		left := oneLeaf
		right := obj.list[index+1]
		child := createChildLeaf(left, right)
		parent := createParentLeaf(left, right)
		childWithParent := createLeafWithParent(child.Head(), parent)
		childrenLeaves = append(childrenLeaves, childWithParent)
	}

	return createLeaves(childrenLeaves)
}

// MarshalJSON converts the instance to JSON
func (obj *concreteLeaves) MarshalJSON() ([]byte, error) {
	jsLeaves := createJSONLeaves(obj)
	js, jsErr := json.Marshal(jsLeaves)
	if jsErr != nil {
		return nil, jsErr
	}

	return js, nil
}

// UnmarshalJSON converts the JSON to an instance
func (obj *concreteLeaves) UnmarshalJSON(data []byte) error {
	jsLeaves := new(jsonLeaves)
	jsErr := json.Unmarshal(data, jsLeaves)
	if jsErr != nil {
		return jsErr
	}

	leaves, leavesErr := createLeavesFromJSON(jsLeaves)
	if leavesErr != nil {
		return leavesErr
	}

	obj.list = leaves.Leaves()
	return nil
}

/*
* Compact
 */

func createJSONCompact(compact Compact) *JSONCompact {
	head := compact.Head().String()
	jsonLeaves := createJSONLeaves(compact.Leaves())
	out := JSONCompact{
		Head:   head,
		Leaves: jsonLeaves,
	}

	return &out
}

type concreteCompact struct {
	head   Hash
	leaves Leaves
}

func createCompact(head Hash, leaves Leaves) Compact {
	out := concreteCompact{
		head:   head,
		leaves: leaves,
	}

	return &out
}

func createCompactFromJSON(jsCompact *JSONCompact) (Compact, error) {
	head, headErr := createHashFromString(jsCompact.Head)
	if headErr != nil {
		return nil, headErr
	}

	leaves, leavesErr := createLeavesFromJSON(jsCompact.Leaves)
	if leavesErr != nil {
		return nil, leavesErr
	}

	out := createCompact(head, leaves)
	return out, nil

}

// Head returns the head hash
func (obj *concreteCompact) Head() Hash {
	return obj.head
}

// Leaves returns the leaves
func (obj *concreteCompact) Leaves() Leaves {
	return obj.leaves
}

// Length returns the length of the compact hashtree
func (obj *concreteCompact) Length() int {
	return len(obj.leaves.Leaves())
}

// MarshalJSON converts the instance to JSON
func (obj *concreteCompact) MarshalJSON() ([]byte, error) {
	jsCompact := createJSONCompact(obj)
	js, jsErr := json.Marshal(jsCompact)
	if jsErr != nil {
		return nil, jsErr
	}

	return js, nil
}

// UnmarshalJSON converts the JSON to an instance
func (obj *concreteCompact) UnmarshalJSON(data []byte) error {
	jsCompact := new(JSONCompact)
	jsErr := json.Unmarshal(data, jsCompact)
	if jsErr != nil {
		return jsErr
	}

	compact, compactErr := createCompactFromJSON(jsCompact)
	if compactErr != nil {
		return compactErr
	}

	obj.head = compact.Head()
	obj.leaves = compact.Leaves()
	return nil
}

/*
* HashTree
 */

// concreteHashTree represents a concrete HashTree implementation
type concreteHashTree struct {
	head   Hash
	parent ParentLeaf
}

func createHashTree(head Hash, parent ParentLeaf) HashTree {
	out := concreteHashTree{
		head:   head,
		parent: parent,
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
func (obj *concreteHashTree) Height() int {
	left := obj.parent.Left()
	return left.Height() + 2
}

// Length returns the hashtree length
func (obj *concreteHashTree) Length() int {
	blockLeaves := obj.parent.BlockLeaves()
	return len(blockLeaves.Leaves())
}

// Head returns the head hash
func (obj *concreteHashTree) Head() Hash {
	return obj.head
}

// Parent returns the parent leaf
func (obj *concreteHashTree) Parent() ParentLeaf {
	return obj.parent
}

// Compact returns the compact version of the hashtree
func (obj *concreteHashTree) Compact() Compact {
	blockLeaves := obj.parent.BlockLeaves()
	return createCompact(obj.head, blockLeaves)
}

// Order orders data that matches the leafs of the HashTree
func (obj *concreteHashTree) Order(data [][]byte) ([][]byte, error) {
	hashed := map[string][]byte{}
	for _, oneData := range data {
		sha := sha256.New()
		sha.Write(oneData)
		hashAsString := hex.EncodeToString(sha.Sum(nil))
		hashed[hashAsString] = oneData
	}

	out := [][]byte{}
	leaves := obj.parent.BlockLeaves().Leaves()
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

// MarshalJSON converts the instance to JSON
func (obj *concreteHashTree) MarshalJSON() ([]byte, error) {
	jsCompact := createJSONCompact(obj.Compact())
	js, jsErr := json.Marshal(jsCompact)
	if jsErr != nil {
		return nil, jsErr
	}

	return js, nil
}

// UnmarshalJSON converts the JSON to an instance
func (obj *concreteHashTree) UnmarshalJSON(data []byte) error {
	jsCompact := new(JSONCompact)
	jsErr := json.Unmarshal(data, jsCompact)
	if jsErr != nil {
		return jsErr
	}

	compact, compactErr := createCompactFromJSON(jsCompact)
	if compactErr != nil {
		return compactErr
	}

	ht := compact.Leaves().HashTree()
	obj.head = ht.Head()
	obj.parent = ht.Parent()
	return nil
}
