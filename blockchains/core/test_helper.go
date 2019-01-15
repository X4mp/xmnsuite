package core

import (
	"math/rand"
	"net"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/core/meta"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/genesis"
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
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request"
	active_request "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active/vote"
	active_vote "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active/vote/active"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/group"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/keyname"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/category"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/link"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/node"
	community_project "github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/project"
	"github.com/xmnservices/xmnsuite/crypto"
)

type simpleRequestVote struct {
	Voter      user.User
	IsApproved bool
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
	met := meta.SDKFunc.Create(meta.CreateParams{})
	node, nodeErr := saveThenSpawnBlockchain(namespace, name, &id, nil, rootPath, port, nodePK, pk.PublicKey(), met)
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

	// compare the genesis instances:
	genesis.CompareGenesisForTests(t, genIns, retGen.(genesis.Genesis))

	// retrieve the genesis by intersect keynames:
	retGenByIntersectKeynames, retGenByIntersectKeynamesErr := repository.RetrieveByIntersectKeynames(representation.MetaData(), []string{"genesis"})
	if retGenByIntersectKeynamesErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retGenByIntersectKeynamesErr.Error())
		return nil, nil, nil, nil
	}

	// compare the genesis instances:
	genesis.CompareGenesisForTests(t, genIns, retGenByIntersectKeynames.(genesis.Genesis))

	// retrieve the genesis partial set by keyname:
	retGenSetByIntersectKeynames, retGenSetByIntersectKeynamesErr := repository.RetrieveSetByKeyname(representation.MetaData(), "genesis", 0, 5)
	if retGenSetByIntersectKeynamesErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retGenSetByIntersectKeynamesErr.Error())
		return nil, nil, nil, nil
	}

	if retGenSetByIntersectKeynames.Index() != 0 {
		t.Errorf("the index was invalid.  Expected: %d, Received: %d", 0, retGenSetByIntersectKeynames.Index())
		return nil, nil, nil, nil
	}

	if retGenSetByIntersectKeynames.Amount() != 1 {
		t.Errorf("the amount was invalid.  Expected: %d, Received: %d", 1, retGenSetByIntersectKeynames.Amount())
		return nil, nil, nil, nil
	}

	if retGenSetByIntersectKeynames.TotalAmount() != 1 {
		t.Errorf("the total amount was invalid.  Expected: %d, Received: %d", 1, retGenSetByIntersectKeynames.TotalAmount())
		return nil, nil, nil, nil
	}

	// compare the genesis instances:
	genInstances := retGenSetByIntersectKeynames.Instances()
	genesis.CompareGenesisForTests(t, genIns, genInstances[0].(genesis.Genesis))

	// save the request group and keyname:
	meta := meta.SDKFunc.Create(meta.CreateParams{})
	entReqs := meta.WriteOnEntityRequest()
	for _, entReq := range entReqs {
		grp := group.SDKFunc.Create(group.CreateParams{
			Name: entReq.RequestedBy().MetaData().Keyname(),
		})

		mp := entReq.Map()
		keynameRepresentation := keyname.SDKFunc.CreateRepresentation()
		for _, oneRepresentation := range mp {
			kname := keyname.SDKFunc.Create(keyname.CreateParams{
				Name:  oneRepresentation.MetaData().Keyname(),
				Group: grp,
			})

			// save the keyname:
			saveKeynameErr := service.Save(kname, keynameRepresentation)
			if saveKeynameErr != nil {
				t.Errorf("the returned error was expected to be nil, error returned: %s", saveKeynameErr.Error())
				return nil, nil, nil, nil
			}

			// retrieve the keyname:
			retKeyname, retKeynameErr := repository.RetrieveByID(keynameRepresentation.MetaData(), kname.ID())
			if retKeynameErr != nil {
				t.Errorf("the returned error was expected to be nil, error returned: %s", retKeynameErr.Error())
				return nil, nil, nil, nil
			}

			// compare the keyname instances:
			keyname.CompareKeynameForTests(t, kname, retKeyname.(keyname.Keyname))
		}

	}

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

