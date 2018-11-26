package core

import (
	"math"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"testing"

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
	rootPath := filepath.Join("./test_files")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// create the representations:
	userRepresentation := user.SDKFunc.CreateRepresentation()
	walletRepresentation := wallet.SDKFunc.CreateRepresentation()

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
	walletRepresentation := wallet.SDKFunc.CreateRepresentation()
	userRepresentation := user.SDKFunc.CreateRepresentation()
	pldgeRepresentation := pledge.SDKFunc.CreateRepresentation()
	transferRepresentation := transfer.SDKFunc.CreateRepresentation()

	rootPath := filepath.Join("./test_files")
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
		FromUser:  genIns.User(),
		NewEntity: userIns,
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
	pldgeRequest := request.SDKFunc.Create(request.CreateParams{
		FromUser:  genIns.User(),
		NewEntity: pldge,
	})

	// create our user vote:
	pldgeRequestVote := vote.SDKFunc.Create(vote.CreateParams{
		Request:    pldgeRequest,
		Voter:      genIns.User(),
		IsApproved: true,
	})

	// save the new wallet request, then save vote:
	saveRequestThenSaveVotesForTests(t, client, pk, repository, pldgeRepresentation, pldgeRequest, []crypto.PrivateKey{pk}, []vote.Vote{
		pldgeRequestVote,
	}, createEntityVoteRouteFunc())

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
		FromUser:  userIns,
		NewEntity: trsf,
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

	// save the request, returns error
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

	rootPath := filepath.Join("./test_files")
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
		FromUser:  genIns.User(),
		NewEntity: userIns,
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
		FromUser:  genIns.User(),
		NewEntity: vldator,
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

	rootPath := filepath.Join("./test_files")
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
		FromUser:  genIns.User(),
		NewEntity: userIns,
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
		FromUser:  genIns.User(),
		NewEntity: trsf,
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
		Keyname:     "projects",
		Title:       "Projects",
		Description: "The XMN projects belongs on that blockchain",
		Node: node.SDKFunc.Create(node.CreateParams{
			Power: rand.Int() % 10,
			IP:    net.ParseIP("127.0.0.1"),
			Port:  123124,
		}),
	})

	// create the representations:
	linkRepresentation := link.SDKFunc.CreateRepresentation()

	rootPath := filepath.Join("./test_files")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// spawn bockchain with genesis instance:
	node, client, _, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, genIns)
	defer node.Stop()

	// create the link request:
	newLinkRequest := request.SDKFunc.Create(request.CreateParams{
		FromUser:  genIns.User(),
		NewEntity: lnk,
	})

	// create our user vote:
	newLinkRequestVote := vote.SDKFunc.Create(vote.CreateParams{
		Request:    newLinkRequest,
		Voter:      genIns.User(),
		IsApproved: true,
	})

	// save the new token request, then save vote:
	saveRequestThenSaveVotesForTests(t, client, pk, repository, linkRepresentation, newLinkRequest, []crypto.PrivateKey{pk}, []vote.Vote{
		newLinkRequestVote,
	}, createTokenVoteRouteFunc())
}
