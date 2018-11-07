package core

import (
	"math/rand"
	"net"
	"path/filepath"
	"testing"

	uuid "github.com/satori/go.uuid"
	ed25519 "github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/xmnservices/xmnsuite/blockchains/core/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/genesis"
	"github.com/xmnservices/xmnsuite/blockchains/core/wallet"
	"github.com/xmnservices/xmnsuite/crypto"
)

func TestSaveGenesis_thenRetrieveByID_Success(t *testing.T) {
	// variables:
	namespace := "xmn"
	name := "core"
	id := uuid.NewV4()
	rootPath := filepath.Join("./test_files")
	port := rand.Int()%9000 + 1000
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	nodePK := ed25519.GenPrivKey()
	ip := net.ParseIP("127.0.0.1")
	defer func() {
		//os.RemoveAll(rootPath)
	}()

	// spawn the blockchain:
	node, nodeErr := spawnBlockchain(namespace, name, &id, rootPath, port, nodePK, pk.PublicKey())
	if nodeErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", nodeErr.Error())
		return
	}

	node.Start()
	defer node.Stop()

	// get the client:
	client, clientErr := connectToBlockchain(ip, port)
	if clientErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", clientErr.Error())
		return
	}

	// create a genesis:
	genIns := genesis.CreateGenesisForTests()

	// create the representation:
	representation := genesis.SDKFunc.CreateRepresentation(genesis.CreateRepresentationParams{
		DepositRepresentation: deposit.SDKFunc.CreateRepresentation(deposit.CreateRepresentationParams{
			WalletRepresentation: wallet.SDKFunc.CreateRepresentation(),
		}),
	})

	// create the entity service:
	entityService := entity.SDKFunc.CreateSDKService(entity.CreateSDKServiceParams{
		PK:     pk,
		Client: client,
	})

	// create the entity repository:
	entityRepository := entity.SDKFunc.CreateSDKRepository(entity.CreateSDKRepositoryParams{
		PK:     pk,
		Client: client,
	})

	// save the genesis:
	saveErr := entityService.Save(genIns, representation)
	if saveErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveErr.Error())
		return
	}

	// retrieve the genesis:
	retGen, retGenErr := entityRepository.RetrieveByID(representation.MetaData(), genIns.ID())
	if retGenErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retGenErr.Error())
		return
	}

	// compare the wallet instances:
	genesis.CompareGenesisForTests(t, genIns, retGen.(genesis.Genesis))
}
