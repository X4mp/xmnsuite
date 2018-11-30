package core

import (
	"math"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/genesis"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/entities/pledge"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/entities/transfer"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/entities/validator"
	"github.com/xmnservices/xmnsuite/blockchains/core/request"
	"github.com/xmnservices/xmnsuite/blockchains/core/request/vote"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/entities/developer"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/entities/developer/entities/milestone"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/entities/developer/entities/project"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/entities/developer/entities/task"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/entities/link"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/entities/node"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/withdrawal"
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

func TestSaveGenesis_createWallet_Success(t *testing.T) {
	// variables:
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	pubKey := pk.PublicKey()
	genIns := genesis.CreateGenesisWithPubKeyForTests(pubKey)

	newWalletPK := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	newWalletPubKey := newWalletPK.PublicKey()
	walletIns := wallet.CreateWalletWithPublicKeyForTests(newWalletPubKey)
	rootPath := filepath.Join("./test_files_TestSaveGenesis_createWallet_Success")
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
	rootPath := filepath.Join("./test_files_TestSaveGenesis_createWalletWithSameCreator_Success")
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
	rootPath := filepath.Join("./test_files_TestSaveGenesis_createWalletAlreadyInGenesis_returnsError")
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
	rootPath := filepath.Join("./test_files_TestSaveGenesis_createWallet_addUserToWallet_addAnotherUserToWallerWithSamePublicKey_saveVotesWithEnoughSharesToPass_returnsError")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// create the user in wallet request:
	req := request.SDKFunc.Create(request.CreateParams{
		FromUser:       genIns.User(),
		NewEntity:      userIns,
		EntityMetaData: user.SDKFunc.CreateMetaData(),
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
		PK:              pk,
		Client:          client,
		CreateRouteFunc: createEntityVoteRouteFunc(),
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
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// create the user in wallet request:
	userInWalletRequest := request.SDKFunc.Create(request.CreateParams{
		FromUser:       genIns.User(),
		NewEntity:      userIns,
		EntityMetaData: user.SDKFunc.CreateMetaData(),
	})

	// create our genesis user vote:
	userInWalletRequestVote := vote.SDKFunc.Create(vote.CreateParams{
		Request:    userInWalletRequest,
		Voter:      genIns.User(),
		IsApproved: true,
	})

	// save the new wallet request, then save vote:
	node, _, _, _, _ := spawnBlockchainWithGenesisThenSaveRequestThenSaveVotesForTests(t, pk, rootPath, genIns, user.SDKFunc.CreateRepresentation(), userInWalletRequest, []crypto.PrivateKey{pk}, []vote.Vote{
		userInWalletRequestVote,
	}, createEntityVoteRouteFunc())

	defer node.Stop()
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

	// create the user in wallet request:
	userInWalletRequest := request.SDKFunc.Create(request.CreateParams{
		FromUser:       genIns.User(),
		NewEntity:      userIns,
		EntityMetaData: user.SDKFunc.CreateMetaData(),
	})

	// create our genesis user vote:
	userInWalletRequestVote := vote.SDKFunc.Create(vote.CreateParams{
		Request:    userInWalletRequest,
		Voter:      genIns.User(),
		IsApproved: true,
	})

	// save the new wallet request, then save vote:
	node, client, _, repository, _ := spawnBlockchainWithGenesisThenSaveRequestThenSaveVotesForTests(t, pk, rootPath, genIns, userRepresentation, userInWalletRequest, []crypto.PrivateKey{pk}, []vote.Vote{
		userInWalletRequestVote,
	}, createEntityVoteRouteFunc())

	defer node.Stop()

	// update the wallet to increase concensus:
	updateWalletRequest := request.SDKFunc.Create(request.CreateParams{
		FromUser: genIns.User(),
		NewEntity: wallet.SDKFunc.Create(wallet.CreateParams{
			ID:              wal.ID(),
			Creator:         wal.Creator(),
			ConcensusNeeded: genIns.User().Shares() + userIns.Shares(),
		}),
		EntityMetaData: wallet.SDKFunc.CreateMetaData(),
	})

	// create our genesis user vote on the wallet update:
	updateWalletRequestVote := vote.SDKFunc.Create(vote.CreateParams{
		Request:    updateWalletRequest,
		Voter:      genIns.User(),
		IsApproved: true,
	})

	// save the new wallet request, then save vote:
	saveRequestThenSaveVotesForTests(t, client, pk, repository, walletRepresentation, updateWalletRequest, []crypto.PrivateKey{pk}, []vote.Vote{
		updateWalletRequestVote,
	}, createEntityVoteRouteFunc())

	// update the wallet to decrease concensus:
	updateAgainWalletRequest := request.SDKFunc.Create(request.CreateParams{
		FromUser: genIns.User(),
		NewEntity: wallet.SDKFunc.Create(wallet.CreateParams{
			ID:              wal.ID(),
			Creator:         wal.Creator(),
			ConcensusNeeded: genIns.User().Shares(),
		}),
		EntityMetaData: wallet.SDKFunc.CreateMetaData(),
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
	saveRequestThenSaveVotesForTests(t, client, pk, repository, walletRepresentation, updateAgainWalletRequest, []crypto.PrivateKey{pk, userPK}, []vote.Vote{
		updateAgainWalletRequestVoteByGenUser,
		updateAgainWalletRequestVoteByNewlyAddedUser,
	}, createEntityVoteRouteFunc())
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
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, client, service, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, genIns)
	defer node.Stop()

	// save the pledge:
	savePledge(t, client, pk, service, repository, genIns.User(), userIns, pldge)

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

	// create the user in wallet request:
	trsfRequest := request.SDKFunc.Create(request.CreateParams{
		FromUser:       userIns,
		NewEntity:      trsf,
		EntityMetaData: transfer.SDKFunc.CreateMetaData(),
	})

	// create our user vote:
	trsfRequestVote := vote.SDKFunc.Create(vote.CreateParams{
		Request:    trsfRequest,
		Voter:      userIns,
		IsApproved: true,
	})

	// create the request service:
	requestService := request.SDKFunc.CreateSDKService(request.CreateSDKServiceParams{
		PK:     walPK,
		Client: client,
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
		CreateRouteFunc: createEntityVoteRouteFunc(),
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

	// create the repreentations:
	walletRepresentation := wallet.SDKFunc.CreateRepresentation()
	userRepresentation := user.SDKFunc.CreateRepresentation()
	validatorRepresentation := validator.SDKFunc.CreateRepresentation()

	rootPath := filepath.Join("./test_files_TestSaveGenesis_createNewWallet_createValidator_Success")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, client, service, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, genIns)
	defer node.Stop()

	// save the new wallet:
	savedWallet := saveEntityThenRetrieveEntityByIDThenDeleteEntityByID(t, walletIns, walletRepresentation, service, repository)

	// compare the wallets:
	wallet.CompareWalletsForTests(t, walletIns.(wallet.Wallet), savedWallet.(wallet.Wallet))

	// create the user in wallet request:
	userInWalletRequest := request.SDKFunc.Create(request.CreateParams{
		FromUser:       genIns.User(),
		NewEntity:      userIns,
		EntityMetaData: user.SDKFunc.CreateMetaData(),
	})

	// create our user vote:
	userInWalletRequestVote := vote.SDKFunc.Create(vote.CreateParams{
		Request:    userInWalletRequest,
		Voter:      genIns.User(),
		IsApproved: true,
	})

	// save the new wallet request, then save vote:
	saveRequestThenSaveVotesForTests(t, client, pk, repository, userRepresentation, userInWalletRequest, []crypto.PrivateKey{pk}, []vote.Vote{
		userInWalletRequestVote,
	}, createEntityVoteRouteFunc())

	// create the user in validator request:
	validatorRequest := request.SDKFunc.Create(request.CreateParams{
		FromUser:       genIns.User(),
		NewEntity:      vldator,
		EntityMetaData: validator.SDKFunc.CreateMetaData(),
	})

	// create our user vote:
	validatorRequestVote := vote.SDKFunc.Create(vote.CreateParams{
		Request:    validatorRequest,
		Voter:      genIns.User(),
		IsApproved: true,
	})

	// save the new wallet request, then save vote:
	saveRequestThenSaveVotesForTests(t, client, pk, repository, validatorRepresentation, validatorRequest, []crypto.PrivateKey{pk}, []vote.Vote{
		validatorRequestVote,
	}, createEntityVoteRouteFunc())
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

	// create the repreentations:
	walletRepresentation := wallet.SDKFunc.CreateRepresentation()
	userRepresentation := user.SDKFunc.CreateRepresentation()
	transferRepresentation := transfer.SDKFunc.CreateRepresentation()

	rootPath := filepath.Join("./test_files_TestSaveGenesis_createNewWallet_createTransfer_Success")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, client, service, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, genIns)
	defer node.Stop()

	// save the new wallet:
	savedWallet := saveEntityThenRetrieveEntityByIDThenDeleteEntityByID(t, walletIns, walletRepresentation, service, repository)

	// compare the wallets:
	wallet.CompareWalletsForTests(t, walletIns.(wallet.Wallet), savedWallet.(wallet.Wallet))

	// create the user in wallet request:
	userInWalletRequest := request.SDKFunc.Create(request.CreateParams{
		FromUser:       genIns.User(),
		NewEntity:      userIns,
		EntityMetaData: user.SDKFunc.CreateMetaData(),
	})

	// create our user vote:
	userInWalletRequestVote := vote.SDKFunc.Create(vote.CreateParams{
		Request:    userInWalletRequest,
		Voter:      genIns.User(),
		IsApproved: true,
	})

	// save the new wallet request, then save vote:
	saveRequestThenSaveVotesForTests(t, client, pk, repository, userRepresentation, userInWalletRequest, []crypto.PrivateKey{pk}, []vote.Vote{
		userInWalletRequestVote,
	}, createEntityVoteRouteFunc())

	// create the user in wallet request:
	trsfRequest := request.SDKFunc.Create(request.CreateParams{
		FromUser:       genIns.User(),
		NewEntity:      trsf,
		EntityMetaData: transfer.SDKFunc.CreateMetaData(),
	})

	// create our user vote:
	trsfRequestVote := vote.SDKFunc.Create(vote.CreateParams{
		Request:    trsfRequest,
		Voter:      genIns.User(),
		IsApproved: true,
	})

	// save the new wallet request, then save vote:
	saveRequestThenSaveVotesForTests(t, client, pk, repository, transferRepresentation, trsfRequest, []crypto.PrivateKey{pk}, []vote.Vote{
		trsfRequestVote,
	}, createEntityVoteRouteFunc())
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

	rootPath := filepath.Join("./test_TestSaveGenesis_CreateLink_voteOnLink_Success")
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

	rootPath := filepath.Join("./test_TestSaveGenesis_CreateLink_voteOnLink_CreateNode_voteOnNode_Success")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, client, service, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, genIns)
	defer node.Stop()

	// save the link:
	saveNode(t, client, pk, service, repository, genIns.User(), lnk, nod)
}

func TestSaveGenesis_createPledge_voteOnPledge_createDeveloper_voteOnDeveloper_Success(t *testing.T) {
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

	dev := developer.SDKFunc.Create(developer.CreateParams{
		Pledge: pldge,
		User:   genIns.User(),
		Name:   "Steve",
		Resume: "this is the content of my resume",
	})

	rootPath := filepath.Join("./test_files_TestSaveGenesis_createPledge_voteOnPledge_createDeveloper_voteOnDeveloper_Success")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, client, service, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, genIns)
	defer node.Stop()

	// save the developer:
	saveDeveloper(t, client, pk, service, repository, genIns.User(), userIns, pldge, dev)
}

func TestSaveGenesis_createProject_voteOnProjectWithDeveloperUser_Success(t *testing.T) {
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

	dev := developer.SDKFunc.Create(developer.CreateParams{
		Pledge: pldge,
		User:   genIns.User(),
		Name:   "Steve",
		Resume: "this is a resume",
	})

	proj := project.SDKFunc.Create(project.CreateParams{
		Title:       "This is a project",
		Description: "This is the project description",
	})

	rootPath := filepath.Join("./test_files_TestSaveGenesis_createProject_voteOnProjectWithDeveloperUser_Success")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, client, service, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, genIns)
	defer node.Stop()

	// save the project:
	saveProject(t, client, pk, service, repository, genIns.User(), userIns, pldge, dev, proj)
}

func TestSaveGenesis_CreateProject_voteOnProject_withoutADeveloperUser_returnsError(t *testing.T) {
	// variables:
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	pubKey := pk.PublicKey()
	genIns := genesis.CreateGenesisWithPubKeyForTests(pubKey)

	proj := project.SDKFunc.Create(project.CreateParams{
		Title:       "This is a project",
		Description: "This is the project description",
	})

	// create the representations:
	projectRepresentation := project.SDKFunc.CreateRepresentation()

	rootPath := filepath.Join("./test_files_TestSaveGenesis_CreateProject_voteOnProject_withoutADeveloperUser_returnsError")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, client, _, _ := spawnBlockchainWithGenesisForTests(t, pk, rootPath, genIns)
	defer node.Stop()

	// create the project request:
	newProjectRequest := request.SDKFunc.Create(request.CreateParams{
		FromUser:       genIns.User(),
		NewEntity:      proj,
		EntityMetaData: project.SDKFunc.CreateMetaData(),
	})

	// create our user vote:
	newProjectRequestVote := vote.SDKFunc.Create(vote.CreateParams{
		Request:    newProjectRequest,
		Voter:      genIns.User(),
		IsApproved: true,
	})

	// create the request service:
	requestService := request.SDKFunc.CreateSDKService(request.CreateSDKServiceParams{
		PK:     pk,
		Client: client,
	})

	// save the request:
	saveProjectErr := requestService.Save(newProjectRequest, projectRepresentation)
	if saveProjectErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveProjectErr.Error())
		return
	}

	// create the vote service:
	voteService := vote.SDKFunc.CreateSDKService(vote.CreateSDKServiceParams{
		PK:              pk,
		Client:          client,
		CreateRouteFunc: createTokenDeveloperVoteRouteFunc(),
	})

	// save the vote, it should returns an error:
	saveProjectVoteErr := voteService.Save(newProjectRequestVote, projectRepresentation)
	if saveProjectVoteErr == nil {
		t.Errorf("the returned error was expected to be valid, nil returned")
	}
}

