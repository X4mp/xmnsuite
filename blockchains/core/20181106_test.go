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
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/genesis"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/pledge"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project/feature"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project/milestone"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project/milestone/task"
	completed_task "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project/milestone/task/completed"
	pledge_task "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project/milestone/task/pledge"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/transfer"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/validator"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request"
	active_request "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active/vote"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/keyname"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/category"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/link"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/node"
	community_project "github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/project"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/withdrawal"
	"github.com/xmnservices/xmnsuite/crypto"
)

func TestSaveGenesis_Success(t *testing.T) {
	// variables:
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	pubKey := pk.PublicKey()
	genIns := genesis.CreateGenesisWithPubKeyForTests(pubKey)
	rootPath := filepath.Join("./test_files_TestSaveGenesis_Success")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, _, _, _ := spawnBlockchainWithGenesisForTests(t, pk, rootPath, genIns)
	defer node.Stop()
}

func TestSaveGenesis_thenRespawnBlockchain_Success(t *testing.T) {
	// variables:
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	pubKey := pk.PublicKey()
	genIns := genesis.CreateGenesisWithPubKeyForTests(pubKey)
	rootPath := filepath.Join("./test_files_TestSaveGenesis_Success")
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
	node, nodeErr := saveThenSpawnBlockchain(namespace, name, &id, nil, rootPath, port, nodePK, pk.PublicKey(), met)
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
		PK:     pk,
		Client: client,
	})

	// create the entity repository:
	entityRepository := entity.SDKFunc.CreateSDKRepository(entity.CreateSDKRepositoryParams{
		PK:     pk,
		Client: client,
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
	secondNode, secondNodeErr := spawnBlockchain(namespace, name, &id, nil, rootPath, port, met)
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

func TestSaveGenesis_createUser_Success(t *testing.T) {
	// variables:
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	pubKey := pk.PublicKey()
	genIns := genesis.CreateGenesisWithPubKeyForTests(pubKey)

	newWalletPK := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	newWalletPubKey := newWalletPK.PublicKey()
	walletIns := wallet.CreateWalletWithPublicKeyForTests(newWalletPubKey)
	userIns := user.CreateUserWithWalletAndPublicKeyAndSharesForTests(walletIns, newWalletPubKey, walletIns.ConcensusNeeded())

	rootPath := filepath.Join("./test_files_TestSaveGenesis_createWallet_Success")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, client, service, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, genIns)
	defer node.Stop()

	// create the user with the new wallet:
	saveUserWithNewWallet(t, client, pk, service, repository, genIns.User(), userIns)
}

func TestSaveGenesis_addUserToWallet_addAnotherUserToWallerWithSamePublicKey_saveVotesWithEnoughSharesToPass_returnsError(t *testing.T) {
	// variables:
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	pubKey := pk.PublicKey()
	genIns := genesis.CreateGenesisWithPubKeyForTests(pubKey)

	userIns := user.CreateUserWithWalletAndPublicKeyAndSharesForTests(genIns.User().Wallet(), genIns.User().Wallet().Creator(), genIns.User().Wallet().ConcensusNeeded())
	rootPath := filepath.Join("./test_files_TestSaveGenesis_createWallet_addUserToWallet_addAnotherUserToWallerWithSamePublicKey_saveVotesWithEnoughSharesToPass_returnsError")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, client, _, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, genIns)
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
		FromUser:   genIns.User(),
		SaveEntity: userIns,
		Reason:     "TEST",
		Keyname:    kname,
	})

	// create the request service:
	requestService := request.SDKFunc.CreateSDKService(request.CreateSDKServiceParams{
		PK:     pk,
		Client: client,
	})

	requestRepository := active_request.SDKFunc.CreateRepository(active_request.CreateRepositoryParams{
		EntityRepository: repository,
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

	// retrieve the request:
	retReq, retReqErr := requestRepository.RetrieveByRequest(req)
	if retReqErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retReqErr.Error())
		return
	}

	// create the request vote:
	reqVote := vote.SDKFunc.Create(vote.CreateParams{
		Request:    retReq,
		Voter:      genIns.User(),
		IsApproved: true,
	})

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
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, client, _, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, genIns)
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
		FromUser:   genIns.User(),
		SaveEntity: userIns,
		Reason:     "TEST",
		Keyname:    kname,
	})

	// create our genesis user vote:
	userInWalletRequestVote := &simpleRequestVote{
		Voter:      genIns.User(),
		IsApproved: true,
	}

	// save the request then save votes:
	saveRequestThenSaveVotesForTests(t, client, pk, repository, user.SDKFunc.CreateRepresentation(), userInWalletRequest, []crypto.PrivateKey{pk}, []*simpleRequestVote{
		userInWalletRequestVote,
	})
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
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// create the representations:
	userRepresentation := user.SDKFunc.CreateRepresentation()
	walletRepresentation := wallet.SDKFunc.CreateRepresentation()

	// spawn bockchain with genesis instance:
	node, client, _, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, genIns)
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
		FromUser:   genIns.User(),
		SaveEntity: userIns,
		Reason:     "TEST",
		Keyname:    kname,
	})

	// create our genesis user vote:
	userInWalletRequestVote := &simpleRequestVote{
		Voter:      genIns.User(),
		IsApproved: true,
	}

	// save the request then save votes:
	saveRequestThenSaveVotesForTests(t, client, pk, repository, userRepresentation, userInWalletRequest, []crypto.PrivateKey{pk}, []*simpleRequestVote{
		userInWalletRequestVote,
	})

	// retrieve the keyname:
	walletKanme, walletKanmeErr := knameRepository.RetrieveByName(wallet.SDKFunc.CreateMetaData().Keyname())
	if walletKanmeErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", walletKanmeErr.Error())
	}

	// update the wallet to increase concensus:
	updateWalletRequest := request.SDKFunc.Create(request.CreateParams{
		FromUser: genIns.User(),
		SaveEntity: wallet.SDKFunc.Create(wallet.CreateParams{
			ID:              wal.ID(),
			Creator:         wal.Creator(),
			ConcensusNeeded: genIns.User().Shares() + userIns.Shares(),
		}),
		Reason:  "TEST",
		Keyname: walletKanme,
	})

	// create our genesis user vote on the wallet update:
	updateWalletRequestVote := &simpleRequestVote{
		Voter:      genIns.User(),
		IsApproved: true,
	}

	// save the new wallet request, then save vote:
	saveRequestThenSaveVotesForTests(t, client, pk, repository, walletRepresentation, updateWalletRequest, []crypto.PrivateKey{pk}, []*simpleRequestVote{
		updateWalletRequestVote,
	})

	// update the wallet to decrease concensus:
	updateAgainWalletRequest := request.SDKFunc.Create(request.CreateParams{
		FromUser: genIns.User(),
		SaveEntity: wallet.SDKFunc.Create(wallet.CreateParams{
			ID:              wal.ID(),
			Creator:         wal.Creator(),
			ConcensusNeeded: genIns.User().Shares(),
		}),
		Reason:  "TEST",
		Keyname: walletKanme,
	})

	// create our genesis user vote on the wallet update:
	updateAgainWalletRequestVoteByGenUser := &simpleRequestVote{
		Voter:      genIns.User(),
		IsApproved: true,
	}

	// create our newly added user vote on the wallet update:
	updateAgainWalletRequestVoteByNewlyAddedUser := &simpleRequestVote{
		Voter:      userIns,
		IsApproved: true,
	}

	// save the new wallet request, then save vote:
	saveRequestThenSaveVotesForTests(t, client, pk, repository, walletRepresentation, updateAgainWalletRequest, []crypto.PrivateKey{pk, userPK}, []*simpleRequestVote{
		updateAgainWalletRequestVoteByGenUser,
		updateAgainWalletRequestVoteByNewlyAddedUser,
	})
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
			Amount: int(math.Floor(float64(genIns.Deposit().Amount() / 2))),
		}),
		To: walletIns,
	})

	// create the repreentations:
	transferRepresentation := transfer.SDKFunc.CreateRepresentation()

	rootPath := filepath.Join("./test_files_TestSaveGenesis_createNewWallet_createPledge_transferPledgeTokens_returnsError")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, client, service, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, genIns)
	defer node.Stop()

	// create the user with the new wallet:
	saveUserWithNewWallet(t, client, pk, service, repository, genIns.User(), userIns)

	// compare the wallets:
	wallet.CompareWalletsForTests(t, walletIns.(wallet.Wallet), userIns.Wallet())

	// save the pledge:
	savePledge(t, client, pk, service, repository, genIns, pldge)

	// transfer the pledge funds, returns error:
	trsf := transfer.SDKFunc.Create(transfer.CreateParams{
		Withdrawal: withdrawal.SDKFunc.Create(withdrawal.CreateParams{
			From:   walletIns,
			Amount: pldge.From().Amount(),
		}),
		Deposit: deposit.SDKFunc.Create(deposit.CreateParams{
			To:     genIns.Deposit().To(),
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
		FromUser:   userIns,
		SaveEntity: trsf,
		Reason:     "TEST",
		Keyname:    kname,
	})

	// create the request service:
	requestService := request.SDKFunc.CreateSDKService(request.CreateSDKServiceParams{
		PK:     walPK,
		Client: client,
	})

	// create the request repository:
	requestRepository := active_request.SDKFunc.CreateRepository(active_request.CreateRepositoryParams{
		EntityRepository: repository,
	})

	// save the request:
	saveRequestErr := requestService.Save(trsfRequest, transferRepresentation)
	if saveRequestErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveRequestErr.Error())
		return
	}

	// retrieve the request:
	retReq, retReqErr := requestRepository.RetrieveByRequest(trsfRequest)
	if retReqErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retReqErr.Error())
		return
	}

	// create our user vote:
	trsfRequestVote := vote.SDKFunc.Create(vote.CreateParams{
		Request:    retReq,
		Voter:      userIns,
		IsApproved: true,
	})

	// create the vote service:
	voteService := vote.SDKFunc.CreateSDKService(vote.CreateSDKServiceParams{
		PK:     walPK,
		Client: client,
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
		IP:     net.ParseIP("127.0.0.1"),
		Port:   8080,
		PubKey: ed25519.GenPrivKey().PubKey(),
		Pledge: pledge.SDKFunc.Create(pledge.CreateParams{
			From: withdrawal.SDKFunc.Create(withdrawal.CreateParams{
				From:   genIns.Deposit().To(),
				Amount: int(math.Floor(float64(genIns.Deposit().Amount() / 2))),
			}),
			To: walletIns,
		}),
	})

	// create the representations:
	validatorRepresentation := validator.SDKFunc.CreateRepresentation()

	rootPath := filepath.Join("./test_files_TestSaveGenesis_createNewWallet_createValidator_Success")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, client, service, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, genIns)
	defer node.Stop()

	// create the user with the new wallet:
	saveUserWithNewWallet(t, client, pk, service, repository, genIns.User(), userIns)

	// compare the wallets:
	wallet.CompareWalletsForTests(t, walletIns.(wallet.Wallet), userIns.Wallet())

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
		FromUser:   genIns.User(),
		SaveEntity: vldator,
		Reason:     "TEST",
		Keyname:    kname,
	})

	// create our user vote:
	validatorRequestVote := &simpleRequestVote{
		Voter:      genIns.User(),
		IsApproved: true,
	}

	// save the new wallet request, then save vote:
	saveRequestThenSaveVotesForTests(t, client, pk, repository, validatorRepresentation, validatorRequest, []crypto.PrivateKey{pk}, []*simpleRequestVote{
		validatorRequestVote,
	})
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
			Amount: int(math.Floor(float64(genIns.Deposit().Amount() / 2))),
		}),
		Deposit: deposit.SDKFunc.Create(deposit.CreateParams{
			To:     walletIns,
			Amount: int(math.Floor(float64(genIns.Deposit().Amount() / 2))),
		}),
	})

	// create the representations:
	transferRepresentation := transfer.SDKFunc.CreateRepresentation()

	rootPath := filepath.Join("./test_files_TestSaveGenesis_createNewWallet_createTransfer_Success")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, client, service, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, genIns)
	defer node.Stop()

	// create the user with the new wallet:
	saveUserWithNewWallet(t, client, pk, service, repository, genIns.User(), userIns)

	// compare the wallets:
	wallet.CompareWalletsForTests(t, walletIns.(wallet.Wallet), userIns.Wallet())

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
		FromUser:   genIns.User(),
		SaveEntity: trsf,
		Reason:     "TEST",
		Keyname:    kname,
	})

	// create our user vote:
	trsfRequestVote := &simpleRequestVote{
		Voter:      genIns.User(),
		IsApproved: true,
	}

	// save the new wallet request, then save vote:
	saveRequestThenSaveVotesForTests(t, client, pk, repository, transferRepresentation, trsfRequest, []crypto.PrivateKey{pk}, []*simpleRequestVote{
		trsfRequestVote,
	})
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
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, client, service, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, genIns)
	defer node.Stop()

	// save the link:
	saveLink(t, client, pk, service, repository, genIns.User(), lnk)
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
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, client, service, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, genIns)
	defer node.Stop()

	// save the link:
	saveNode(t, client, pk, service, repository, genIns.User(), lnk, nod)
}

