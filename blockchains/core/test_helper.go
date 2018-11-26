package core

import (
	"fmt"
	"math/rand"
	"net"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/genesis"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/request"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/request/vote"
	"github.com/xmnservices/xmnsuite/crypto"
)

func createEntityVoteRouteFunc() vote.CreateRouteFn {
	return func(ins vote.Vote, rep entity.Representation) (string, error) {
		return fmt.Sprintf("/%s/requests/%s", rep.MetaData().Keyname(), ins.Request().ID().String()), nil
	}
}

func createTokenVoteRouteFunc() vote.CreateRouteFn {
	return func(ins vote.Vote, rep entity.Representation) (string, error) {
		return fmt.Sprintf("/token/requests/%s/%s", ins.Request().ID().String(), rep.MetaData().Keyname()), nil
	}
}

func spawnBlockchainForTests(t *testing.T, pk crypto.PrivateKey, rootPath string) (applications.Node, applications.Client, entity.Service, entity.Repository) {
	// variables:
	namespace := "xmn"
	name := "core"
	id := uuid.NewV4()
	port := rand.Int()%9000 + 1000
	nodePK := ed25519.GenPrivKey()
	ip := net.ParseIP("127.0.0.1")

	// spawn the blockchain:
	node, nodeErr := spawnBlockchain(namespace, name, &id, rootPath, port, nodePK, pk.PublicKey())
	if nodeErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", nodeErr.Error())
		return nil, nil, nil, nil
	}

	// start the node:
	node.Start()

	// get the client:
	client, clientErr := connectToBlockchain(ip, port)
	if clientErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", clientErr.Error())
		return nil, nil, nil, nil
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

	// returns:
	return node, client, entityService, entityRepository
}

func spawnBlockchainWithGenesisForTests(t *testing.T, pk crypto.PrivateKey, rootPath string, genIns genesis.Genesis) (applications.Node, applications.Client, entity.Service, entity.Repository) {

	// sopawn the blockchain:
	node, client, service, repository := spawnBlockchainForTests(t, pk, rootPath)

	// create the representation:
	representation := genesis.SDKFunc.CreateRepresentation()

	// save the genesis:
	saveErr := service.Save(genIns, representation)
	if saveErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveErr.Error())
		return nil, nil, nil, nil
	}

	// retrieve the genesis:
	retGen, retGenErr := repository.RetrieveByID(representation.MetaData(), genIns.ID())
	if retGenErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retGenErr.Error())
		return nil, nil, nil, nil
	}

	// compare the wallet instances:
	genesis.CompareGenesisForTests(t, genIns, retGen.(genesis.Genesis))

	// returns:
	return node, client, service, repository
}

func saveEntityThenRetrieveEntityByIDThenDeleteEntityByID(
	t *testing.T,
	ins entity.Entity,
	representation entity.Representation,
	service entity.Service,
	repository entity.Repository,
) entity.Entity {
	// save the first entity:
	firstSaveErr := service.Save(ins, representation)
	if firstSaveErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", firstSaveErr.Error())
		return nil
	}

	// retrieve the entity by ID:
	retInsID, retInsIDErr := repository.RetrieveByID(representation.MetaData(), ins.ID())
	if retInsIDErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retInsIDErr.Error())
		return nil
	}

	if retInsID == nil {
		t.Errorf("the returned entity was expected to be valid, nil returned")
		return nil
	}

	// delete the entity:
	delErr := service.Delete(ins, representation)
	if delErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", delErr.Error())
		return nil
	}

	// return the retrieved entity:
	return retInsID
}

func spawnBlockchainWithGenesisThenSaveRequestThenSaveVotesForTests(
	t *testing.T,
	pk crypto.PrivateKey,
	rootPath string,
	gen genesis.Genesis,
	representation entity.Representation,
	req request.Request,
	votesPK []crypto.PrivateKey,
	reqVotes []vote.Vote,
	createRouteFunc vote.CreateRouteFn,
) (applications.Node, applications.Client, entity.Service, entity.Repository, request.Service) {
	// spawn bockchain with genesis instance:
	node, client, service, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, gen)

	// save the request then save votes:
	requestService := saveRequestThenSaveVotesForTests(t, client, pk, repository, representation, req, votesPK, reqVotes, createRouteFunc)

	// return:
	return node, client, service, repository, requestService
}

func saveRequestThenSaveVotesForTests(
	t *testing.T,
	client applications.Client,
	pk crypto.PrivateKey,
	repository entity.Repository,
	representation entity.Representation,
	req request.Request,
	votesPK []crypto.PrivateKey,
	reqVotes []vote.Vote,
	createRouteFunc vote.CreateRouteFn,
) request.Service {
	// create the metadata:
	requestMetaData := request.SDKFunc.CreateMetaData()
	voteMetaData := vote.SDKFunc.CreateMetaData()

	// create the request service:
	requestService := request.SDKFunc.CreateSDKService(request.CreateSDKServiceParams{
		PK:     pk,
		Client: client,
	})

	// save the request:
	saveRequestErr := requestService.Save(req, representation)
	if saveRequestErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveRequestErr.Error())
		return nil
	}

	// retrieve the request and compare them:
	retRequest, retRequesterr := repository.RetrieveByID(requestMetaData, req.ID())
	if retRequesterr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retRequesterr.Error())
		return nil
	}

	// compare the requests:
	request.CompareRequestForTests(t, req, retRequest.(request.Request))

	// save the votes:
	for index, oneVote := range reqVotes {
		// create the vote service:
		oneVoteService := vote.SDKFunc.CreateSDKService(vote.CreateSDKServiceParams{
			PK:              votesPK[index],
			Client:          client,
			CreateRouteFunc: createRouteFunc,
		})

		// save the vote:
		savedVoteErr := oneVoteService.Save(oneVote, representation)
		if savedVoteErr != nil {
			t.Errorf("the returned error was expected to be nil, error returned: %s", savedVoteErr.Error())
			return nil
		}

		if (index + 1) < len(reqVotes) {
			retVote, retVoteErr := repository.RetrieveByID(voteMetaData, oneVote.ID())
			if retVoteErr != nil {
				t.Errorf("the returned error was expected to be nil, error returned: %s", retVoteErr.Error())
				return nil
			}

			// compare the votes:
			vote.CompareVoteForTests(t, oneVote, retVote.(vote.Vote))
		}
	}

	// the request should no longer exists:
	_, retRequestErr := repository.RetrieveByID(requestMetaData, req.ID())
	if retRequestErr == nil {
		t.Errorf("the returned error was expected to be valid, nil returned")
		return nil
	}

	// the votes should no longer exists:
	for _, oneVote := range reqVotes {
		_, retVoteErr := repository.RetrieveByID(voteMetaData, oneVote.ID())
		if retVoteErr == nil {
			t.Errorf("the returned error was expected to be valid, nil returned")
			return nil
		}
	}

	// the new entity should now be an entity:
	_, retInsErr := repository.RetrieveByID(representation.MetaData(), req.New().ID())
	if retInsErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retInsErr.Error())
		return nil
	}

	return requestService
}