func saveRequestThenSaveVotesForTests(
	t *testing.T,
	client applications.Client,
	pk crypto.PrivateKey,
	repository entity.Repository,
	representation entity.Representation,
	req request.Request,
	votesPK []crypto.PrivateKey,
	reqVotes []*simpleRequestVote,
) request.Service {
	// create the metadata:
	requestMetaData := active_request.SDKFunc.CreateMetaData()

	// create the request service:
	requestService := request.SDKFunc.CreateSDKService(request.CreateSDKServiceParams{
		PK:     pk,
		Client: client,
	})

	// create the request repository:
	requestRepository := active_request.SDKFunc.CreateRepository(active_request.CreateRepositoryParams{
		EntityRepository: repository,
	})

	voteRepository := active_vote.SDKFunc.CreateRepository(active_vote.CreateRepositoryParams{
		EntityRepository: repository,
	})

	// save the request:
	saveRequestErr := requestService.Save(req, representation)
	if saveRequestErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveRequestErr.Error())
		return nil
	}

	// retrieve the request and compare them:
	retRequest, retRequesterr := requestRepository.RetrieveByRequest(req)
	if retRequesterr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retRequesterr.Error())
		return nil
	}

	// compare the requests:
	request.CompareRequestForTests(t, req, retRequest.Request().(request.Request))

	// save the votes:
	crVotes := []vote.Vote{}
	for index, oneVote := range reqVotes {
		// create the vote service:
		oneVoteService := vote.SDKFunc.CreateSDKService(vote.CreateSDKServiceParams{
			PK:     votesPK[index],
			Client: client,
		})

		// create the vote:
		vot := vote.SDKFunc.Create(vote.CreateParams{
			Request:    retRequest,
			Voter:      oneVote.Voter,
			Reason:     "TEST",
			IsApproved: oneVote.IsApproved,
			IsNeutral:  false,
		})

		// save the vote:
		savedVoteErr := oneVoteService.Save(vot, representation)
		if savedVoteErr != nil {
			t.Errorf("the returned error was expected to be nil, error returned: %s", savedVoteErr.Error())
			return nil
		}

		if (index + 1) < len(reqVotes) {
			retVote, retVoteErr := voteRepository.RetrieveByVote(vot)
			if retVoteErr != nil {
				t.Errorf("the returned error was expected to be nil, error returned: %s", retVoteErr.Error())
				return nil
			}

			// compare the votes:
			vote.CompareVoteForTests(t, vot, retVote.(active_vote.Vote).Vote())
		}

		crVotes = append(crVotes, vot)
	}

	// the request should no longer exists:
	_, retRequestErr := repository.RetrieveByID(requestMetaData, req.ID())
	if retRequestErr == nil {
		t.Errorf("the returned error was expected to be valid, nil returned")
		return nil
	}

	// the votes should no longer exists:
	for _, oneVote := range crVotes {
		_, retVoteErr := voteRepository.RetrieveByVote(oneVote)
		if retVoteErr == nil {
			t.Errorf("the returned error was expected to be valid, nil returned")
			return nil
		}
	}

	// the new entity should now be an entity:
	_, retInsErr := repository.RetrieveByID(representation.MetaData(), req.Save().ID())
	if retInsErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retInsErr.Error())
		return nil
	}

	return requestService
}

func saveUserWithNewWallet(
	t *testing.T,
	client applications.Client,
	pk crypto.PrivateKey,
	service entity.Service,
	repository entity.Repository,
	fromUser user.User,
	newUser user.User,
) request.Service {

	// create the representations:
	userRepresentation := user.SDKFunc.CreateRepresentation()
	metaData := userRepresentation.MetaData()

	// retrieve the keyname:
	knameRepository := keyname.SDKFunc.CreateRepository(keyname.CreateRepositoryParams{
		EntityRepository: repository,
	})

	kname, knameErr := knameRepository.RetrieveByName(metaData.Keyname())
	if knameErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", knameErr.Error())
	}

	// create the wallet request:
	newUserRequest := request.SDKFunc.Create(request.CreateParams{
		FromUser:   fromUser,
		SaveEntity: newUser,
		Reason:     "TEST",
		Keyname:    kname,
	})

	// create our user vote:
	newUserRequestVote := &simpleRequestVote{
		Voter:      fromUser,
		IsApproved: true,
	}

	// save the new token request, then save vote:
	requestService := saveRequestThenSaveVotesForTests(t, client, pk, repository, userRepresentation, newUserRequest, []crypto.PrivateKey{pk}, []*simpleRequestVote{
		newUserRequestVote,
	})

	// returns:
	return requestService
}