func TestSaveGenesis_createCategory_thenCreateCategoryWithParent_Success(t *testing.T) {
	// variables:
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	pubKey := pk.PublicKey()
	genIns := genesis.CreateGenesisWithPubKeyForTests(pubKey)

	cat := category.CreateCategoryForTests()
	catWithParent := category.CreateCategoryWithParentForTests(cat)

	rootPath := filepath.Join("./test_files_TestSaveGenesis_createCategory_thenCreateCategoryWithParent_Success")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, client, service, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, genIns)
	defer node.Stop()

	// create the category:
	saveCategory(t, client, pk, service, repository, genIns, cat)

	// create the category with parent:
	saveCategory(t, client, pk, service, repository, genIns, catWithParent)
}

func TestSaveGenesis_saveCategory_saveProposal_Success(t *testing.T) {
	// variables:
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	pubKey := pk.PublicKey()
	genIns := genesis.CreateGenesisWithPubKeyForTests(pubKey)

	cat := category.CreateCategoryForTests()
	prop := proposal.CreateProposalWithCategoryForTests(cat)

	rootPath := filepath.Join("./test_files_TestSaveGenesis_saveCategory_saveProposal_Success")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, client, service, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, genIns)
	defer node.Stop()

	// save the category:
	saveCategory(t, client, pk, service, repository, genIns, cat)

	// save proposal:
	saveProposal(t, client, pk, service, repository, genIns, prop)
}

