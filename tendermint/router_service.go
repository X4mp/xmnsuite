package tendermint

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	datastore "github.com/XMNBlockchain/datamint/datastore"
	keys "github.com/XMNBlockchain/datamint/keys"
	router "github.com/XMNBlockchain/datamint/router"
	config "github.com/tendermint/tendermint/config"
	log "github.com/tendermint/tendermint/libs/log"
	nm "github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	core_grpc "github.com/tendermint/tendermint/rpc/grpc"
	rpcclient "github.com/tendermint/tendermint/rpc/lib/client"
)

type routerService struct {
	rootDir  string
	blkChain Blockchain
	rter     router.Router
}

func createRouterService(rootDir string, blkChain Blockchain, rter router.Router) RouterService {
	out := routerService{
		rootDir:  rootDir,
		blkChain: blkChain,
		rter:     rter,
	}

	return &out
}

// Spawn spawns a new blockchain application
func (obj *routerService) Spawn() (router.Router, error) {

	//create the datastore and keys:
	k := keys.SDKFunc.Create()
	store := datastore.SDKFunc.Create()

	//bind the custom application to the middle application:
	gen := obj.blkChain.GetGenesis()
	keynamePrefix := strings.Replace(gen.GetPath().String(), string(filepath.Separator), "-", -1)
	middleApp, middleAppErr := createABCIApplication("stateKey", keynamePrefix, []byte(gen.GetHead()), k, store, obj.rter)
	if middleAppErr != nil {
		return nil, middleAppErr
	}

	//create the config:
	dirPath := filepath.Join(obj.rootDir, gen.GetPath().String())
	conf := config.DefaultConfig().SetRoot(dirPath)

	//add the GRPC listen address:
	conf.RPC.GRPCListenAddress = "tcp://0.0.0.0:2350"

	// create the node:
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	logger = log.NewFilter(logger, log.AllowError())
	pvFile := conf.PrivValidatorFile()
	pv := privval.LoadFilePV(pvFile)
	papp := proxy.NewLocalClientCreator(middleApp)
	node, nodeErr := nm.NewNode(conf, pv, papp,
		nm.DefaultGenesisDocProviderFunc(conf),
		nm.DefaultDBProvider,
		nm.DefaultMetricsProvider(conf.Instrumentation),
		logger)

	if nodeErr != nil {
		return nil, nodeErr
	}

	// start the node:
	startErr := node.Start()
	if startErr != nil {
		return nil, startErr
	}

	//conf.RPC.GRPCListenAddress = fmt.Sprintf("%s%d", conf.RPC.ListenAddress, 1)

	//wait for RPC and GRPC:
	waitForRPC(conf.RPC.ListenAddress)
	waitForGRPC(conf.RPC.GRPCListenAddress)

	client, clientErr := obj.Connect(conf.RPC.ListenAddress)
	if clientErr != nil {
		return nil, clientErr
	}

	return client, nil
}

// Connect connects to an external blockchain
func (obj *routerService) Connect(ipAddress string) (router.Router, error) {
	out, outErr := createRPCRouter(ipAddress)
	if outErr != nil {
		return nil, outErr
	}

	return out, outErr
}

func waitForRPC(laddr string) {
	client := rpcclient.NewJSONRPCClient(laddr)
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

func waitForGRPC(grpcAddr string) {
	client := core_grpc.StartGRPCClient(grpcAddr)
	for {
		_, err := client.Ping(context.Background(), &core_grpc.RequestPing{})
		if err == nil {
			return
		}

		fmt.Println("error", err)
		time.Sleep(time.Millisecond)
	}
}
