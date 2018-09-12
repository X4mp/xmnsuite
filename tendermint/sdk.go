package tendermint

import (
	"time"

	uuid "github.com/satori/go.uuid"
	crypto "github.com/tendermint/tendermint/crypto"
	applications "github.com/xmnservices/xmnsuite/applications"
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

// ApplicationService represents an application service
type ApplicationService interface {
	Spawn(rootDir string, blkChain Blockchain, apps applications.Applications) (applications.Node, error)
	Connect(ipAddress string) (applications.Client, error)
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

// SDKFunc represents the tendermint interval blockchains SDK functions
var SDKFunc = struct {
	CreatePath               func(params CreatePathParams) Path
	CreateBlockchain         func(params CreateBlockchainParams) Blockchain
	CreateBlockchainService  func(params CreateBlockchainServiceParams) BlockchainService
	CreateApplicationService func() ApplicationService
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
	CreateApplicationService: func() ApplicationService {
		serv := createApplicationService()
		return serv
	},
}
