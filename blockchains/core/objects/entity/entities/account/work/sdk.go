package work

import (
	"errors"
	"math/rand"
	"time"

	"github.com/oelmekki/matrix"
)

// Work represents the work
type Work interface {
	Timestamp() time.Time
	Input() []matrix.Matrix
	Output() []matrix.Matrix
	PartialVerify(amount int) error
}

// Normalized represents a normalized work
type Normalized interface {
}

// CreateWorkParams represents the CreateWork params
type CreateWorkParams struct {
	TS  time.Time
	IN  []matrix.Matrix
	OUT []matrix.Matrix
}

// GenerateParams represents the Generate params
type GenerateParams struct {
	MatrixSize   int
	MatrixAmount int
	MaxElement   int
}

// SDKFunc represents the Work SDK func
var SDKFunc = struct {
	Create      func(params CreateWorkParams) Work
	Generate    func(params GenerateParams) Work
	Normalize   func(ins Work) Normalized
	Denormalize func(ins interface{}) Work
}{
	Create: func(params CreateWorkParams) Work {
		out, outErr := createWork(params.TS, params.IN, params.OUT)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	Generate: func(params GenerateParams) Work {
		inputs := []matrix.Matrix{}
		for i := 0; i < params.MatrixAmount; i++ {
			input := matrix.GenerateMatrix(params.MatrixSize, params.MatrixSize)
			for rowsIndex := 0; rowsIndex < params.MatrixSize; rowsIndex++ {
				for colsIndex := 0; colsIndex < params.MatrixSize; colsIndex++ {
					element := float64((rand.Int() % (maxElement - minElement)) + minElement)
					input.SetAt(rowsIndex, colsIndex, element)
				}
			}

			inputs = append(inputs, input)
		}

		curTS := time.Now().UTC()
		outputs := []matrix.Matrix{}
		for _, oneInput := range inputs {
			output, outputErr := generate(curTS, oneInput)
			if outputErr != nil {
				panic(outputErr)
			}

			outputs = append(outputs, output)
		}

		out, outErr := createWork(curTS, inputs, outputs)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	Normalize: func(ins Work) Normalized {
		out, outErr := createNormalizedWork(ins)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	Denormalize: func(ins interface{}) Work {
		if data, ok := ins.([]byte); ok {
			ptr := new(normalizedWork)
			jsErr := cdc.UnmarshalJSON(data, ptr)
			if jsErr != nil {
				panic(jsErr)
			}

			out, outErr := fromNormalizedToWork(ptr)
			if outErr != nil {
				panic(outErr)
			}

			return out
		}

		if normalized, ok := ins.(Normalized); ok {
			out, outErr := fromNormalizedToWork(normalized)
			if outErr != nil {
				panic(outErr)
			}

			return out
		}

		panic(errors.New("the given instance cannot be denormalized to a Work instance"))
	},
}
