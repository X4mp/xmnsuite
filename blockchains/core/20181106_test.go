package core

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/genesis"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/request"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/request/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/request/vote"
	"github.com/xmnservices/xmnsuite/crypto"
)

func TestSaveGenesis_Success(t *testing.T) {
	// variables:
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	pubKey := pk.PublicKey()
	genIns := genesis.CreateGenesisWithPubKeyForTests(pubKey)
	rootPath := filepath.Join("./test_files")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, _, _, _ := spawnBlockchainWithGenesisForTests(t, pk, rootPath, genIns)
	defer node.Stop()
}

func TestSaveGenesis_createSameGenesisInstance_returnsError(t *testing.T) {
	// variables:
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	pubKey := pk.PublicKey()
	genIns := genesis.CreateGenesisWithPubKeyForTests(pubKey)
	rootPath := filepath.Join("./test_files")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, _, service, _ := spawnBlockchainWithGenesisForTests(t, pk, rootPath, genIns)
	defer node.Stop()

	// create the representation:
	representation := genesis.SDKFunc.CreateRepresentation()

	// save the genesis:
	saveErr := service.Save(genIns, representation)
	if saveErr == nil {
		t.Errorf("the returned error was expected to be an error, nil returned")
		return
	}
}

func TestSaveGenesis_createWallet_Success(t *testing.T) {
	// variables:
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	pubKey := pk.PublicKey()
	genIns := genesis.CreateGenesisWithPubKeyForTests(pubKey)

	newWalletPK := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	newWalletPubKey := newWalletPK.PublicKey()
	walletIns := wallet.CreateWalletWithPublicKeyForTests(newWalletPubKey)
	rootPath := filepath.Join("./test_files")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, _, service, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, genIns)
	defer node.Stop()

	// save the new wallet:
	savedWallet := saveEntityThenRetrieveEntityByIDThenDeleteEntityByID(t, walletIns, wallet.SDKFunc.CreateRepresentation(), service, repository)

	// compare the wallets:
	wallet.CompareWalletsForTests(t, walletIns.(wallet.Wallet), savedWallet.(wallet.Wallet))
}

func TestSaveGenesis_createWalletWithSameCreator_Success(t *testing.T) {
	// variables:
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	pubKey := pk.PublicKey()
	genIns := genesis.CreateGenesisWithPubKeyForTests(pubKey)

	walletIns := wallet.CreateWalletWithPublicKeyForTests(genIns.User().Wallet().Creator())
	rootPath := filepath.Join("./test_files")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, _, service, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, genIns)
	defer node.Stop()

	// save the new wallet:
	savedWallet := saveEntityThenRetrieveEntityByIDThenDeleteEntityByID(t, walletIns, wallet.SDKFunc.CreateRepresentation(), service, repository)

	// compare the wallets:
	wallet.CompareWalletsForTests(t, walletIns.(wallet.Wallet), savedWallet.(wallet.Wallet))
}

func TestSaveGenesis_createWalletAlreadyInGenesis_returnsError(t *testing.T) {
	// variables:
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	pubKey := pk.PublicKey()
	genIns := genesis.CreateGenesisWithPubKeyForTests(pubKey)
	rootPath := filepath.Join("./test_files")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, _, service, _ := spawnBlockchainWithGenesisForTests(t, pk, rootPath, genIns)
	defer node.Stop()

	// save the genesis wallet again:
	saveErr := service.Save(genIns.User().Wallet(), wallet.SDKFunc.CreateRepresentation())
	if saveErr == nil {
		t.Errorf("the returned error was expected to be an error, nil returned")
		return
	}
}

func TestSaveGenesis_createWallet_addUserToWallet_addAnotherUserToWallerWithSamePublicKey_saveVotesWithEnoughSharesToPass_returnsError(t *testing.T) {
	// variables:
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	pubKey := pk.PublicKey()
	genIns := genesis.CreateGenesisWithPubKeyForTests(pubKey)

	userIns := user.CreateUserWithWalletAndPublicKeyAndSharesForTests(genIns.User().Wallet(), genIns.User().Wallet().Creator(), genIns.User().Wallet().ConcensusNeeded())
	rootPath := filepath.Join("./test_files")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// create the user in wallet request:
	req := request.SDKFunc.Create(request.CreateParams{
		FromUser:  genIns.User(),
		NewEntity: userIns,
	})

	// create the request vote:
	reqVote := vote.SDKFunc.Create(vote.CreateParams{
		Request:    req,
		Voter:      genIns.User(),
		IsApproved: true,
	})

	// spawn bockchain with genesis instance:
	node, client, _, _ := spawnBlockchainWithGenesisForTests(t, pk, rootPath, genIns)
	defer node.Stop()

	// create the request service:
	requestService := request.SDKFunc.CreateSDKService(request.CreateSDKServiceParams{
		PK:     pk,
		Client: client,
	})

	// create the vote service:
	voteService := vote.SDKFunc.CreateSDKService(vote.CreateSDKServiceParams{
		PK:     pk,
		Client: client,
	})

	// save the request, returns an error due to the duplicate pubKey on user, of same wallet:
	saveRequestErr := requestService.Save(req, user.SDKFunc.CreateRepresentation())
	if saveRequestErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveRequestErr.Error())
		return
	}

	// save the vote, returns an error:
	savedVoteErr := voteService.Save(reqVote, user.SDKFunc.CreateRepresentation())
	if savedVoteErr == nil {
		t.Errorf("the returned error was expected to be valid, nil returned")
		return
	}
}

func TestSaveGenesis_createNewUserOnWallet_Success(t *testing.T) {
	// variables:
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	pubKey := pk.PublicKey()
	genIns := genesis.CreateGenesisWithPubKeyForTests(pubKey)

	userIns := user.CreateUserWithWalletForTests(genIns.User().Wallet())
	rootPath := filepath.Join("./test_files")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// create the user in wallet request:
	userInWalletRequest := request.SDKFunc.Create(request.CreateParams{
		FromUser:  genIns.User(),
		NewEntity: userIns,
	})

	// create our genesis user vote:
	userInWalletRequestVote := vote.SDKFunc.Create(vote.CreateParams{
		Request:    userInWalletRequest,
		Voter:      genIns.User(),
		IsApproved: true,
	})

	// save the new wallet request, then save vote:
	node, _, _, _, _, _ := spawnBlockchainWithGenesisThenSaveRequestThenSaveVotesForTests(t, pk, rootPath, genIns, user.SDKFunc.CreateRepresentation(), user.SDKFunc.CreateRepresentation(), userInWalletRequest, []vote.Vote{
		userInWalletRequestVote,
	})

	defer node.Stop()
}