func TestSaveGenesis_saveCategory_saveProposal_saveCommunityProject_Success(t *testing.T) {
	// variables:
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	pubKey := pk.PublicKey()
	genIns := genesis.CreateGenesisWithPubKeyForTests(pubKey)

	cat := category.CreateCategoryForTests()
	prop := proposal.CreateProposalWithCategoryForTests(cat)
	communityProject := community_project.CreateProjectWithProposalForTests(prop)

	rootPath := filepath.Join("./test_files_TestSaveGenesis_saveCategory_saveProposal_saveCommunityProject_Success")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, client, service, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, genIns)
	defer node.Stop()

	// save the category:
	saveCategory(t, client, pk, service, repository, genIns, cat)

	// save proposal:
	saveProposal(t, client, pk, service, repository, genIns, prop)

	// save approved project:
	saveCommunityProject(t, client, pk, service, repository, genIns, communityProject)
}

func TestSaveGenesis_saveCategory_saveProposal_saveCommunityProject_saveProject_Success(t *testing.T) {
	// variables:
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	pubKey := pk.PublicKey()
	genIns := genesis.CreateGenesisWithPubKeyForTests(pubKey)

	cat := category.CreateCategoryForTests()
	prop := proposal.CreateProposalWithCategoryForTests(cat)
	communityProject := community_project.CreateProjectWithProposalForTests(prop)
	ownerUser := user.CreateUserForTests()
	managerUser := user.CreateUserForTests()
	linkerUser := user.CreateUserForTests()
	proj := project.CreateProjectWithCommunityProjectAndWallets(communityProject, ownerUser.Wallet(), managerUser.Wallet(), linkerUser.Wallet())

	rootPath := filepath.Join("./test_files_TestSaveGenesis_saveCategory_saveProposal_saveCommunityProject_saveProject_Success")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, client, service, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, genIns)
	defer node.Stop()

	// save the wallets in project:
	saveUserWithNewWallet(t, client, pk, service, repository, genIns.User(), ownerUser)
	saveUserWithNewWallet(t, client, pk, service, repository, genIns.User(), managerUser)
	saveUserWithNewWallet(t, client, pk, service, repository, genIns.User(), linkerUser)

	// save a transfer between our genesis and the manager:
	saveTransfer(t, client, pk, service, repository, genIns, transfer.SDKFunc.Create(transfer.CreateParams{
		Withdrawal: withdrawal.SDKFunc.Create(withdrawal.CreateParams{
			From:   genIns.Deposit().To(),
			Amount: 50,
		}),
		Deposit: deposit.SDKFunc.Create(deposit.CreateParams{
			To:     proj.Manager(),
			Amount: 50,
		}),
	}))

	// save a transfer between our genesis and the linker:
	saveTransfer(t, client, pk, service, repository, genIns, transfer.SDKFunc.Create(transfer.CreateParams{
		Withdrawal: withdrawal.SDKFunc.Create(withdrawal.CreateParams{
			From:   genIns.Deposit().To(),
			Amount: 50,
		}),
		Deposit: deposit.SDKFunc.Create(deposit.CreateParams{
			To:     proj.Linker(),
			Amount: 50,
		}),
	}))

	// make the manager pledge:
	savePledge(t, client, pk, service, repository, genIns, pledge.SDKFunc.Create(pledge.CreateParams{
		From: withdrawal.SDKFunc.Create(withdrawal.CreateParams{
			From:   proj.Manager(),
			Amount: proj.Project().Proposal().ManagerPledgeNeeded(),
		}),
		To: proj.Owner(),
	}))

	// make the linker pledge:
	savePledge(t, client, pk, service, repository, genIns, pledge.SDKFunc.Create(pledge.CreateParams{
		From: withdrawal.SDKFunc.Create(withdrawal.CreateParams{
			From:   proj.Linker(),
			Amount: proj.Project().Proposal().LinkerPledgeNeeded(),
		}),
		To: proj.Owner(),
	}))

	// save the category:
	saveCategory(t, client, pk, service, repository, genIns, cat)

	// save proposal:
	saveProposal(t, client, pk, service, repository, genIns, prop)

	// save approved project:
	saveCommunityProject(t, client, pk, service, repository, genIns, communityProject)

	// save project:
	saveProject(t, client, pk, service, repository, genIns, proj)
}