func saveLink(
	t *testing.T,
	client applications.Client,
	pk crypto.PrivateKey,
	service entity.Service,
	repository entity.Repository,
	fromUser user.User,
	lnk link.Link,
) request.Service {

	// create the representations:
	linkRepresentation := link.SDKFunc.CreateRepresentation()

	// retrieve the keyname:
	knameRepository := keyname.SDKFunc.CreateRepository(keyname.CreateRepositoryParams{
		EntityRepository: repository,
	})

	kname, knameErr := knameRepository.RetrieveByName(link.SDKFunc.CreateMetaData().Keyname())
	if knameErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", knameErr.Error())
	}

	// create the link request:
	newLinkRequest := request.SDKFunc.Create(request.CreateParams{
		FromUser:   fromUser,
		SaveEntity: lnk,
		Reason:     "TEST",
		Keyname:    kname,
	})

	// create our user vote:
	newLinkRequestVote := &simpleRequestVote{
		Voter:      fromUser,
		IsApproved: true,
	}

	// save the new token request, then save vote:
	requestService := saveRequestThenSaveVotesForTests(t, client, pk, repository, linkRepresentation, newLinkRequest, []crypto.PrivateKey{pk}, []*simpleRequestVote{
		newLinkRequestVote,
	})

	// returns:
	return requestService
}

func saveNode(
	t *testing.T,
	client applications.Client,
	pk crypto.PrivateKey,
	service entity.Service,
	repository entity.Repository,
	fromUser user.User,
	lnk link.Link,
	nod node.Node,
) request.Service {

	// create the representations:
	nodeRepresentation := node.SDKFunc.CreateRepresentation()

	// save the link:
	saveLink(t, client, pk, service, repository, fromUser, lnk)

	// retrieve the keyname:
	knameRepository := keyname.SDKFunc.CreateRepository(keyname.CreateRepositoryParams{
		EntityRepository: repository,
	})

	kname, knameErr := knameRepository.RetrieveByName(node.SDKFunc.CreateMetaData().Keyname())
	if knameErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", knameErr.Error())
	}

	// create the node request:
	newNodeRequest := request.SDKFunc.Create(request.CreateParams{
		FromUser:   fromUser,
		SaveEntity: nod,
		Reason:     "TEST",
		Keyname:    kname,
	})

	// create our user vote:
	newNodeRequestVote := &simpleRequestVote{
		Voter:      fromUser,
		IsApproved: true,
	}

	// save the new token request, then save vote:
	requestService := saveRequestThenSaveVotesForTests(t, client, pk, repository, nodeRepresentation, newNodeRequest, []crypto.PrivateKey{pk}, []*simpleRequestVote{
		newNodeRequestVote,
	})

	// returns:
	return requestService
}

func savePledge(
	t *testing.T,
	client applications.Client,
	pk crypto.PrivateKey,
	service entity.Service,
	repository entity.Repository,
	fromGen genesis.Genesis,
	pldge pledge.Pledge,
) request.Service {

	// create the repreentations:
	pldgeRepresentation := pledge.SDKFunc.CreateRepresentation()

	// retrieve the keyname:
	knameRepository := keyname.SDKFunc.CreateRepository(keyname.CreateRepositoryParams{
		EntityRepository: repository,
	})

	kname, knameErr := knameRepository.RetrieveByName(pledge.SDKFunc.CreateMetaData().Keyname())
	if knameErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", knameErr.Error())
	}

	// create the user in wallet request:
	pldgeRequest := request.SDKFunc.Create(request.CreateParams{
		FromUser:   fromGen.User(),
		SaveEntity: pldge,
		Reason:     "TEST",
		Keyname:    kname,
	})

	// create our user vote:
	pldgeRequestVote := &simpleRequestVote{
		Voter:      fromGen.User(),
		IsApproved: true,
	}

	// save the new wallet request, then save vote:
	requestService := saveRequestThenSaveVotesForTests(t, client, pk, repository, pldgeRepresentation, pldgeRequest, []crypto.PrivateKey{pk}, []*simpleRequestVote{
		pldgeRequestVote,
	})

	// returns:
	return requestService
}

