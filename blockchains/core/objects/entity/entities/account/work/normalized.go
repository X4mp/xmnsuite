package work

import (
	"strconv"
	"time"

	"github.com/oelmekki/matrix"
)

type normalizedWork struct {
	Timestamp time.Time    `json:"timestamp"`
	Input     [][][]string `json:"input"`
	Output    [][][]string `json:"output"`
}

func createNormalizedWork(ins Work) (*normalizedWork, error) {
	input := fromFloatToStringMatrix(ins.Input())
	output := fromFloatToStringMatrix(ins.Output())
	out := normalizedWork{
		Timestamp: ins.Timestamp(),
		Input:     input,
		Output:    output,
	}

	return &out, nil
}

func fromFloatToStringMatrix(in []matrix.Matrix) [][][]string {
	out := [][][]string{}
	for index, oneMatrix := range in {
		out = append(out, [][]string{})
		amountRows := oneMatrix.Rows()
		amountCols := oneMatrix.Cols()
		for rowsIndex := 0; rowsIndex < amountRows; rowsIndex++ {
			out[index] = append(out[index], []string{})
			for colsIndex := 0; colsIndex < amountCols; colsIndex++ {
				element := oneMatrix.At(rowsIndex, colsIndex)
				out[index][rowsIndex] = append(out[index][rowsIndex], strconv.FormatFloat(element, 'f', 0, 64))
			}
		}
	}

	return out
}
