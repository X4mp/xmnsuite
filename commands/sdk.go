package commands

import (
	"net"
	"net/url"
	"strconv"

	uuid "github.com/satori/go.uuid"
	tcrypto "github.com/tendermint/tendermint/crypto"
	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/genesis"
	"github.com/xmnservices/xmnsuite/blockchains/core/meta"
	"github.com/xmnservices/xmnsuite/crypto"
)

// Command represents a cli command
type Command interface {
	Execute() (applications.Node, error)
}

// Constants represents the CLI constants
type Constants interface {
	Namespace() string
	Name() string
	ID() *uuid.UUID
	RoutePrefix() string
	RouterRoleKey() string
}

// Configs represents the CLI configs
type Configs interface {
	Constants() Constants
	Port() int
	NodePrivateKey() tcrypto.PrivKey
	BlockchainRootDirectory() string
	DatabaseFilePath() string
	Meta() meta.Meta
}

// GenesisConfigs represents the genesis configs
type GenesisConfigs interface {
	Configs() Configs
	GenesisTransaction() genesis.Genesis
	RootPrivateKey() crypto.PrivateKey
}

// StartConfigs represents the start configs
type StartConfigs interface {
	Configs() Configs
	HasPeers() bool
	Peers() []Node
}

// Node represents a blockchain node
type Node interface {
	IP() net.IP
	Port() int
	String() string
}

// CreateConstantsParams represents the create constants params
type CreateConstantsParams struct {
	Namespace     string
	Name          string
	ID            string
	RoutePrefix   string
	RouterRoleKey string
}

// CreateParams represents the create configs params
type CreateParams struct {
	Constants               CreateConstantsParams
	Port                    int
	NodePrivateKey          string
	BlockchainRootDirectory string
	DatabaseFilePath        string
	Meta                    meta.Meta
}

// CreateGenesisParams represents the create genesis params
type CreateGenesisParams struct {
	Configs            CreateParams
	RootPrivateKey     string
	GenesisTransaction genesis.Genesis
}

// CreateStartParams represents the create start params
type CreateStartParams struct {
	Configs CreateParams
	Peers   []string
}

// SDKFunc represents the cli SDK func
var SDKFunc = struct {
	CreateGenesis func(params CreateGenesisParams) Command
	CreateStart   func(params CreateStartParams) Command
}{
	CreateGenesis: func(params CreateGenesisParams) Command {
		// create the configs:
		conf := createConfigsFromParams(params.Configs)

		// create the root public key:
		rootPrivKey, rootPrivKeyErr := createRootPrivateKey(params.RootPrivateKey)
		if rootPrivKeyErr != nil {
			panic(rootPrivKeyErr)
		}

		// create the genesis configs:
		genConfigs, genConfigsErr := createGenesisConfigs(conf, rootPrivKey, params.GenesisTransaction)
		if genConfigsErr != nil {
			panic(genConfigsErr)
		}

		// create the command:
		cmd, cmdErr := createGenesisCommand(genConfigs)
		if cmdErr != nil {
			panic(cmdErr)
		}

		return cmd
	},
	CreateStart: func(params CreateStartParams) Command {
		// create the configs:
		conf := createConfigsFromParams(params.Configs)

		// create the nodes, if any:
		nodes := []Node{}
		if params.Peers != nil && len(params.Peers) > 0 {
			for _, onePeers := range params.Peers {
				ur, urErr := url.Parse(onePeers)
				if urErr != nil {
					panic(urErr)
				}

				port, portErr := strconv.Atoi(ur.Port())
				if portErr != nil {
					panic(portErr)
				}

				ip := net.ParseIP(ur.Host)
				nod, nodErr := createNode(ip, port)
				if nodErr != nil {
					panic(nodErr)
				}

				nodes = append(nodes, nod)
			}
		}

		// create the start configs:
		startConf, startConfErr := createStartConfigs(conf, nodes)
		if startConfErr != nil {
			panic(startConfErr)
		}

		out, outErr := createStartCommand(startConf)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
}