func TestSaveGenesis_saveCategory_saveProposal_saveCommunityProject_saveProject_saveMilestone_Success(t *testing.T) {
	// variables:
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	pubKey := pk.PublicKey()
	genIns := genesis.CreateGenesisWithPubKeyForTests(pubKey)

	cat := category.CreateCategoryForTests()
	prop := proposal.CreateProposalWithCategoryForTests(cat)
	communityProject := community_project.CreateProjectWithProposalForTests(prop)
	ownerUser := user.CreateUserForTests()
	managerUser := user.CreateUserForTests()
	linkerUser := user.CreateUserForTests()
	proj := project.CreateProjectWithCommunityProjectAndWallets(communityProject, ownerUser.Wallet(), managerUser.Wallet(), linkerUser.Wallet())
	mils := milestone.CreateMilestoneWithProjectForTests(proj)

	rootPath := filepath.Join("./test_files_TestSaveGenesis_saveCategory_saveProposal_saveCommunityProject_saveProject_saveMilestone_Success")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, client, service, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, genIns)
	defer node.Stop()

	// save the wallets in project:
	saveUserWithNewWallet(t, client, pk, service, repository, genIns.User(), ownerUser)
	saveUserWithNewWallet(t, client, pk, service, repository, genIns.User(), managerUser)
	saveUserWithNewWallet(t, client, pk, service, repository, genIns.User(), linkerUser)

	// save a transfer between our genesis and the manager:
	saveTransfer(t, client, pk, service, repository, genIns, transfer.SDKFunc.Create(transfer.CreateParams{
		Withdrawal: withdrawal.SDKFunc.Create(withdrawal.CreateParams{
			From:   genIns.Deposit().To(),
			Amount: 50,
		}),
		Deposit: deposit.SDKFunc.Create(deposit.CreateParams{
			To:     proj.Manager(),
			Amount: 50,
		}),
	}))

	// save a transfer between our genesis and the linker:
	saveTransfer(t, client, pk, service, repository, genIns, transfer.SDKFunc.Create(transfer.CreateParams{
		Withdrawal: withdrawal.SDKFunc.Create(withdrawal.CreateParams{
			From:   genIns.Deposit().To(),
			Amount: 50,
		}),
		Deposit: deposit.SDKFunc.Create(deposit.CreateParams{
			To:     proj.Linker(),
			Amount: 50,
		}),
	}))

	// make the manager pledge:
	savePledge(t, client, pk, service, repository, genIns, pledge.SDKFunc.Create(pledge.CreateParams{
		From: withdrawal.SDKFunc.Create(withdrawal.CreateParams{
			From:   proj.Manager(),
			Amount: proj.Project().Proposal().ManagerPledgeNeeded(),
		}),
		To: proj.Owner(),
	}))

	// make the linker pledge:
	savePledge(t, client, pk, service, repository, genIns, pledge.SDKFunc.Create(pledge.CreateParams{
		From: withdrawal.SDKFunc.Create(withdrawal.CreateParams{
			From:   proj.Linker(),
			Amount: proj.Project().Proposal().LinkerPledgeNeeded(),
		}),
		To: proj.Owner(),
	}))

	// save the category:
	saveCategory(t, client, pk, service, repository, genIns, cat)

	// save proposal:
	saveProposal(t, client, pk, service, repository, genIns, prop)

	// save approved project:
	saveCommunityProject(t, client, pk, service, repository, genIns, communityProject)

	// save project:
	saveProject(t, client, pk, service, repository, genIns, proj)

	// save milestone:
	saveMilestone(t, client, pk, service, repository, genIns, mils)
}

