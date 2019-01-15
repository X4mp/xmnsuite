package balance

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/withdrawal"
)

type repository struct {
	depositRepository    deposit.Repository
	withdrawalRepository withdrawal.Repository
}

func createRepository(depositRepository deposit.Repository, withdrawalRepository withdrawal.Repository) Repository {
	out := repository{
		depositRepository:    depositRepository,
		withdrawalRepository: withdrawalRepository,
	}
	return &out
}

// RetrieveByWallet retrieves a Balance instance by wallet
func (app *repository) RetrieveByWallet(wal wallet.Wallet) (Balance, error) {
	// retrieve all the withdrawals related to our wallet:
	withs, withsErr := app.withdrawalRepository.RetrieveSetByFromWallet(wal)
	if withsErr != nil {
		return nil, withsErr
	}

	// retrieve all the deposits related to our wallet:
	deps, depsErr := app.depositRepository.RetrieveSetByToWallet(wal)
	if depsErr != nil {
		return nil, depsErr
	}

	return app.calculate(wal, withs, deps)
}

func (app *repository) calculate(
	wal wallet.Wallet,
	withs []withdrawal.Withdrawal,
	deps []deposit.Deposit,
) (Balance, error) {
	// calculate the withdrawals amount:
	withAmount := 0
	for _, oneWithdrawalIns := range withs {
		withAmount += oneWithdrawalIns.(withdrawal.Withdrawal).Amount()
	}

	// calculate the deposits amount:
	depAmount := 0
	for _, oneDepIns := range deps {
		depAmount += oneDepIns.(deposit.Deposit).Amount()
	}

	// create the balance:
	total := depAmount - withAmount
	bal := createBalance(wal, total)
	return bal, nil
}
