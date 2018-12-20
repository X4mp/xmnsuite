package link

import (
	"log"
	"time"

	"github.com/montanaflynn/stats"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/validator"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/link"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/node"
	"github.com/xmnservices/xmnsuite/blockchains/tendermint"
	"github.com/xmnservices/xmnsuite/crypto"
)

type application struct {
	pk                                crypto.PrivateKey
	linkAmountToRetrievePerBatch      int
	validatorAmountToRetrievePerBatch int
	sleepAfterUpdateDuration          time.Duration
	entityService                     entity.Service
	linkRepository                    link.Repository
	nodeRepository                    node.Repository
	nodeRepresentation                entity.Representation
	stop                              bool
}

func createApplication(
	pk crypto.PrivateKey,
	linkAmountToRetrievePerBatch int,
	validatorAmountToRetrievePerBatch int,
	sleepAfterUpdateDuration time.Duration,
	entityService entity.Service,
	linkRepository link.Repository,
	nodeRepository node.Repository,
	nodeRepresentation entity.Representation,
) Daemon {
	out := application{
		pk: pk,
		linkAmountToRetrievePerBatch:      linkAmountToRetrievePerBatch,
		validatorAmountToRetrievePerBatch: validatorAmountToRetrievePerBatch,
		sleepAfterUpdateDuration:          sleepAfterUpdateDuration,
		entityService:                     entityService,
		linkRepository:                    linkRepository,
		nodeRepository:                    nodeRepository,
		nodeRepresentation:                nodeRepresentation,
	}

	return &out
}

// Start starts the link daemon
func (app *application) Start() error {

	app.stop = false

	for {

		// sleep some time:
		log.Printf("Waiting %f seconds...", app.sleepAfterUpdateDuration.Seconds())
		time.Sleep(app.sleepAfterUpdateDuration)

		// if we must stop:
		if app.stop {
			return nil
		}

		// retrieve the current links from the database:
		index := 0
		retPartialSet, retPartialSetErr := app.linkRepository.RetrieveSet(index, app.linkAmountToRetrievePerBatch)
		if retPartialSetErr != nil {
			log.Printf("there was an error while retrieving Link instances (index: %d, amount: %d): %s", index, app.linkAmountToRetrievePerBatch, retPartialSetErr.Error())
			continue
		}

		// for each link, download the nodes, on all nodes, and create the real node list based on the power of everyone:
		lnks := retPartialSet.Instances()
		for _, oneLinkIns := range lnks {
			if lnk, ok := oneLinkIns.(link.Link); ok {
				// retrieve the nodes related to the link:
				nodes, nodesErr := app.nodeRepository.RetrieveByLink(lnk)
				if nodesErr != nil {
					log.Printf("there was an error while retrieving Node instances related to Link (ID: %s): %s", lnk.ID().String(), nodesErr.Error())
					continue
				}

				// if there is no node, continue:
				if len(nodes) <= 0 {
					log.Printf("the link (ID: %s) contain no nodes", lnk.ID().String())
					continue
				}

				// retrieve the validators:
				validators := app.retrieveLinkValidators(nodes)

				// convert the fetched validators to nodes:
				newNodes, newNodesErr := app.convert(lnk, validators)
				if newNodesErr != nil {
					log.Printf("there was an error while converting fetched Validator instances to Node instances for Link (ID: %s): %s", lnk.ID().String(), newNodesErr.Error())
					continue
				}

				// update the link nodes in the database:
				updateErr := app.updateDB(lnk, nodes, newNodes)
				if updateErr != nil {
					log.Printf("there was an error while updating nodes on Link (ID: %s): %s", lnk.ID().String(), updateErr.Error())
					continue
				}
			}

			// log
			log.Printf("the entity (ID: %s) was expected to be a Link instance", oneLinkIns.ID().String())
		}

	}
}

// Stop stops the link daemon
func (app *application) Stop() error {
	app.stop = true
	return nil
}

func (app *application) updateDB(lnk link.Link, prevNodes []node.Node, newNodes []node.Node) error {
	// delete the old nodes:
	for _, oneNode := range prevNodes {
		delNodeErr := app.entityService.Delete(oneNode, app.nodeRepresentation)
		if delNodeErr != nil {
			log.Printf("there was an error while deleting an old Node (ID: %s) on Link (ID: %s): %s", oneNode.ID().String(), lnk.ID().String(), delNodeErr.Error())
		}
	}

	// save the new nodes:
	for _, oneNode := range newNodes {
		saveNodeErr := app.entityService.Save(oneNode, app.nodeRepresentation)
		if saveNodeErr != nil {
			log.Printf("there was an error while saving a new Node (ID: %s) on Link (ID: %s): %s", oneNode.ID().String(), lnk.ID().String(), saveNodeErr.Error())
		}
	}

	return nil
}

func (app *application) convert(lnk link.Link, valMap map[string][]validator.Validator) ([]node.Node, error) {
	// we make the per-node map:
	vals := map[string][]validator.Validator{}
	for _, oneList := range valMap {
		for _, oneVal := range oneList {
			nodeIDAsString := oneVal.ID().String()
			vals[nodeIDAsString] = append(vals[nodeIDAsString], oneVal)
		}
	}

	// for each validator, re-create the validator instances with a median of its powers:
	out := []node.Node{}
	for _, oneValList := range vals {
		pows := []float64{}
		for _, oneVal := range oneValList {
			pows = append(pows, float64(oneVal.Pledge().From().Amount()))
		}

		med, medErr := stats.Median(pows)
		if medErr != nil {
			return nil, medErr
		}

		rounded, roundedErr := stats.Round(med, 0)
		if roundedErr != nil {
			return nil, roundedErr
		}

		// create the validator:
		oneValidator := oneValList[:1][0]
		nod := node.SDKFunc.Create(node.CreateParams{
			ID:    oneValidator.ID(),
			Link:  lnk,
			Power: int(rounded),
			IP:    oneValidator.IP(),
			Port:  oneValidator.Port(),
		})

		out = append(out, nod)
	}

	// re-order the list:
	return out, nil
}

func (app *application) retrieveLinkValidators(nods []node.Node) map[string][]validator.Validator {
	// retrieve the node validators:
	validators := map[string][]validator.Validator{}
	for _, oneNode := range nods {
		// create the validator repository:
		validatorRepository := validator.SDKFunc.CreateSDKRepository(entity.SDKFunc.CreateSDKRepository(entity.CreateSDKRepositoryParams{
			PK: app.pk,
			Client: tendermint.SDKFunc.CreateClient(tendermint.CreateClientParams{
				IP:   oneNode.IP(),
				Port: oneNode.Port(),
			}),
		}))

		// retrieve the validators:
		valIndex := 0
		nodeValidators := []validator.Validator{}
		for {
			valPS, valPSErr := validatorRepository.RetrieveSet(valIndex, app.validatorAmountToRetrievePerBatch)
			if valPSErr != nil {
				log.Printf("there was an error while retrieving %d Validator instances for Node (ID: %s): %s", app.validatorAmountToRetrievePerBatch, oneNode.ID().String(), valPSErr.Error())
				break
			}

			if valPS.IsLast() {
				break
			}

			valsIns := valPS.Instances()
			for _, oneValidatorIns := range valsIns {
				if val, ok := oneValidatorIns.(validator.Validator); ok {
					nodeValidators = append(nodeValidators, val)
				}
			}
		}

		// add the validators to the map:
		validators[oneNode.ID().String()] = nodeValidators

	}

	return validators
}