func saveCategory(
	t *testing.T,
	client applications.Client,
	pk crypto.PrivateKey,
	service entity.Service,
	repository entity.Repository,
	fromGen genesis.Genesis,
	cat category.Category,
) request.Service {

	// create the representation:
	categoryRepresentation := category.SDKFunc.CreateRepresentation()

	// retrieve the keyname:
	knameRepository := keyname.SDKFunc.CreateRepository(keyname.CreateRepositoryParams{
		EntityRepository: repository,
	})

	kname, knameErr := knameRepository.RetrieveByName(category.SDKFunc.CreateMetaData().Keyname())
	if knameErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", knameErr.Error())
	}

	// create the user in wallet request:
	catRequest := request.SDKFunc.Create(request.CreateParams{
		FromUser:   fromGen.User(),
		SaveEntity: cat,
		Reason:     "TEST",
		Keyname:    kname,
	})

	// create our user vote:
	catRequestVote := &simpleRequestVote{
		Voter:      fromGen.User(),
		IsApproved: true,
	}

	// save the new wallet request, then save vote:
	requestService := saveRequestThenSaveVotesForTests(t, client, pk, repository, categoryRepresentation, catRequest, []crypto.PrivateKey{pk}, []*simpleRequestVote{
		catRequestVote,
	})

	// returns:
	return requestService
}

func saveProposal(
	t *testing.T,
	client applications.Client,
	pk crypto.PrivateKey,
	service entity.Service,
	repository entity.Repository,
	fromGen genesis.Genesis,
	prop proposal.Proposal,
) request.Service {

	// create the representation:
	proposalRepresentation := proposal.SDKFunc.CreateRepresentation()

	// retrieve the keyname:
	knameRepository := keyname.SDKFunc.CreateRepository(keyname.CreateRepositoryParams{
		EntityRepository: repository,
	})

	kname, knameErr := knameRepository.RetrieveByName(proposal.SDKFunc.CreateMetaData().Keyname())
	if knameErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", knameErr.Error())
	}

	// create the user in wallet request:
	proposalRequest := request.SDKFunc.Create(request.CreateParams{
		FromUser:   fromGen.User(),
		SaveEntity: prop,
		Reason:     "TEST",
		Keyname:    kname,
	})

	// create our user vote:
	proposalRequestVote := &simpleRequestVote{
		Voter:      fromGen.User(),
		IsApproved: true,
	}

	// save the new wallet request, then save vote:
	requestService := saveRequestThenSaveVotesForTests(t, client, pk, repository, proposalRepresentation, proposalRequest, []crypto.PrivateKey{pk}, []*simpleRequestVote{
		proposalRequestVote,
	})

	// returns:
	return requestService
}

func saveCommunityProject(
	t *testing.T,
	client applications.Client,
	pk crypto.PrivateKey,
	service entity.Service,
	repository entity.Repository,
	fromGen genesis.Genesis,
	proj community_project.Project,
) request.Service {

	// create the representation:
	projectRepresentation := community_project.SDKFunc.CreateRepresentation()

	// retrieve the keyname:
	knameRepository := keyname.SDKFunc.CreateRepository(keyname.CreateRepositoryParams{
		EntityRepository: repository,
	})

	kname, knameErr := knameRepository.RetrieveByName(community_project.SDKFunc.CreateMetaData().Keyname())
	if knameErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", knameErr.Error())
	}

	// create the user in wallet request:
	projRequest := request.SDKFunc.Create(request.CreateParams{
		FromUser:   fromGen.User(),
		SaveEntity: proj,
		Reason:     "TEST",
		Keyname:    kname,
	})

	// create our user vote:
	projRequestVote := &simpleRequestVote{
		Voter:      fromGen.User(),
		IsApproved: true,
	}

	// save the new wallet request, then save vote:
	requestService := saveRequestThenSaveVotesForTests(t, client, pk, repository, projectRepresentation, projRequest, []crypto.PrivateKey{pk}, []*simpleRequestVote{
		projRequestVote,
	})

	// returns:
	return requestService
}

