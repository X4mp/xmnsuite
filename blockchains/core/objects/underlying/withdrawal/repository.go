package withdrawal

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
)

type repository struct {
	entityRepository   entity.Repository
	withdrawalMetaData entity.MetaData
}

func createRepository(entityRepository entity.Repository, withdrawalMetaData entity.MetaData) Repository {
	out := repository{
		entityRepository:   entityRepository,
		withdrawalMetaData: withdrawalMetaData,
	}

	return &out
}

// RetrieveSetByFromWalletAndToken retrieves a withdrawal set related to a wallet
func (app *repository) RetrieveSetByFromWallet(wal wallet.Wallet) ([]Withdrawal, error) {
	keynames := []string{
		retrieveAllWithdrawalsKeyname(),
		retrieveWithdrawalsByFromWalletIDKeyname(wal.ID()),
	}

	ps, psErr := app.entityRepository.RetrieveSetByIntersectKeynames(app.withdrawalMetaData, keynames, 0, -1)
	if psErr != nil {
		return nil, psErr
	}

	ins := ps.Instances()
	withs := []Withdrawal{}
	for _, oneWithIns := range ins {
		if onWith, ok := oneWithIns.(Withdrawal); ok {
			withs = append(withs, onWith)
			continue
		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid Withdrawal instance", oneWithIns.ID().String())
		return nil, errors.New(str)
	}

	return withs, nil
}
