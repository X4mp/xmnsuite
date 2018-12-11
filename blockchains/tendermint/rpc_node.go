package tendermint

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	nm "github.com/tendermint/tendermint/node"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	rpcclient "github.com/tendermint/tendermint/rpc/lib/client"
	applications "github.com/xmnservices/xmnsuite/blockchains/applications"
)

type rpcNode struct {
	rpcAddress string
	node       *nm.Node
	dbPath     string
}

func createRPCNode(rpcAddress string, node *nm.Node, dbPath string) applications.Node {
	out := rpcNode{
		rpcAddress: rpcAddress,
		node:       node,
		dbPath:     dbPath,
	}

	return &out
}

// GetAddress returns the address
func (app *rpcNode) GetAddress() string {
	return app.rpcAddress
}

// GetClient returns a client connected to the current node
func (app *rpcNode) GetClient() (applications.Client, error) {
	return createRPCClient(app.rpcAddress)
}

// Start starts the node
func (app *rpcNode) Start() error {
	startErr := app.node.Start()
	if startErr != nil {
		return startErr
	}

	//wait for RPC and GRPC:
	app.waitForRPC()
	return nil
}

// Stop stops the node
func (app *rpcNode) Stop() error {
	stopErr := app.node.Stop()
	if stopErr != nil {
		return stopErr
	}

	for {
		if !app.node.IsRunning() {
			// remove all the LOCK files:
			return app.removeLockFiles(app.dbPath)
		}

		fmt.Println("node still running, waiting...")
		time.Sleep(time.Second)
	}
}

func (app *rpcNode) waitForRPC() {
	client := rpcclient.NewJSONRPCClient(app.rpcAddress)
	ctypes.RegisterAmino(client.Codec())
	result := new(ctypes.ResultStatus)
	for {
		_, err := client.Call("status", map[string]interface{}{}, result)
		if err == nil {
			return
		}

		fmt.Println("error", err)
		time.Sleep(time.Millisecond)
	}
}

func (app *rpcNode) removeLockFiles(path string) error {
	files, filesErr := ioutil.ReadDir(path)
	if filesErr != nil {
		return filesErr
	}

	for _, oneFile := range files {
		name := oneFile.Name()
		subPath := filepath.Join(path, name)
		if oneFile.IsDir() {
			subDirErr := app.removeLockFiles(subPath)
			if subDirErr != nil {
				return subDirErr
			}

			continue
		}

		if name == "LOCK" {
			remErr := os.Remove(subPath)
			if remErr != nil {
				return remErr
			}
		}
	}

	return nil
}
