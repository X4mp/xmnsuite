package core

import (
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
	reqRepresentation entity.Representation,
	voteRepresentation entity.Representation,
	req request.Request,
	reqVotes []vote.Vote,
) (applications.Node, applications.Client, entity.Service, entity.Repository, request.Service, vote.Service) {
	// spawn bockchain with genesis instance:
	node, client, service, repository := spawnBlockchainWithGenesisForTests(t, pk, rootPath, gen)

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

	// save the request:
	saveRequestErr := requestService.Save(req, reqRepresentation)
	if saveRequestErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveRequestErr.Error())
		return nil, nil, nil, nil, nil, nil
	}

	// save the votes:
	for _, oneVote := range reqVotes {
		// save the vote:
		savedVoteErr := voteService.Save(oneVote, voteRepresentation)
		if savedVoteErr != nil {
			t.Errorf("the returned error was expected to be nil, error returned: %s", savedVoteErr.Error())
			return nil, nil, nil, nil, nil, nil
		}
	}

	return node, client, service, repository, requestService, voteService
}
