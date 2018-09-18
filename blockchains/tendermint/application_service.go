package tendermint

import (
	"os"
	"path/filepath"

	config "github.com/tendermint/tendermint/config"
	log "github.com/tendermint/tendermint/libs/log"
	nm "github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
	applications "github.com/xmnservices/xmnsuite/blockchains/applications"
)

type applicationService struct {
}

func createApplicationService() ApplicationService {
	out := applicationService{}
	return &out
}

// Spawn spawns a new blockchain application
func (obj *applicationService) Spawn(rootDir string, blkChain Blockchain, apps applications.Applications) (applications.Node, error) {

	// retrieve the genesis block:
	gen := blkChain.GetGenesis()

	//create the abci application:
	abciApp, abciAppErr := createABCIApplication(apps)
	if abciAppErr != nil {
		return nil, abciAppErr
	}

	//create the config:
	dirPath := filepath.Join(rootDir, gen.GetPath().String())
	conf := config.DefaultConfig().SetRoot(dirPath)

	// create the node:
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	logger = log.NewFilter(logger, log.AllowError())
	pvFile := conf.PrivValidatorFile()
	pv := privval.LoadFilePV(pvFile)
	papp := proxy.NewLocalClientCreator(abciApp)
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
