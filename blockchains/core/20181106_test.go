package core

import (
	"math"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/xmnservices/xmnsuite/blockchains/core/meta"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/pledge"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/transfer"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/validator"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/work"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/genesis"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/keyname"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/vote"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/link"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/node"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/withdrawal"
	"github.com/xmnservices/xmnsuite/crypto"
)

func TestSaveGenesis_Success(t *testing.T) {
	// variables:
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	pubKey := pk.PublicKey()
	genIns := genesis.CreateGenesisWithPubKeyForTests(pubKey)
	rootPath := filepath.Join("./test_files_TestSaveGenesis_Success")
	routePrefix := "/some-route-prefix"
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, _, _, _ := spawnBlockchainWithGenesisForTests(t, pk, rootPath, routePrefix, genIns)
	defer node.Stop()
}

func TestSaveGenesis_thenRespawnBlockchain_Success(t *testing.T) {
	// variables:
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	pubKey := pk.PublicKey()
	genIns := genesis.CreateGenesisWithPubKeyForTests(pubKey)
	rootPath := filepath.Join("./test_files_TestSaveGenesis_Success")
	routePrefix := "/some-route-prefix"
	namespace := "xmn"
	name := "core"
	id := uuid.NewV4()
	port := rand.Int()%9000 + 1000
	nodePK := ed25519.GenPrivKey()
	met := meta.SDKFunc.Create(meta.CreateParams{})

	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, nodeErr := saveThenSpawnBlockchain(namespace, name, &id, nil, rootPath, routePrefix, port, nodePK, pk.PublicKey(), met)
	if nodeErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", nodeErr.Error())
		return
	}

	// start the node:
	startErr := node.Start()
	if startErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", startErr.Error())
		return
	}

	// get the client:
	client, clientErr := node.GetClient()
	if clientErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", clientErr.Error())
		return
	}

	// create the entity service:
	entityService := entity.SDKFunc.CreateSDKService(entity.CreateSDKServiceParams{
		PK:          pk,
		Client:      client,
		RoutePrefix: routePrefix,
	})

	// create the entity repository:
	entityRepository := entity.SDKFunc.CreateSDKRepository(entity.CreateSDKRepositoryParams{
		PK:          pk,
		Client:      client,
		RoutePrefix: routePrefix,
	})

	// create the representation:
	representation := genesis.SDKFunc.CreateRepresentation()

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

	// stop de node:
	stopErr := node.Stop()
	if stopErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", stopErr.Error())
		return
	}

	// spawn the blockchain again:
	secondNode, secondNodeErr := spawnBlockchain(namespace, name, &id, nil, rootPath, routePrefix, port, met)
	if secondNodeErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", secondNodeErr.Error())
		return
	}

	// start the node:
	secondStartErr := secondNode.Start()
	if secondStartErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", secondStartErr.Error())
		return
	}

	// stop the node:
	secondStopErr := secondNode.Stop()
	if secondStopErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", secondStopErr.Error())
		return
	}
}

func TestSaveGenesis_createSameGenesisInstance_returnsError(t *testing.T) {
	// variables:
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	pubKey := pk.PublicKey()
	genIns := genesis.CreateGenesisWithPubKeyForTests(pubKey)
	rootPath := filepath.Join("./test_files_TestSaveGenesis_createSameGenesisInstance_returnsError")
	routePrefix := "/some-route-prefix"
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, _, service, _ := spawnBlockchainWithGenesisForTests(t, pk, rootPath, routePrefix, genIns)
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

func TestSaveGenesis_createAccount_Success(t *testing.T) {
	// variables:
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	pubKey := pk.PublicKey()
	genIns := genesis.CreateGenesisWithPubKeyForTests(pubKey)

	newWalletPK := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	newWalletPubKey := newWalletPK.PublicKey()
	walletIns := wallet.CreateWalletWithPublicKeyForTests(newWalletPubKey)
	userIns := user.CreateUserWithWalletAndPublicKeyAndSharesForTests(walletIns, newWalletPubKey, walletIns.ConcensusNeeded())

	rootPath := filepath.Join("./test_files_TestSaveGenesis_createWallet_Success")
	routePrefix := "/some-route-prefix"
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, _, _, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, routePrefix, genIns)
	defer node.Stop()

	client, clientErr := node.GetClient()
	if clientErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", clientErr.Error())
		return
	}

	// create the service:
	accountService := account.SDKFunc.CreateSDKService(account.CreateSDKServiceParams{
		PK:          pk,
		Client:      client,
		RoutePrefix: routePrefix,
	})

	// save the account
	acc := saveAccountForTests(t, userIns, genIns, accountService, repository)
	if acc == nil {
		return
	}
}

