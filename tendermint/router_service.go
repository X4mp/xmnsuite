package tendermint

import (
	"os"
	"path/filepath"
	"strings"

	datastore "github.com/XMNBlockchain/datamint/datastore"
	keys "github.com/XMNBlockchain/datamint/keys"
	router "github.com/XMNBlockchain/datamint/router"
	config "github.com/tendermint/tendermint/config"
	log "github.com/tendermint/tendermint/libs/log"
	nm "github.com/tendermint/tendermint/node"
	privval "github.com/tendermint/tendermint/privval"
	proxy "github.com/tendermint/tendermint/proxy"
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

	//logger:
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	logger = log.NewFilter(logger, log.AllowError())

	//pv file:
	pvFile := conf.PrivValidatorFile()

	pv := privval.LoadFilePV(pvFile)

	//local client:
	papp := proxy.NewLocalClientCreator(middleApp)

	//create the node:
	node, nodeErr := nm.NewNode(
		conf,
		pv,
		papp,
		nm.DefaultGenesisDocProviderFunc(conf),
		nm.DefaultDBProvider,
		nm.DefaultMetricsProvider(conf.Instrumentation),
		logger,
	)

	if nodeErr != nil {
		return nil, nodeErr
	}

	//start the node:
	nodeStartErr := node.Start()
	if nodeStartErr != nil {
		return nil, nodeStartErr
	}

	//create the client:
	client, clientErr := obj.Connect(conf.RPC.ListenAddress)
	if clientErr != nil {
		return nil, clientErr
	}

	//returns the client:
	return client, nil
}

// Connect connects to an external blockchain
func (obj *routerService) Connect(ipAddress string) (router.Router, error) {
	out, outErr := createGRPCRouter(ipAddress)
	return out, outErr
}
