package hashtree

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"testing"
	"time"

	helpers "github.com/xmnservices/xmnsuite/helpers"
	convert "github.com/xmnservices/xmnsuite/tests"
)

// we must also split data, create a tree, create a compact tree, and pass the shuffled data to it, to get it back in order
// when passing an invalid amount of blocks to the CreateHashTree, returns an error (1, for example.)
func createTreeAndTest(t *testing.T, text string, delimiter string, height int) {

	shuf := func(v [][]byte) {
		f := reflect.Swapper(v)
		n := len(v)
		r := rand.New(rand.NewSource(time.Now().Unix()))
		for i := 0; i < n; i++ {
			f(r.Intn(n), r.Intn(n))
		}
	}

	splittedData := bytes.Split([]byte(text), []byte(delimiter))
	splittedDataLength := len(splittedData)
	splittedDataLengthPowerOfTwo := int(math.Pow(2, math.Ceil(math.Log(float64(splittedDataLength))/math.Log(2))))
	tree, treeErr := createHashTreeFromBlocks(splittedData)

	if tree == nil {
		t.Errorf("the returned instance was expected to be an instance, nil returned")
		return
	}

	if treeErr != nil {
		t.Errorf("the returned error was expected to be nil, valid error returned: %s", treeErr.Error())
		return
	}

	secondTree, secondTreeErr := createHashTreeFromBlocks(splittedData)
	if secondTreeErr != nil {
		t.Errorf("the returned error was expected to be nil, valid error returned: %s", secondTreeErr.Error())
		return
	}

	if tree.Head().String() != secondTree.Head().String() {
		t.Errorf("the tree hashes changed even if they were build with the same data: First: %s, Second: %s", tree.Head().String(), secondTree.Head().String())
		return
	}

	treeHeight := tree.Height()
	if treeHeight != height {
		t.Errorf("the binary tree's height should be %d because it contains %d data blocks, %d given", height, len(splittedData), treeHeight)
		return
	}

	treeLength := tree.Length()
	if treeLength != splittedDataLengthPowerOfTwo {
		t.Errorf("the HashTree should have a length of %d, %d given", splittedDataLengthPowerOfTwo, treeLength)
		return
	}

	compact := tree.Compact()
	compactLength := compact.Length()
	if splittedDataLengthPowerOfTwo != compactLength {
		t.Errorf("the CompactHashTree should have a length of %d, %d given", splittedDataLengthPowerOfTwo, compactLength)
		return
	}

	if !tree.Head().Compare(compact.Head()) {
		t.Errorf("the HashTree root hash: %x is not the same as the CompactHashTree root hash: %x", tree.Head().Get(), compact.Head().Get())
		return
	}

	shuffledData := make([][]byte, len(splittedData))
	copy(shuffledData, splittedData)
	shuf(shuffledData)

	reOrderedSplittedData, reOrderedSplittedDataErr := tree.Order(shuffledData)
	if reOrderedSplittedDataErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", reOrderedSplittedDataErr.Error())
		return
	}

	if !reflect.DeepEqual(splittedData, reOrderedSplittedData) {
		t.Errorf("the re-ordered data is invalid")
		return
	}
}

func TestCreateHashTree_Success(t *testing.T) {
	createTreeAndTest(t, "this|is", "|", 2)                                                                                                                       //2 blocks
	createTreeAndTest(t, "this|is|some|data|separated|by|delimiters|asfsf", "|", 4)                                                                               //8 blocks
	createTreeAndTest(t, "this|is|some|data|separated|by|delimiters|asfsf|another", "|", 5)                                                                       //9 blocks, rounded up to 16
	createTreeAndTest(t, "this|is|some|data|separated|by|delimiters|asfsf|another|lol", "|", 5)                                                                   //10 blocks, rounded up to 16
	createTreeAndTest(t, "this|is|some|data|separated|by|delimiters|asfsf|asfasdf|asdfasdf|asdfasdf|asdfasdf|asdfasdf|asdfasdf|asdfasfd|sdfasd", "|", 5)          //16 blocks
	createTreeAndTest(t, "this|is|some|data|separated|by|delimiters|asfsf|asfasdf|asdfasdf|asdfasdf|asdfasdf|asdfasdf|asdfasdf|asdfasfd|sdfasd|dafgsagf", "|", 6) //17 blocks, rounded up to 32
}

func TestCreateHashTree_withOneBlock_returnsError(t *testing.T) {

	//variables:
	text := "this"
	delimiter := "|"

	splittedData := bytes.Split([]byte(text), []byte(delimiter))
	tree, treeErr := createHashTreeFromBlocks(splittedData)
	orderedData, orderedDataErr := tree.Order(splittedData)

	if orderedDataErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", orderedDataErr.Error())
		return
	}

	if !reflect.DeepEqual(orderedData, splittedData) {
		t.Errorf("the ordered data was invalid")
		return
	}

	if treeErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned")
		return
	}

	if tree.Height() != 2 {
		t.Errorf("the height of the tree was edxpected to be 2, %d returned", tree.Height())
		return
	}

	if tree.Length() != 2 {
		t.Errorf("the length of the tree was edxpected to be 2, %d returned", tree.Length())
		return
	}
}

func TestCreate_convertToJSON_convertToBinary_backAndForth_Success(t *testing.T) {

	//variables:
	jsEmpty := new(hashTree)
	binEmpty := new(hashTree)
	r := rand.New(rand.NewSource(99))
	blks := [][]byte{
		[]byte("this"),
		[]byte("is"),
		[]byte("some"),
		[]byte("blocks"),
		[]byte(fmt.Sprintf("some rand number to make it unique: %d", r.Int())),
	}

	//execute:
	h, htErr := createHashTreeFromBlocks(blks)
	if htErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", htErr.Error())
		return
	}

	// convert with amino:
	convert.ConvertToJSON(t, h, jsEmpty, cdc)
	convert.ConvertToBinary(t, h, binEmpty, cdc)

	// convert with GOB:
	data, dataErr := helpers.GetBytes(h)
	if dataErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", dataErr.Error())
		return
	}

	ptr := new(hashTree)
	gobErr := helpers.Marshal(data, ptr)
	if gobErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", gobErr.Error())
		return
	}

	if !h.Head().Compare(ptr.Head()) {
		t.Errorf("there was an error while converting the hashtree backandforth using gob")
		return
	}
}