func TestSaveGenesis_createAccount_withWorkMatrixTooSmall_returnsError(t *testing.T) {
	// variables:
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	pubKey := pk.PublicKey()
	genIns := genesis.CreateGenesisWithPubKeyForTests(pubKey)

	newWalletPK := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	newWalletPubKey := newWalletPK.PublicKey()
	walletIns := wallet.CreateWalletWithPublicKeyForTests(newWalletPubKey)
	userIns := user.CreateUserWithWalletAndPublicKeyAndSharesForTests(walletIns, newWalletPubKey, walletIns.ConcensusNeeded())

	ac := account.SDKFunc.Create(account.CreateAccountParams{
		User: userIns,
		Work: work.SDKFunc.Generate(work.GenerateParams{
			MatrixSize:   1,
			MatrixAmount: 1,
		}),
	})

	rootPath := filepath.Join("./test_files_TestSaveGenesis_createWallet_Success")
	routePrefix := "/some-route-prefix"
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, _, _, _ := spawnBlockchainWithGenesisForTests(t, pk, rootPath, routePrefix, genIns)
	defer node.Stop()

	client, clientErr := node.GetClient()
	if clientErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", clientErr.Error())
		return
	}

	// create the service:
	accountService := account.SDKFunc.CreateSDKService(account.CreateSDKServiceParams{
		PK:          pk,
		Client:      client,
		RoutePrefix: routePrefix,
	})

	// save the account:
	saveErr := accountService.Save(ac, genIns.GazPriceInMatrixWorkKb())
	if saveErr == nil {
		t.Errorf("the returned error was expected to be an error, nil returned")
		return
	}
}

func TestSaveGenesis_createWalletWithSameCreator_Success(t *testing.T) {
	// variables:
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	pubKey := pk.PublicKey()
	genIns := genesis.CreateGenesisWithPubKeyForTests(pubKey)
	walletIns := wallet.CreateWalletWithPublicKeyForTests(genIns.User().Wallet().Creator())
	userIns := user.CreateUserWithWalletAndPublicKeyAndSharesForTests(walletIns, walletIns.Creator(), walletIns.ConcensusNeeded())

	rootPath := filepath.Join("./test_files_TestSaveGenesis_createWalletWithSameCreator_Success")
	routePrefix := "/some-route-prefix"
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, _, _, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, routePrefix, genIns)
	defer node.Stop()

	client, clientErr := node.GetClient()
	if clientErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", clientErr.Error())
		return
	}

	// create the service:
	accountService := account.SDKFunc.CreateSDKService(account.CreateSDKServiceParams{
		PK:          pk,
		Client:      client,
		RoutePrefix: routePrefix,
	})

	// save the account
	acc := saveAccountForTests(t, userIns, genIns, accountService, repository)
	if acc == nil {
		return
	}

	// compare the wallets:
	wallet.CompareWalletsForTests(t, walletIns.(wallet.Wallet), acc.User().Wallet())
}