func saveProject(
	t *testing.T,
	client applications.Client,
	pk crypto.PrivateKey,
	service entity.Service,
	repository entity.Repository,
	fromGen genesis.Genesis,
	proj project.Project,
) request.Service {

	// create the representation:
	projectRepresentation := project.SDKFunc.CreateRepresentation()

	// retrieve the keyname:
	knameRepository := keyname.SDKFunc.CreateRepository(keyname.CreateRepositoryParams{
		EntityRepository: repository,
	})

	kname, knameErr := knameRepository.RetrieveByName(project.SDKFunc.CreateMetaData().Keyname())
	if knameErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", knameErr.Error())
	}

	// create the user in wallet request:
	projRequest := request.SDKFunc.Create(request.CreateParams{
		FromUser:   fromGen.User(),
		SaveEntity: proj,
		Reason:     "TEST",
		Keyname:    kname,
	})

	// create our user vote:
	projRequestVote := &simpleRequestVote{
		Voter:      fromGen.User(),
		IsApproved: true,
	}

	// save the new wallet request, then save vote:
	requestService := saveRequestThenSaveVotesForTests(t, client, pk, repository, projectRepresentation, projRequest, []crypto.PrivateKey{pk}, []*simpleRequestVote{
		projRequestVote,
	})

	// returns:
	return requestService
}

func saveMilestone(
	t *testing.T,
	client applications.Client,
	pk crypto.PrivateKey,
	service entity.Service,
	repository entity.Repository,
	fromGen genesis.Genesis,
	mils milestone.Milestone,
) request.Service {

	// create the representation:
	milestoneRepresentation := milestone.SDKFunc.CreateRepresentation()

	// retrieve the keyname:
	knameRepository := keyname.SDKFunc.CreateRepository(keyname.CreateRepositoryParams{
		EntityRepository: repository,
	})

	kname, knameErr := knameRepository.RetrieveByName(milestone.SDKFunc.CreateMetaData().Keyname())
	if knameErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", knameErr.Error())
	}

	// create the user in wallet request:
	milsRequest := request.SDKFunc.Create(request.CreateParams{
		FromUser:   fromGen.User(),
		SaveEntity: mils,
		Reason:     "TEST",
		Keyname:    kname,
	})

	// create our user vote:
	milsRequestVote := &simpleRequestVote{
		Voter:      fromGen.User(),
		IsApproved: true,
	}

	// save the new wallet request, then save vote:
	requestService := saveRequestThenSaveVotesForTests(t, client, pk, repository, milestoneRepresentation, milsRequest, []crypto.PrivateKey{pk}, []*simpleRequestVote{
		milsRequestVote,
	})

	// returns:
	return requestService
}

func saveFeature(
	t *testing.T,
	client applications.Client,
	pk crypto.PrivateKey,
	service entity.Service,
	repository entity.Repository,
	fromGen genesis.Genesis,
	feat feature.Feature,
) request.Service {

	// create the representation:
	featureRepresentation := feature.SDKFunc.CreateRepresentation()

	// retrieve the keyname:
	knameRepository := keyname.SDKFunc.CreateRepository(keyname.CreateRepositoryParams{
		EntityRepository: repository,
	})

	kname, knameErr := knameRepository.RetrieveByName(feature.SDKFunc.CreateMetaData().Keyname())
	if knameErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", knameErr.Error())
	}

	// create the user in wallet request:
	featureRequest := request.SDKFunc.Create(request.CreateParams{
		FromUser:   fromGen.User(),
		SaveEntity: feat,
		Reason:     "TEST",
		Keyname:    kname,
	})

	// create our user vote:
	featureRequestVote := &simpleRequestVote{
		Voter:      fromGen.User(),
		IsApproved: true,
	}

	// save the new wallet request, then save vote:
	requestService := saveRequestThenSaveVotesForTests(t, client, pk, repository, featureRepresentation, featureRequest, []crypto.PrivateKey{pk}, []*simpleRequestVote{
		featureRequestVote,
	})

	// returns:
	return requestService
}

func saveTask(
	t *testing.T,
	client applications.Client,
	pk crypto.PrivateKey,
	service entity.Service,
	repository entity.Repository,
	fromGen genesis.Genesis,
	tsk task.Task,
) request.Service {

	// create the representation:
	taskRepresentation := task.SDKFunc.CreateRepresentation()

	// retrieve the keyname:
	knameRepository := keyname.SDKFunc.CreateRepository(keyname.CreateRepositoryParams{
		EntityRepository: repository,
	})

	kname, knameErr := knameRepository.RetrieveByName(task.SDKFunc.CreateMetaData().Keyname())
	if knameErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", knameErr.Error())
	}

	// create the user in wallet request:
	taskRequest := request.SDKFunc.Create(request.CreateParams{
		FromUser:   fromGen.User(),
		SaveEntity: tsk,
		Reason:     "TEST",
		Keyname:    kname,
	})

	// create our user vote:
	taskRequestVote := &simpleRequestVote{
		Voter:      fromGen.User(),
		IsApproved: true,
	}

	// save the new wallet request, then save vote:
	requestService := saveRequestThenSaveVotesForTests(t, client, pk, repository, taskRepresentation, taskRequest, []crypto.PrivateKey{pk}, []*simpleRequestVote{
		taskRequestVote,
	})

	// returns:
	return requestService
}

