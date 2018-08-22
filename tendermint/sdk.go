package tendermint

import (
	"time"

	router "github.com/XMNBlockchain/datamint/router"
	uuid "github.com/satori/go.uuid"
	crypto "github.com/tendermint/tendermint/crypto"
)

/*
 * Blockchain
 */

// Validator represents a validator
type Validator interface {
	GetName() string
	GetPower() int
	GetPubKey() crypto.PubKey
}

// Path represents a blockchain relative path
type Path interface {
	GetNamespace() string
	GetName() string
	GetID() *uuid.UUID
	String() string
}

// Genesis represents the genesis block
type Genesis interface {
	GetHead() []byte
	GetPath() Path
	GetValidators() []Validator
	CreatedOn() time.Time
}

// PrivateValidator represents a private validator
type PrivateValidator interface {
	GetAddress() string
	GetPubKey() crypto.PubKey
	GetPrivKey() crypto.PrivKey
	GetLastHeight() int64
	GetLastRound() int64
	GetLastStep() int8
}

// Blockchain represents the blockchain initial phase
type Blockchain interface {
	GetGenesis() Genesis
	GetPK() crypto.PrivKey
	GetPV() PrivateValidator
}

// BlockchainService represents the blockchain service
type BlockchainService interface {
	Retrieve(path Path) (Blockchain, error)
	Save(blkChain Blockchain) error
	Delete(path Path) error
}

/*
 * Application
 */

// RouterService represents an application service
type RouterService interface {
	Spawn() (router.Router, error)
	Connect(ipAddress string) (router.Router, error)
}

/*
 * Params
 */

// CreatePathParams represents the params of the CreatePath SDK func
type CreatePathParams struct {
	Namespace string
	Name      string
	ID        *uuid.UUID
}

// CreateBlockchainParams represents the params of the CreateBlockchain SDK func
type CreateBlockchainParams struct {
	Namespace string
	Name      string
	ID        *uuid.UUID
	PrivKey   crypto.PrivKey
}

// CreateBlockchainServiceParams represents the params of the CreateBlockchainService SDK func
type CreateBlockchainServiceParams struct {
	RootDirPath string
}

// CreateRouterServiceParams represents the params of the CreateRouterService SDK func
type CreateRouterServiceParams struct {
	RootDir  string
	BlkChain Blockchain
	Router   router.Router
}

// SDKFunc represents the tendermint interval blockchains SDK functions
var SDKFunc = struct {
	CreatePath              func(params CreatePathParams) Path
	CreateBlockchain        func(params CreateBlockchainParams) Blockchain
	CreateBlockchainService func(params CreateBlockchainServiceParams) BlockchainService
	CreateRouterService     func(params CreateRouterServiceParams) RouterService
}{
	CreatePath: func(params CreatePathParams) Path {
		return createPath(params.Namespace, params.Name, params.ID)
	},
	CreateBlockchain: func(params CreateBlockchainParams) Blockchain {
		if params.PrivKey != nil {
			out, outErr := generateBlockchainWithPrivateKey(params.Namespace, params.Name, params.ID, params.PrivKey)
			if outErr != nil {
				panic(outErr)
			}

			return out
		}

		out, outErr := generateBlockchain(params.Namespace, params.Name, params.ID)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	CreateBlockchainService: func(params CreateBlockchainServiceParams) BlockchainService {
		return createBlockchainService(params.RootDirPath)
	},
	CreateRouterService: func(params CreateRouterServiceParams) RouterService {
		return createRouterService(params.RootDir, params.BlkChain, params.Router)
	},
}