func TestSaveGenesis_addUserToWallet_addAnotherUserToWallerWithSamePublicKey_saveVotesWithEnoughSharesToPass_returnsError(t *testing.T) {
	// variables:
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	pubKey := pk.PublicKey()
	genIns := genesis.CreateGenesisWithPubKeyForTests(pubKey)

	userIns := user.CreateUserWithWalletAndPublicKeyAndSharesForTests(genIns.User().Wallet(), genIns.User().Wallet().Creator(), genIns.User().Wallet().ConcensusNeeded())
	rootPath := filepath.Join("./test_files_TestSaveGenesis_createWallet_addUserToWallet_addAnotherUserToWallerWithSamePublicKey_saveVotesWithEnoughSharesToPass_returnsError")
	routePrefix := "/some-route-prefix"
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, client, _, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, routePrefix, genIns)
	defer node.Stop()

	// retrieve the keyname:
	knameRepository := keyname.SDKFunc.CreateRepository(keyname.CreateRepositoryParams{
		EntityRepository: repository,
	})

	kname, knameErr := knameRepository.RetrieveByName(user.SDKFunc.CreateMetaData().Keyname())
	if knameErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", knameErr.Error())
	}

	// create the user in wallet request:
	req := request.SDKFunc.Create(request.CreateParams{
		FromUser:  genIns.User(),
		NewEntity: userIns,
		Reason:    "TEST",
		Keyname:   kname,
	})

	// create the request vote:
	reqVote := vote.SDKFunc.Create(vote.CreateParams{
		Request:    req,
		Voter:      genIns.User(),
		IsApproved: true,
	})

	// create the request service:
	requestService := request.SDKFunc.CreateSDKService(request.CreateSDKServiceParams{
		PK:          pk,
		Client:      client,
		RoutePrefix: routePrefix,
	})

	// create the vote service:
	voteService := vote.SDKFunc.CreateSDKService(vote.CreateSDKServiceParams{
		PK:              pk,
		Client:          client,
		CreateRouteFunc: createWalletVoteRouteFunc(routePrefix),
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
	rootPath := filepath.Join("./test_files_TestSaveGenesis_createNewUserOnWallet_Success")
	routePrefix := "/some-route-prefix"
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, client, _, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, routePrefix, genIns)
	defer node.Stop()

	// retrieve the keyname:
	knameRepository := keyname.SDKFunc.CreateRepository(keyname.CreateRepositoryParams{
		EntityRepository: repository,
	})

	kname, knameErr := knameRepository.RetrieveByName(user.SDKFunc.CreateMetaData().Keyname())
	if knameErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", knameErr.Error())
	}

	// create the user in wallet request:
	userInWalletRequest := request.SDKFunc.Create(request.CreateParams{
		FromUser:  genIns.User(),
		NewEntity: userIns,
		Reason:    "TEST",
		Keyname:   kname,
	})

	// create our genesis user vote:
	userInWalletRequestVote := vote.SDKFunc.Create(vote.CreateParams{
		Request:    userInWalletRequest,
		Voter:      genIns.User(),
		IsApproved: true,
	})

	// save the request then save votes:
	saveRequestThenSaveVotesForTests(t, routePrefix, client, pk, repository, user.SDKFunc.CreateRepresentation(), userInWalletRequest, []crypto.PrivateKey{pk}, []vote.Vote{
		userInWalletRequestVote,
	}, createWalletVoteRouteFunc(routePrefix))
}

func TestSaveGenesis_addUserToWallet_increaseTheNeededConcensus_voteUsingTwoUsers_Success(t *testing.T) {
	// variables:
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	pubKey := pk.PublicKey()
	genIns := genesis.CreateGenesisWithPubKeyForTests(pubKey)

	wal := genIns.User().Wallet()
	userPK := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	userIns := user.CreateUserWithWalletAndPublicKeyForTests(wal, userPK.PublicKey())
	rootPath := filepath.Join("./test_files_TestSaveGenesis_addUserToWallet_increaseTheNeededConcensus_voteUsingTwoUsers_Success")
	routePrefix := "/some-route-prefix"
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// create the representations:
	userRepresentation := user.SDKFunc.CreateRepresentation()
	walletRepresentation := wallet.SDKFunc.CreateRepresentation()

	// spawn bockchain with genesis instance:
	node, client, _, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, routePrefix, genIns)
	defer node.Stop()

	// retrieve the keyname:
	knameRepository := keyname.SDKFunc.CreateRepository(keyname.CreateRepositoryParams{
		EntityRepository: repository,
	})

	kname, knameErr := knameRepository.RetrieveByName(user.SDKFunc.CreateMetaData().Keyname())
	if knameErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", knameErr.Error())
	}

	// create the user in wallet request:
	userInWalletRequest := request.SDKFunc.Create(request.CreateParams{
		FromUser:  genIns.User(),
		NewEntity: userIns,
		Reason:    "TEST",
		Keyname:   kname,
	})

	// create our genesis user vote:
	userInWalletRequestVote := vote.SDKFunc.Create(vote.CreateParams{
		Request:    userInWalletRequest,
		Voter:      genIns.User(),
		IsApproved: true,
	})

	// save the request then save votes:
	saveRequestThenSaveVotesForTests(t, routePrefix, client, pk, repository, userRepresentation, userInWalletRequest, []crypto.PrivateKey{pk}, []vote.Vote{
		userInWalletRequestVote,
	}, createWalletVoteRouteFunc(routePrefix))

	// retrieve the keyname:
	walletKanme, walletKanmeErr := knameRepository.RetrieveByName(wallet.SDKFunc.CreateMetaData().Keyname())
	if walletKanmeErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", walletKanmeErr.Error())
	}

	// update the wallet to increase concensus:
	updateWalletRequest := request.SDKFunc.Create(request.CreateParams{
		FromUser: genIns.User(),
		NewEntity: wallet.SDKFunc.Create(wallet.CreateParams{
			ID:              wal.ID(),
			Creator:         wal.Creator(),
			ConcensusNeeded: genIns.User().Shares() + userIns.Shares(),
		}),
		Reason:  "TEST",
		Keyname: walletKanme,
	})

	// create our genesis user vote on the wallet update:
	updateWalletRequestVote := vote.SDKFunc.Create(vote.CreateParams{
		Request:    updateWalletRequest,
		Voter:      genIns.User(),
		IsApproved: true,
	})

	// save the new wallet request, then save vote:
	saveRequestThenSaveVotesForTests(t, routePrefix, client, pk, repository, walletRepresentation, updateWalletRequest, []crypto.PrivateKey{pk}, []vote.Vote{
		updateWalletRequestVote,
	}, createWalletVoteRouteFunc(routePrefix))

	// update the wallet to decrease concensus:
	updateAgainWalletRequest := request.SDKFunc.Create(request.CreateParams{
		FromUser: genIns.User(),
		NewEntity: wallet.SDKFunc.Create(wallet.CreateParams{
			ID:              wal.ID(),
			Creator:         wal.Creator(),
			ConcensusNeeded: genIns.User().Shares(),
		}),
		Reason:  "TEST",
		Keyname: walletKanme,
	})

	// create our genesis user vote on the wallet update:
	updateAgainWalletRequestVoteByGenUser := vote.SDKFunc.Create(vote.CreateParams{
		Request:    updateAgainWalletRequest,
		Voter:      genIns.User(),
		IsApproved: true,
	})

	// create our newly added user vote on the wallet update:
	updateAgainWalletRequestVoteByNewlyAddedUser := vote.SDKFunc.Create(vote.CreateParams{
		Request:    updateAgainWalletRequest,
		Voter:      userIns,
		IsApproved: true,
	})

	// save the new wallet request, then save vote:
	saveRequestThenSaveVotesForTests(t, routePrefix, client, pk, repository, walletRepresentation, updateAgainWalletRequest, []crypto.PrivateKey{pk, userPK}, []vote.Vote{
		updateAgainWalletRequestVoteByGenUser,
		updateAgainWalletRequestVoteByNewlyAddedUser,
	}, createWalletVoteRouteFunc(routePrefix))
}