func TestSaveGenesis_saveCategory_saveProposal_saveCommunityProject_saveProject_saveFeature_Success(t *testing.T) {
	// variables:
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	pubKey := pk.PublicKey()
	genIns := genesis.CreateGenesisWithPubKeyForTests(pubKey)

	cat := category.CreateCategoryForTests()
	prop := proposal.CreateProposalWithCategoryForTests(cat)
	communityProject := community_project.CreateProjectWithProposalForTests(prop)
	ownerUser := user.CreateUserForTests()
	managerUser := user.CreateUserForTests()
	linkerUser := user.CreateUserForTests()
	proj := project.CreateProjectWithCommunityProjectAndWallets(communityProject, ownerUser.Wallet(), managerUser.Wallet(), linkerUser.Wallet())
	feat := feature.CreateFeatureWithProjectAndCreatedByUser(proj, genIns.User())

	rootPath := filepath.Join("./test_files_TestSaveGenesis_saveCategory_saveProposal_saveCommunityProject_saveProject_saveFeature_Success")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, client, service, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, genIns)
	defer node.Stop()

	// save the wallets in project:
	saveUserWithNewWallet(t, client, pk, service, repository, genIns.User(), ownerUser)
	saveUserWithNewWallet(t, client, pk, service, repository, genIns.User(), managerUser)
	saveUserWithNewWallet(t, client, pk, service, repository, genIns.User(), linkerUser)

	// save a transfer between our genesis and the manager:
	saveTransfer(t, client, pk, service, repository, genIns, transfer.SDKFunc.Create(transfer.CreateParams{
		Withdrawal: withdrawal.SDKFunc.Create(withdrawal.CreateParams{
			From:   genIns.Deposit().To(),
			Amount: 50,
		}),
		Deposit: deposit.SDKFunc.Create(deposit.CreateParams{
			To:     proj.Manager(),
			Amount: 50,
		}),
	}))

	// save a transfer between our genesis and the linker:
	saveTransfer(t, client, pk, service, repository, genIns, transfer.SDKFunc.Create(transfer.CreateParams{
		Withdrawal: withdrawal.SDKFunc.Create(withdrawal.CreateParams{
			From:   genIns.Deposit().To(),
			Amount: 50,
		}),
		Deposit: deposit.SDKFunc.Create(deposit.CreateParams{
			To:     proj.Linker(),
			Amount: 50,
		}),
	}))

	// make the manager pledge:
	savePledge(t, client, pk, service, repository, genIns, pledge.SDKFunc.Create(pledge.CreateParams{
		From: withdrawal.SDKFunc.Create(withdrawal.CreateParams{
			From:   proj.Manager(),
			Amount: proj.Project().Proposal().ManagerPledgeNeeded(),
		}),
		To: proj.Owner(),
	}))

	// make the linker pledge:
	savePledge(t, client, pk, service, repository, genIns, pledge.SDKFunc.Create(pledge.CreateParams{
		From: withdrawal.SDKFunc.Create(withdrawal.CreateParams{
			From:   proj.Linker(),
			Amount: proj.Project().Proposal().LinkerPledgeNeeded(),
		}),
		To: proj.Owner(),
	}))

	// save the category:
	saveCategory(t, client, pk, service, repository, genIns, cat)

	// save proposal:
	saveProposal(t, client, pk, service, repository, genIns, prop)

	// save approved project:
	saveCommunityProject(t, client, pk, service, repository, genIns, communityProject)

	// save project:
	saveProject(t, client, pk, service, repository, genIns, proj)

	// save feature:
	saveFeature(t, client, pk, service, repository, genIns, feat)
}

