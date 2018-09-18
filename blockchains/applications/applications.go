package applications

import (
	"errors"
	"fmt"
)

type applications struct {
	apps []Application
}

func createApplications(apps []Application) Applications {
	out := applications{
		apps: apps,
	}

	return &out
}

// RetrieveBlockIndex retrieves the highest block index
func (app *applications) RetrieveBlockIndex() int64 {
	highestBlockIndex := int64(0)
	for _, oneApp := range app.apps {
		blkHeight := oneApp.GetBlockIndex()
		if blkHeight > highestBlockIndex {
			highestBlockIndex = blkHeight
		}
	}

	return highestBlockIndex
}

// RetrieveByBlockIndex retrieves an application by its block index
func (app *applications) RetrieveByBlockIndex(blkIndex int64) (Application, error) {
	for _, oneApp := range app.apps {
		fromBlockIndex := oneApp.FromBlockIndex()
		toBlockIndex := oneApp.ToBlockIndex()

		if blkIndex >= fromBlockIndex {
			if toBlockIndex == -1 {
				return oneApp, nil
			}

			if blkIndex < toBlockIndex {
				return oneApp, nil
			}
		}
	}

	str := fmt.Sprintf("the block index (%d) has no matching application", blkIndex)
	return nil, errors.New(str)
}