func savePledgeTask(
	t *testing.T,
	client applications.Client,
	pk crypto.PrivateKey,
	service entity.Service,
	repository entity.Repository,
	fromGen genesis.Genesis,
	tsk pledge_task.Task,
) request.Service {

	// create the representation:
	taskRepresentation := pledge_task.SDKFunc.CreateRepresentation()

	// retrieve the keyname:
	knameRepository := keyname.SDKFunc.CreateRepository(keyname.CreateRepositoryParams{
		EntityRepository: repository,
	})

	kname, knameErr := knameRepository.RetrieveByName(pledge_task.SDKFunc.CreateMetaData().Keyname())
	if knameErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", knameErr.Error())
	}

	// create the user in wallet request:
	taskRequest := request.SDKFunc.Create(request.CreateParams{
		FromUser:   fromGen.User(),
		SaveEntity: tsk,
		Reason:     "TEST",
		Keyname:    kname,
	})

	// create our user vote:
	taskRequestVote := &simpleRequestVote{
		Voter:      fromGen.User(),
		IsApproved: true,
	}

	// save the new wallet request, then save vote:
	requestService := saveRequestThenSaveVotesForTests(t, client, pk, repository, taskRepresentation, taskRequest, []crypto.PrivateKey{pk}, []*simpleRequestVote{
		taskRequestVote,
	})

	// returns:
	return requestService
}

func saveCompletedTask(
	t *testing.T,
	client applications.Client,
	pk crypto.PrivateKey,
	service entity.Service,
	repository entity.Repository,
	fromGen genesis.Genesis,
	tsk completed_task.Task,
) request.Service {

	// create the representation:
	taskRepresentation := completed_task.SDKFunc.CreateRepresentation()

	// retrieve the keyname:
	knameRepository := keyname.SDKFunc.CreateRepository(keyname.CreateRepositoryParams{
		EntityRepository: repository,
	})

	kname, knameErr := knameRepository.RetrieveByName(completed_task.SDKFunc.CreateMetaData().Keyname())
	if knameErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", knameErr.Error())
	}

	// create the user in wallet request:
	taskRequest := request.SDKFunc.Create(request.CreateParams{
		FromUser:   fromGen.User(),
		SaveEntity: tsk,
		Reason:     "TEST",
		Keyname:    kname,
	})

	// create our user vote:
	taskRequestVote := &simpleRequestVote{
		Voter:      fromGen.User(),
		IsApproved: true,
	}

	// save the new wallet request, then save vote:
	requestService := saveRequestThenSaveVotesForTests(t, client, pk, repository, taskRepresentation, taskRequest, []crypto.PrivateKey{pk}, []*simpleRequestVote{
		taskRequestVote,
	})

	// returns:
	return requestService
}

func saveTransfer(
	t *testing.T,
	client applications.Client,
	pk crypto.PrivateKey,
	service entity.Service,
	repository entity.Repository,
	fromGen genesis.Genesis,
	trsf transfer.Transfer,
) request.Service {

	// create the representation:
	transferRepresentation := transfer.SDKFunc.CreateRepresentation()

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
		FromUser:   fromGen.User(),
		SaveEntity: trsf,
		Reason:     "TEST",
		Keyname:    kname,
	})

	// create our user vote:
	trsfRequestVote := &simpleRequestVote{
		Voter:      fromGen.User(),
		IsApproved: true,
	}

	// save the new wallet request, then save vote:
	requestService := saveRequestThenSaveVotesForTests(t, client, pk, repository, transferRepresentation, trsfRequest, []crypto.PrivateKey{pk}, []*simpleRequestVote{
		trsfRequestVote,
	})

	// returns:
	return requestService
}