func TestSaveGenesis_createNewWallet_createPledge_transferPledgeTokens_returnsError(t *testing.T) {
	// variables:
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	pubKey := pk.PublicKey()
	genIns := genesis.CreateGenesisWithPubKeyForTests(pubKey)

	walPK := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	walPubKey := walPK.PublicKey()
	walletIns := wallet.CreateWalletWithPublicKeyForTests(walPubKey)
	userIns := user.CreateUserWithWalletAndPublicKeyAndSharesForTests(walletIns, walPubKey, genIns.Deposit().Amount()*2)

	pldge := pledge.SDKFunc.Create(pledge.CreateParams{
		From: withdrawal.SDKFunc.Create(withdrawal.CreateParams{
			From:   genIns.Deposit().To(),
			Token:  genIns.Deposit().Token(),
			Amount: int(math.Floor(float64(genIns.Deposit().Amount() / 2))),
		}),
		To: walletIns,
	})

	// create the repreentations:
	transferRepresentation := transfer.SDKFunc.CreateRepresentation()

	rootPath := filepath.Join("./test_files_TestSaveGenesis_createNewWallet_createPledge_transferPledgeTokens_returnsError")
	routePrefix := "/some-route-prefix"
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, client, service, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, routePrefix, genIns)
	defer node.Stop()

	// save the pledge:
	savePledge(t, routePrefix, client, pk, service, repository, genIns, userIns, pldge)

	// transfer the pledge funds, returns error:
	trsf := transfer.SDKFunc.Create(transfer.CreateParams{
		Withdrawal: withdrawal.SDKFunc.Create(withdrawal.CreateParams{
			From:   walletIns,
			Token:  genIns.Deposit().Token(),
			Amount: pldge.From().Amount(),
		}),
		Deposit: deposit.SDKFunc.Create(deposit.CreateParams{
			To:     genIns.Deposit().To(),
			Token:  genIns.Deposit().Token(),
			Amount: pldge.From().Amount(),
		}),
	})

	// retrieve the keyname:
	knameRepository := keyname.SDKFunc.CreateRepository(keyname.CreateRepositoryParams{
		EntityRepository: repository,
	})

	kname, knameErr := knameRepository.RetrieveByName(transfer.SDKFunc.CreateMetaData().Keyname())
	if knameErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", knameErr.Error())
	}

	// create the user in wallet request:
	trsfRequest := request.SDKFunc.Create(request.CreateParams{
		FromUser:  userIns,
		NewEntity: trsf,
		Reason:    "TEST",
		Keyname:   kname,
	})

	// create our user vote:
	trsfRequestVote := vote.SDKFunc.Create(vote.CreateParams{
		Request:    trsfRequest,
		Voter:      userIns,
		IsApproved: true,
	})

	// create the request service:
	requestService := request.SDKFunc.CreateSDKService(request.CreateSDKServiceParams{
		PK:          walPK,
		Client:      client,
		RoutePrefix: routePrefix,
	})

	// save the request:
	saveRequestErr := requestService.Save(trsfRequest, transferRepresentation)
	if saveRequestErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveRequestErr.Error())
		return
	}

	// create the vote service:
	voteService := vote.SDKFunc.CreateSDKService(vote.CreateSDKServiceParams{
		PK:              walPK,
		Client:          client,
		CreateRouteFunc: createWalletVoteRouteFunc(routePrefix),
	})

	// save the vote, it should returns an error:
	savedVoteErr := voteService.Save(trsfRequestVote, transferRepresentation)
	if savedVoteErr == nil {
		t.Errorf("the returned error was expected to be valid, nil returned")
	}
}

