package work

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/oelmekki/matrix"
)

type work struct {
	TS  time.Time       `json:"timestamp"`
	IN  []matrix.Matrix `json:"input"`
	OUT []matrix.Matrix `json:"output"`
}

func createWork(ts time.Time, in []matrix.Matrix, out []matrix.Matrix) (Work, error) {
	wk := work{
		TS:  ts,
		IN:  in,
		OUT: out,
	}

	return &wk, nil
}

func fromNormalizedToWork(ins Normalized) (Work, error) {
	if normalized, ok := ins.(*normalizedWork); ok {
		input, inputErr := fromStringToFloatMatrix(normalized.Input)
		if inputErr != nil {
			return nil, inputErr
		}

		output, outputErr := fromStringToFloatMatrix(normalized.Output)
		if outputErr != nil {
			return nil, outputErr
		}

		out, outErr := createWork(normalized.Timestamp, input, output)
		if outErr != nil {
			return nil, outErr
		}

		return out, nil
	}

	str := fmt.Sprintf("the given normalized instance is not a valid Work")
	return nil, errors.New(str)
}

func fromStringToFloatMatrix(in [][][]string) ([]matrix.Matrix, error) {
	out := []matrix.Matrix{}
	for index, oneMatrix := range in {
		amountRows := len(oneMatrix)
		amountCols := len(oneMatrix[0])
		out = append(out, matrix.GenerateMatrix(amountRows, amountCols))
		for rowsIndex, oneRow := range oneMatrix {
			for colsIndex, oneElement := range oneRow {
				elementAsFloat, elementAsFloatErr := strconv.ParseFloat(oneElement, 64)
				if elementAsFloatErr != nil {
					return nil, elementAsFloatErr
				}

				out[index].SetAt(rowsIndex, colsIndex, elementAsFloat)

			}
		}
	}

	return out, nil
}

// Timestamp returns the timestamp
func (obj *work) Timestamp() time.Time {
	return obj.TS
}

// Input returns the input matrix
func (obj *work) Input() []matrix.Matrix {
	return obj.IN
}

// Output returns the output matrix
func (obj *work) Output() []matrix.Matrix {
	return obj.OUT
}

// PartialVerify verifies partially the output
func (obj *work) PartialVerify(amount int) error {
	isInList := func(element int, indexes []int) bool {
		for _, oneIndex := range indexes {
			if oneIndex == element {
				return true
			}
		}

		return false
	}

	getUniqueRandIndexes := func(amount int, maxSize int) []int {
		indexes := []int{}
		for i := 0; i < amount; i++ {
			isThere := true
			var idx int
			for isThere {
				idx = rand.Int() % maxSize
				isThere = isInList(idx, indexes)
			}

			indexes = append(indexes, idx)
		}

		return indexes
	}

	if amount < len(obj.IN) {
		str := fmt.Sprintf("the amount of elements to verify (%d) cannot be bigger than the input matrix array (lenght: %d)", amount, len(obj.IN))
		return errors.New(str)
	}

	// verify each matrix:
	ts := obj.Timestamp()
	indexes := getUniqueRandIndexes(amount, len(obj.IN))
	for _, oneIndex := range indexes {
		matxErr := verifyMatrix(ts, obj.IN[oneIndex], obj.OUT[oneIndex])
		if matxErr != nil {
			str := fmt.Sprintf("there was an error wile veryfing the input matrix (index: %d): %s", oneIndex, matxErr.Error())
			return errors.New(str)
		}
	}

	return nil
}