func TestSaveGenesis_saveCategory_saveProposal_saveCommunityProject_saveProject_saveTask_Success(t *testing.T) {
	// variables:
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	pubKey := pk.PublicKey()
	genIns := genesis.CreateGenesisWithPubKeyForTests(pubKey)

	cat := category.CreateCategoryForTests()
	prop := proposal.CreateProposalWithCategoryForTests(cat)
	communityProject := community_project.CreateProjectWithProposalForTests(prop)
	ownerUser := user.CreateUserForTests()
	managerUser := user.CreateUserForTests()
	linkerUser := user.CreateUserForTests()
	proj := project.CreateProjectWithCommunityProjectAndWallets(communityProject, ownerUser.Wallet(), managerUser.Wallet(), linkerUser.Wallet())
	mils := milestone.CreateMilestoneWithProjectForTests(proj)
	tsk := task.CreateTaskWithMilestoneAndUser(mils, genIns.User())

	rootPath := filepath.Join("./test_files_TestSaveGenesis_saveCategory_saveProposal_saveCommunityProject_saveProject_saveTask_Success")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, client, service, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, genIns)
	defer node.Stop()

	// save the wallets in project:
	saveUserWithNewWallet(t, client, pk, service, repository, genIns.User(), ownerUser)
	saveUserWithNewWallet(t, client, pk, service, repository, genIns.User(), managerUser)
	saveUserWithNewWallet(t, client, pk, service, repository, genIns.User(), linkerUser)

	// save a transfer between our genesis and the manager:
	saveTransfer(t, client, pk, service, repository, genIns, transfer.SDKFunc.Create(transfer.CreateParams{
		Withdrawal: withdrawal.SDKFunc.Create(withdrawal.CreateParams{
			From:   genIns.Deposit().To(),
			Amount: 50,
		}),
		Deposit: deposit.SDKFunc.Create(deposit.CreateParams{
			To:     proj.Manager(),
			Amount: 50,
		}),
	}))

	// save a transfer between our genesis and the linker:
	saveTransfer(t, client, pk, service, repository, genIns, transfer.SDKFunc.Create(transfer.CreateParams{
		Withdrawal: withdrawal.SDKFunc.Create(withdrawal.CreateParams{
			From:   genIns.Deposit().To(),
			Amount: 50,
		}),
		Deposit: deposit.SDKFunc.Create(deposit.CreateParams{
			To:     proj.Linker(),
			Amount: 50,
		}),
	}))

	// make the manager pledge:
	savePledge(t, client, pk, service, repository, genIns, pledge.SDKFunc.Create(pledge.CreateParams{
		From: withdrawal.SDKFunc.Create(withdrawal.CreateParams{
			From:   proj.Manager(),
			Amount: proj.Project().Proposal().ManagerPledgeNeeded(),
		}),
		To: proj.Owner(),
	}))

	// make the linker pledge:
	savePledge(t, client, pk, service, repository, genIns, pledge.SDKFunc.Create(pledge.CreateParams{
		From: withdrawal.SDKFunc.Create(withdrawal.CreateParams{
			From:   proj.Linker(),
			Amount: proj.Project().Proposal().LinkerPledgeNeeded(),
		}),
		To: proj.Owner(),
	}))

	// save the category:
	saveCategory(t, client, pk, service, repository, genIns, cat)

	// save proposal:
	saveProposal(t, client, pk, service, repository, genIns, prop)

	// save approved project:
	saveCommunityProject(t, client, pk, service, repository, genIns, communityProject)

	// save project:
	saveProject(t, client, pk, service, repository, genIns, proj)

	// save milestone:
	saveMilestone(t, client, pk, service, repository, genIns, mils)

	// save task:
	saveTask(t, client, pk, service, repository, genIns, tsk)
}

