package core

import (
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/xmnservices/xmnsuite/blockchains/core/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/genesis"
	"github.com/xmnservices/xmnsuite/blockchains/core/pledge"
	"github.com/xmnservices/xmnsuite/blockchains/core/request"
	"github.com/xmnservices/xmnsuite/blockchains/core/token"
	"github.com/xmnservices/xmnsuite/blockchains/core/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/withdrawal"
	"github.com/xmnservices/xmnsuite/crypto"
)

func TestSaveGenesis_thenRetrieveByID_Success(t *testing.T) {
	// variables:
	genIns := genesis.CreateGenesisForTests()
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	rootPath := filepath.Join("./test_files")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, _, _, _ := spawnBlockchainWithGenesisForTests(t, pk, rootPath, genIns)
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

	genIns := genesis.CreateGenesisForTests()
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	rootPath := filepath.Join("./test_files")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, _, service, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, genIns)
	defer node.Stop()

	// test all instances:
	for _, oneTestEntity := range testEntities {
		// prepare:
		if oneTestEntity.Prepare != nil {
			oneTestEntity.Prepare(repository, service)
		}

		// execute:
		retIns := saveEntityThenRetrieveEntityByIDThenDeleteEntityByID(t, oneTestEntity.Ins, oneTestEntity.Representation, service, repository)
		if retIns == nil {
			return
		}

		// compare:
		oneTestEntity.Compare(t, oneTestEntity.Ins, retIns)

		// teardown:
		if oneTestEntity.Teardown != nil {
			oneTestEntity.Teardown(repository, service)
		}
	}
}

func TestSaveGenesis_savePledgeRequest_saveVotesOnRequest_Success(t *testing.T) {
	pledgeAmount := rand.Int() % 50
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})

	tok := token.SDKFunc.Create(token.CreateParams{
		Symbol:      "XMND",
		Name:        "XMN Dollar",
		Description: "This is the XMN dollar",
	})

	fromWallet := wallet.SDKFunc.Create(wallet.CreateParams{
		ConcensusNeeded: rand.Int() % 30,
	})

	toWallet := wallet.SDKFunc.Create(wallet.CreateParams{
		ConcensusNeeded: rand.Int() % 200,
	})

	fromUser := user.SDKFunc.Create(user.CreateParams{
		PubKey: pk.PublicKey(),
		Shares: 5,
		Wallet: fromWallet,
	})

	gen := genesis.SDKFunc.Create(genesis.CreateParams{
		GazPricePerKb:         2,
		MaxAmountOfValidators: 20,
		User: fromUser,
		Deposit: deposit.SDKFunc.Create(deposit.CreateParams{
			To:     fromWallet,
			Token:  tok,
			Amount: pledgeAmount + rand.Int()%200,
		}),
	})

	pldge := pledge.SDKFunc.Create(pledge.CreateParams{
		From: withdrawal.SDKFunc.Create(withdrawal.CreateParams{
			From: wallet.SDKFunc.Create(wallet.CreateParams{
				ConcensusNeeded: rand.Int() % 30,
			}),
			Token:  tok,
			Amount: pledgeAmount,
		}),
		To: toWallet,
	})

	pledgeRequest := request.SDKFunc.Create(request.CreateParams{
		FromUser:  fromUser,
		NewEntity: pldge,
	})

	rootPath := filepath.Join("./test_files")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, client, _, _ := spawnBlockchainWithGenesisForTests(t, pk, rootPath, gen)
	defer node.Stop()

	// create the request service:
	reqService := request.SDKFunc.CreateSDKService(request.CreateSDKServiceParams{
		PK:     pk,
		Client: client,
	})

	// save the request:
	saveErr := reqService.Save(pledgeRequest, pledge.SDKFunc.CreateRepresentation())
	if saveErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveErr.Error())
		return
	}
}