func TestSaveGenesis_createNewWallet_createValidator_Success(t *testing.T) {
	// variables:
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	pubKey := pk.PublicKey()
	genIns := genesis.CreateGenesisWithPubKeyForTests(pubKey)

	walPK := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	walPubKey := walPK.PublicKey()
	walletIns := wallet.CreateWalletWithPublicKeyForTests(walPubKey)
	userIns := user.CreateUserWithWalletForTests(walletIns)

	vldator := validator.SDKFunc.Create(validator.CreateParams{
		PubKey: ed25519.GenPrivKey().PubKey(),
		Pledge: pledge.SDKFunc.Create(pledge.CreateParams{
			From: withdrawal.SDKFunc.Create(withdrawal.CreateParams{
				From:   genIns.Deposit().To(),
				Token:  genIns.Deposit().Token(),
				Amount: int(math.Floor(float64(genIns.Deposit().Amount() / 2))),
			}),
			To: walletIns,
		}),
	})

	// create the representations:
	validatorRepresentation := validator.SDKFunc.CreateRepresentation()

	rootPath := filepath.Join("./test_files_TestSaveGenesis_createNewWallet_createValidator_Success")
	routePrefix := "/some-route-prefix"
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, client, _, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, routePrefix, genIns)
	defer node.Stop()

	// create the service:
	accountService := account.SDKFunc.CreateSDKService(account.CreateSDKServiceParams{
		PK:          pk,
		Client:      client,
		RoutePrefix: routePrefix,
	})

	// save the account
	acc := saveAccountForTests(t, userIns, genIns, accountService, repository)
	if acc == nil {
		return
	}

	// compare the wallets:
	wallet.CompareWalletsForTests(t, walletIns.(wallet.Wallet), acc.User().Wallet())

	// retrieve the keyname:
	knameRepository := keyname.SDKFunc.CreateRepository(keyname.CreateRepositoryParams{
		EntityRepository: repository,
	})

	kname, knameErr := knameRepository.RetrieveByName(validator.SDKFunc.CreateMetaData().Keyname())
	if knameErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", knameErr.Error())
	}

	// create the user in validator request:
	validatorRequest := request.SDKFunc.Create(request.CreateParams{
		FromUser:  genIns.User(),
		NewEntity: vldator,
		Reason:    "TEST",
		Keyname:   kname,
	})

	// create our user vote:
	validatorRequestVote := vote.SDKFunc.Create(vote.CreateParams{
		Request:    validatorRequest,
		Voter:      genIns.User(),
		IsApproved: true,
	})

	// save the new wallet request, then save vote:
	saveRequestThenSaveVotesForTests(t, routePrefix, client, pk, repository, validatorRepresentation, validatorRequest, []crypto.PrivateKey{pk}, []vote.Vote{
		validatorRequestVote,
	}, createWalletVoteRouteFunc(routePrefix))
}