func TestSaveGenesis_createMilestone_voteOnMilestone_Success(t *testing.T) {
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

	dev := developer.SDKFunc.Create(developer.CreateParams{
		Pledge: pldge,
		User:   genIns.User(),
		Name:   "Steve",
		Resume: "this is a resume",
	})

	proj := project.SDKFunc.Create(project.CreateParams{
		Title:       "This is a project",
		Description: "This is the project description",
	})

	mils := milestone.SDKFunc.Create(milestone.CreateParams{
		Project:     proj,
		Title:       "This is a milestone",
		Description: "This is a milestone description",
		CreatedOn:   time.Now().UTC(),
		DueOn:       time.Now().Add(time.Second * 60 * 60 * 24),
	})

	rootPath := filepath.Join("./test_files_TestSaveGenesis_createMilestone_voteOnMilestone_Success")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, client, service, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, genIns)
	defer node.Stop()

	// save the milestone:
	saveMilestone(t, client, pk, service, repository, genIns.User(), userIns, pldge, dev, proj, mils)
}

func TestSaveGenesis_createTask_voteOnTask_Success(t *testing.T) {
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

	dev := developer.SDKFunc.Create(developer.CreateParams{
		Pledge: pldge,
		User:   genIns.User(),
		Name:   "Steve",
		Resume: "this is a resume",
	})

	proj := project.SDKFunc.Create(project.CreateParams{
		Title:       "This is a project",
		Description: "This is the project description",
	})

	mils := milestone.SDKFunc.Create(milestone.CreateParams{
		Project:     proj,
		Title:       "This is a milestone",
		Description: "This is a milestone description",
		CreatedOn:   time.Now().UTC(),
		DueOn:       time.Now().Add(time.Second * 60 * 60 * 24),
	})

	tsk := task.SDKFunc.Create(task.CreateParams{
		Milestone:   mils,
		Creator:     dev,
		Title:       "This is a task title",
		Description: "This is a task description",
		CreatedOn:   time.Now().Add(time.Second * 60 * 60 * 1),
		DueOn:       time.Now().Add(time.Second * 60 * 60 * 3),
	})

	rootPath := filepath.Join("./test_files_TestSaveGenesis_createTask_voteOnTask_Success")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, client, service, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, genIns)
	defer node.Stop()

	// save the task:
	saveTask(t, client, pk, service, repository, genIns.User(), userIns, pldge, dev, proj, mils, tsk)
}
