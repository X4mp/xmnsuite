package work

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/oelmekki/matrix"
)

const (
	minElement = 1
	maxElement = 10000
)

var hashValues = map[string]int{
	"a": 1,
	"b": 2,
	"c": 3,
	"d": 4,
	"e": 5,
	"f": 6,
	"g": 7,
	"h": 8,
	"i": 9,
	"j": 10,
	"k": 11,
	"l": 12,
	"m": 13,
	"n": 14,
	"o": 15,
	"p": 16,
	"q": 17,
	"r": 18,
	"s": 19,
	"t": 20,
	"u": 21,
	"v": 22,
	"w": 23,
	"x": 24,
	"y": 25,
	"z": 26,
	"0": 27,
	"1": 28,
	"2": 29,
	"3": 30,
	"4": 31,
	"5": 32,
	"6": 33,
	"7": 34,
	"8": 35,
	"9": 36,
}

func verifyMatrix(ts time.Time, input matrix.Matrix, output matrix.Matrix) error {
	// generate then compare:
	hashedMatrix, hashedMatrixErr := generate(ts, input)
	if hashedMatrixErr != nil {
		return hashedMatrixErr
	}

	if !output.EqualTo(hashedMatrix) {
		return errors.New("the output matrix is invalid")
	}

	return nil
}

func generate(ts time.Time, input matrix.Matrix) (matrix.Matrix, error) {
	// make sure the matrix contains numbers between our boundaries:
	rows := input.Rows()
	cols := input.Cols()
	for rowIndex := 0; rowIndex < rows; rowIndex++ {
		for colsIndex := 0; colsIndex < cols; colsIndex++ {
			element := input.At(rowIndex, colsIndex)
			if element < minElement {
				str := fmt.Sprintf("the element (%f) at index (row: %d, cols: %d) is smaller than the minimum: %d", element, rowIndex, colsIndex, minElement)
				return nil, errors.New(str)
			}

			if element >= maxElement {
				str := fmt.Sprintf("the element (%f) at index (row: %d, cols: %d) is bigger than the maximum: %d", element, rowIndex, colsIndex, maxElement-1)
				return nil, errors.New(str)
			}
		}
	}

	// multiply the scalar:
	scalarMultipliedMatrix, scalarMultipliedMatrixErr := input.ScalarMultiply(float64(ts.Second()))
	if scalarMultipliedMatrixErr != nil {
		return nil, scalarMultipliedMatrixErr
	}

	// execute the dot product:
	dotProductedMatrix, dotProductedMatrixErr := input.DotProduct(scalarMultipliedMatrix)
	if dotProductedMatrixErr != nil {
		return nil, dotProductedMatrixErr
	}

	// hash each number:
	hashedRows := dotProductedMatrix.Rows()
	hashedCols := dotProductedMatrix.Cols()
	hashedMatrix := matrix.GenerateMatrix(hashedRows, hashedCols)
	for rowIndex := 0; rowIndex < hashedRows; rowIndex++ {
		for colsIndex := 0; colsIndex < hashedCols; colsIndex++ {
			// retrieve the element:
			element := dotProductedMatrix.At(rowIndex, colsIndex)

			// hash the element:
			h := sha256.New()
			h.Write([]byte(strconv.FormatFloat(element, 'f', 5, 64)))
			hashAsBytes := h.Sum(nil)

			// convert the hash to string:
			str := hex.EncodeToString(hashAsBytes)

			// for each letter, get the number:
			value := 0
			for _, cr := range str {
				keyname := fmt.Sprintf("%c", cr)
				if oneValue, ok := hashValues[keyname]; ok {
					value += oneValue
				}
			}

			// add the number in the matrix:
			hashedMatrix.SetAt(rowIndex, colsIndex, float64(value))

		}
	}

	// return:
	return hashedMatrix, nil
}