func TestSaveGenesis_createNewWallet_createTransfer_Success(t *testing.T) {
	// variables:
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	pubKey := pk.PublicKey()
	genIns := genesis.CreateGenesisWithPubKeyForTests(pubKey)

	walPK := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	walPubKey := walPK.PublicKey()
	walletIns := wallet.CreateWalletWithPublicKeyForTests(walPubKey)
	userIns := user.CreateUserWithWalletForTests(walletIns)

	trsf := transfer.SDKFunc.Create(transfer.CreateParams{
		Withdrawal: withdrawal.SDKFunc.Create(withdrawal.CreateParams{
			From:   genIns.Deposit().To(),
			Token:  genIns.Deposit().Token(),
			Amount: int(math.Floor(float64(genIns.Deposit().Amount() / 2))),
		}),
		Deposit: deposit.SDKFunc.Create(deposit.CreateParams{
			To:     walletIns,
			Token:  genIns.Deposit().Token(),
			Amount: int(math.Floor(float64(genIns.Deposit().Amount() / 2))),
		}),
	})

	// create the representations:
	transferRepresentation := transfer.SDKFunc.CreateRepresentation()

	rootPath := filepath.Join("./test_files_TestSaveGenesis_createNewWallet_createTransfer_Success")
	routePrefix := "/some-route-prefix"
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, client, _, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, routePrefix, genIns)
	defer node.Stop()

	// create the service:
	accountService := account.SDKFunc.CreateSDKService(account.CreateSDKServiceParams{
		PK:          pk,
		Client:      client,
		RoutePrefix: routePrefix,
	})

	// save the account
	acc := saveAccountForTests(t, userIns, genIns, accountService, repository)
	if acc == nil {
		return
	}

	// compare the wallets:
	wallet.CompareWalletsForTests(t, walletIns.(wallet.Wallet), acc.User().Wallet())

	// retrieve the keyname:
	knameRepository := keyname.SDKFunc.CreateRepository(keyname.CreateRepositoryParams{
		EntityRepository: repository,
	})

	kname, knameErr := knameRepository.RetrieveByName(transfer.SDKFunc.CreateMetaData().Keyname())
	if knameErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", knameErr.Error())
	}

	// create the user in wallet request:
	trsfRequest := request.SDKFunc.Create(request.CreateParams{
		FromUser:  genIns.User(),
		NewEntity: trsf,
		Reason:    "TEST",
		Keyname:   kname,
	})

	// create our user vote:
	trsfRequestVote := vote.SDKFunc.Create(vote.CreateParams{
		Request:    trsfRequest,
		Voter:      genIns.User(),
		IsApproved: true,
	})

	// save the new wallet request, then save vote:
	saveRequestThenSaveVotesForTests(t, routePrefix, client, pk, repository, transferRepresentation, trsfRequest, []crypto.PrivateKey{pk}, []vote.Vote{
		trsfRequestVote,
	}, createWalletVoteRouteFunc(routePrefix))
}

func TestSaveGenesis_CreateLink_voteOnLink_Success(t *testing.T) {
	// variables:
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	pubKey := pk.PublicKey()
	genIns := genesis.CreateGenesisWithPubKeyForTests(pubKey)

	lnk := link.SDKFunc.Create(link.CreateParams{
		Title:       "Projects",
		Description: "The XMN projects belongs on that blockchain",
	})

	rootPath := filepath.Join("./test_files_TestSaveGenesis_CreateLink_voteOnLink_Success")
	routePrefix := "/some-route-prefix"
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, client, service, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, routePrefix, genIns)
	defer node.Stop()

	// save the link:
	saveLink(t, routePrefix, client, pk, service, repository, genIns.User(), lnk)
}

func TestSaveGenesis_CreateLink_voteOnLink_CreateNode_voteOnNode_Success(t *testing.T) {
	// variables:
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	pubKey := pk.PublicKey()
	genIns := genesis.CreateGenesisWithPubKeyForTests(pubKey)

	lnk := link.SDKFunc.Create(link.CreateParams{
		Title:       "Projects",
		Description: "The XMN projects belongs on that blockchain",
	})

	nod := node.SDKFunc.Create(node.CreateParams{
		Power: rand.Int() % 10,
		IP:    net.ParseIP("127.0.0.1"),
		Port:  123124,
		Link:  lnk,
	})

	rootPath := filepath.Join("./test_files_TestSaveGenesis_CreateLink_voteOnLink_CreateNode_voteOnNode_Success")
	routePrefix := "/some-route-prefix"
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, client, service, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, routePrefix, genIns)
	defer node.Stop()

	// save the link:
	saveNode(t, routePrefix, client, pk, service, repository, genIns.User(), lnk, nod)
}
