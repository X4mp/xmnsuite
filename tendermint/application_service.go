package tendermint

import (
	"os"
	"path/filepath"

	applications "github.com/XMNBlockchain/datamint/applications"
	config "github.com/tendermint/tendermint/config"
	log "github.com/tendermint/tendermint/libs/log"
	nm "github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
)

type applicationService struct {
	rootDir  string
	blkChain Blockchain
	apps     applications.Applications
}

func createApplicationService(rootDir string, blkChain Blockchain, apps applications.Applications) (ApplicationService, error) {
	out := applicationService{
		rootDir:  rootDir,
		blkChain: blkChain,
		apps:     apps,
	}

	return &out, nil
}

// Spawn spawns a new blockchain application
func (obj *applicationService) Spawn() (applications.Node, error) {

	// retrieve the genesis block:
	gen := obj.blkChain.GetGenesis()

	//create the abci application:
	abciApp, abciAppErr := createABCIApplication(obj.apps)
	if abciAppErr != nil {
		return nil, abciAppErr
	}

	//create the config:
	dirPath := filepath.Join(obj.rootDir, gen.GetPath().String())
	conf := config.DefaultConfig().SetRoot(dirPath)

	// create the node:
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	logger = log.NewFilter(logger, log.AllowError())
	pvFile := conf.PrivValidatorFile()
	pv := privval.LoadFilePV(pvFile)
	papp := proxy.NewLocalClientCreator(abciApp)
	node, nodeErr := nm.NewNode(conf, pv, papp,
		nm.DefaultGenesisDocProviderFunc(conf),
		nm.DefaultDBProvider,
		nm.DefaultMetricsProvider(conf.Instrumentation),
		logger)

	if nodeErr != nil {
		return nil, nodeErr
	}

	// create the node:
	out := createRPCNode(conf.RPC.ListenAddress, node)
	return out, nil
}

// Connect connects to an external blockchain
func (obj *applicationService) Connect(ipAddress string) (applications.Client, error) {
	out, outErr := createRPCClient(ipAddress)
	if outErr != nil {
		return nil, outErr
	}

	return out, outErr
}