func TestSaveGenesis_saveCategory_saveProposal_saveCommunityProject_saveProject_saveTask_savePledgeTask_Success(t *testing.T) {
	// variables:
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	pubKey := pk.PublicKey()
	genIns := genesis.CreateGenesisWithPubKeyForTests(pubKey)

	cat := category.CreateCategoryForTests()
	prop := proposal.CreateProposalWithCategoryForTests(cat)
	communityProject := community_project.CreateProjectWithProposalForTests(prop)
	ownerUser := user.CreateUserForTests()
	managerUser := user.CreateUserForTests()
	linkerUser := user.CreateUserForTests()
	proj := project.CreateProjectWithCommunityProjectAndWallets(communityProject, ownerUser.Wallet(), managerUser.Wallet(), linkerUser.Wallet())

	mils := milestone.CreateMilestoneWithProjectForTests(proj)
	tsk := task.CreateTaskWithMilestoneAndUser(mils, genIns.User())
	pldge := pledge.CreatePledgeWithWalletForTests(genIns.User().Wallet(), tsk.Milestone().Wallet(), tsk.PledgeNeeded())
	pledgeTask := pledge_task.CreateTaskWithMilestoneTaskAndPledge(tsk, pldge)

	rootPath := filepath.Join("./test_files_TestSaveGenesis_saveCategory_saveProposal_saveCommunityProject_saveProject_saveTask_savePledgeTask_Success")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, client, service, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, genIns)
	defer node.Stop()

	// save the wallets in project:
	saveUserWithNewWallet(t, client, pk, service, repository, genIns.User(), ownerUser)
	saveUserWithNewWallet(t, client, pk, service, repository, genIns.User(), managerUser)
	saveUserWithNewWallet(t, client, pk, service, repository, genIns.User(), linkerUser)

	// save a transfer between our genesis and the manager:
	saveTransfer(t, client, pk, service, repository, genIns, transfer.SDKFunc.Create(transfer.CreateParams{
		Withdrawal: withdrawal.SDKFunc.Create(withdrawal.CreateParams{
			From:   genIns.Deposit().To(),
			Amount: 50,
		}),
		Deposit: deposit.SDKFunc.Create(deposit.CreateParams{
			To:     proj.Manager(),
			Amount: 50,
		}),
	}))

	// save a transfer between our genesis and the linker:
	saveTransfer(t, client, pk, service, repository, genIns, transfer.SDKFunc.Create(transfer.CreateParams{
		Withdrawal: withdrawal.SDKFunc.Create(withdrawal.CreateParams{
			From:   genIns.Deposit().To(),
			Amount: 50,
		}),
		Deposit: deposit.SDKFunc.Create(deposit.CreateParams{
			To:     proj.Linker(),
			Amount: 50,
		}),
	}))

	// make the manager pledge:
	savePledge(t, client, pk, service, repository, genIns, pledge.SDKFunc.Create(pledge.CreateParams{
		From: withdrawal.SDKFunc.Create(withdrawal.CreateParams{
			From:   proj.Manager(),
			Amount: proj.Project().Proposal().ManagerPledgeNeeded(),
		}),
		To: proj.Owner(),
	}))

	// make the linker pledge:
	savePledge(t, client, pk, service, repository, genIns, pledge.SDKFunc.Create(pledge.CreateParams{
		From: withdrawal.SDKFunc.Create(withdrawal.CreateParams{
			From:   proj.Linker(),
			Amount: proj.Project().Proposal().LinkerPledgeNeeded(),
		}),
		To: proj.Owner(),
	}))

	// save the category:
	saveCategory(t, client, pk, service, repository, genIns, cat)

	// save proposal:
	saveProposal(t, client, pk, service, repository, genIns, prop)

	// save approved project:
	saveCommunityProject(t, client, pk, service, repository, genIns, communityProject)

	// save project:
	saveProject(t, client, pk, service, repository, genIns, proj)

	// save milestone:
	saveMilestone(t, client, pk, service, repository, genIns, mils)

	// save task:
	saveTask(t, client, pk, service, repository, genIns, tsk)

	// save pledge task:
	savePledgeTask(t, client, pk, service, repository, genIns, pledgeTask)
}

