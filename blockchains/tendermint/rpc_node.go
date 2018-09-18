package tendermint

import (
	"fmt"
	"time"

	nm "github.com/tendermint/tendermint/node"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	rpcclient "github.com/tendermint/tendermint/rpc/lib/client"
	applications "github.com/xmnservices/xmnsuite/blockchains/applications"
)

type rpcNode struct {
	rpcAddress string
	node       *nm.Node
}

func createRPCNode(rpcAddress string, node *nm.Node) applications.Node {
	out := rpcNode{
		rpcAddress: rpcAddress,
		node:       node,
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
	err := app.node.Stop()
	return err
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
