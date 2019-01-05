package deposit

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token"
)

type repository struct {
	entityRepository entity.Repository
	depositMetaData  entity.MetaData
}

func createRepository(entityRepository entity.Repository, depositMetaData entity.MetaData) Repository {
	out := repository{
		entityRepository: entityRepository,
		depositMetaData:  depositMetaData,
	}

	return &out
}

// RetrieveSetByToWalletAndToken retrieves a deposit set related to a wallet and token:
func (app *repository) RetrieveSetByToWalletAndToken(wal wallet.Wallet, tok token.Token) ([]Deposit, error) {
	keynames := []string{
		retrieveDepositsByToWalletIDKeyname(wal.ID()),
		retrieveDepositsByTokenIDKeyname(tok.ID()),
	}

	ps, psErr := app.entityRepository.RetrieveSetByIntersectKeynames(app.depositMetaData, keynames, 0, -1)
	if psErr != nil {
		return nil, psErr
	}

	ins := ps.Instances()
	deps := []Deposit{}
	for _, oneDepIns := range ins {
		if oneDep, ok := oneDepIns.(Deposit); ok {
			deps = append(deps, oneDep)
			continue
		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid Deposit instance", oneDepIns.ID().String())
		return nil, errors.New(str)
	}

	return deps, nil
}