func TestSaveGenesis_saveCategory_saveProposal_saveCommunityProject_saveProject_saveTask_saveCompletedTask_Success(t *testing.T) {
	// variables:
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	pubKey := pk.PublicKey()
	genIns := genesis.CreateGenesisWithPubKeyForTests(pubKey)

	cat := category.CreateCategoryForTests()
	prop := proposal.CreateProposalWithCategoryForTests(cat)
	communityProject := community_project.CreateProjectWithProposalForTests(prop)
	ownerUser := user.CreateUserForTests()
	managerUser := user.CreateUserForTests()
	linkerUser := user.CreateUserForTests()
	proj := project.CreateProjectWithCommunityProjectAndWallets(communityProject, ownerUser.Wallet(), managerUser.Wallet(), linkerUser.Wallet())

	mils := milestone.CreateMilestoneWithProjectForTests(proj)
	tsk := task.CreateTaskWithMilestoneAndUser(mils, genIns.User())
	pldge := pledge.CreatePledgeWithWalletForTests(genIns.User().Wallet(), tsk.Milestone().Wallet(), tsk.PledgeNeeded())
	pledgeTask := pledge_task.CreateTaskWithMilestoneTaskAndPledge(tsk, pldge)
	completedTask := completed_task.CreateTaskWithMilestoneTask(tsk)

	rootPath := filepath.Join("./test_files_TestSaveGenesis_saveCategory_saveProposal_saveCommunityProject_saveProject_saveTask_saveCompletedTask_Success")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, client, service, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, genIns)
	defer node.Stop()

	// save the wallets in project:
	saveUserWithNewWallet(t, client, pk, service, repository, genIns.User(), ownerUser)
	saveUserWithNewWallet(t, client, pk, service, repository, genIns.User(), managerUser)
	saveUserWithNewWallet(t, client, pk, service, repository, genIns.User(), linkerUser)

	// save a transfer between our genesis and the manager:
	saveTransfer(t, client, pk, service, repository, genIns, transfer.SDKFunc.Create(transfer.CreateParams{
		Withdrawal: withdrawal.SDKFunc.Create(withdrawal.CreateParams{
			From:   genIns.Deposit().To(),
			Amount: 50,
		}),
		Deposit: deposit.SDKFunc.Create(deposit.CreateParams{
			To:     proj.Manager(),
			Amount: 50,
		}),
	}))

	// save a transfer between our genesis and the linker:
	saveTransfer(t, client, pk, service, repository, genIns, transfer.SDKFunc.Create(transfer.CreateParams{
		Withdrawal: withdrawal.SDKFunc.Create(withdrawal.CreateParams{
			From:   genIns.Deposit().To(),
			Amount: 50,
		}),
		Deposit: deposit.SDKFunc.Create(deposit.CreateParams{
			To:     proj.Linker(),
			Amount: 50,
		}),
	}))

	// make the manager pledge:
	savePledge(t, client, pk, service, repository, genIns, pledge.SDKFunc.Create(pledge.CreateParams{
		From: withdrawal.SDKFunc.Create(withdrawal.CreateParams{
			From:   proj.Manager(),
			Amount: proj.Project().Proposal().ManagerPledgeNeeded(),
		}),
		To: proj.Owner(),
	}))

	// make the linker pledge:
	savePledge(t, client, pk, service, repository, genIns, pledge.SDKFunc.Create(pledge.CreateParams{
		From: withdrawal.SDKFunc.Create(withdrawal.CreateParams{
			From:   proj.Linker(),
			Amount: proj.Project().Proposal().LinkerPledgeNeeded(),
		}),
		To: proj.Owner(),
	}))

	// save the category:
	saveCategory(t, client, pk, service, repository, genIns, cat)

	// save proposal:
	saveProposal(t, client, pk, service, repository, genIns, prop)

	// save approved project:
	saveCommunityProject(t, client, pk, service, repository, genIns, communityProject)

	// save project:
	saveProject(t, client, pk, service, repository, genIns, proj)

	// save milestone:
	saveMilestone(t, client, pk, service, repository, genIns, mils)

	// save task:
	saveTask(t, client, pk, service, repository, genIns, tsk)

	// save pledge task:
	savePledgeTask(t, client, pk, service, repository, genIns, pledgeTask)

	// save the completed task:
	saveCompletedTask(t, client, pk, service, repository, genIns, completedTask)
}
