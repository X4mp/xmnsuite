package core

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/pledge"
	"github.com/xmnservices/xmnsuite/blockchains/core/token"
	"github.com/xmnservices/xmnsuite/blockchains/core/wallet"
	"github.com/xmnservices/xmnsuite/crypto"
)

func TestSaveGenesis_thenRetrieveByID_Success(t *testing.T) {
	// variables:
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	rootPath := filepath.Join("./test_files")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, _, _ := spawnBlockchainWithGenesisForTests(t, pk, rootPath)
	defer node.Stop()
}

func TestSaveGenesis_thenCRUD_Success(t *testing.T) {

	pldge := pledge.CreatePledgeForTests()

	// variables:
	testEntities := []struct {
		Ins            entity.Entity
		Representation entity.Representation
		Prepare        func(repository entity.Repository, service entity.Service)
		Teardown       func(repository entity.Repository, service entity.Service)
		Compare        func(t *testing.T, first entity.Entity, second entity.Entity)
	}{
		{
			Ins:            wallet.CreateWalletForTests(),
			Representation: wallet.SDKFunc.CreateRepresentation(),
			Compare: func(t *testing.T, first entity.Entity, second entity.Entity) {
				wallet.CompareWalletsForTests(t, first.(wallet.Wallet), second.(wallet.Wallet))
			},
		},
		{
			Ins:            token.CreateTokenForTests(),
			Representation: token.SDKFunc.CreateRepresentation(),
			Compare: func(t *testing.T, first entity.Entity, second entity.Entity) {
				token.CompareTokensForTests(t, first.(token.Token), second.(token.Token))
			},
		},
		{
			Ins:            pldge,
			Representation: pledge.SDKFunc.CreateRepresentation(),
			Prepare: func(repository entity.Repository, service entity.Service) {
				service.Save(pldge.From().From(), wallet.SDKFunc.CreateRepresentation())
				service.Save(pldge.From().Token(), token.SDKFunc.CreateRepresentation())
			},
			Compare: func(t *testing.T, first entity.Entity, second entity.Entity) {
				pledge.ComparePledgesForTests(t, first.(pledge.Pledge), second.(pledge.Pledge))
			},
		},
	}

	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	rootPath := filepath.Join("./test_files")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, service, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath)
	defer node.Stop()

	// test all instances:
	for _, oneTestEntity := range testEntities {
		// prepare:
		if oneTestEntity.Prepare != nil {
			oneTestEntity.Prepare(repository, service)
		}

		// execute:
		retIns := saveEntityThenRetrieveEntityByIDThenDeleteEntityByID(t, oneTestEntity.Ins, oneTestEntity.Representation, service, repository)
		oneTestEntity.Compare(t, oneTestEntity.Ins, retIns)

		// teardown:
		if oneTestEntity.Teardown != nil {
			oneTestEntity.Teardown(repository, service)
		}
	}
}
