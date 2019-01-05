package balance

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token"
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

// RetrieveByWalletAndToken retrieves a Balance instance by wallet and token
func (app *repository) RetrieveByWalletAndToken(wal wallet.Wallet, tok token.Token) (Balance, error) {
	// retrieve all the withdrawals related to our wallet and token:
	withs, withsErr := app.withdrawalRepository.RetrieveSetByFromWalletAndToken(wal, tok)
	if withsErr != nil {
		return nil, withsErr
	}

	// retrieve all the deposits related to our wallet and token:
	deps, depsErr := app.depositRepository.RetrieveSetByToWalletAndToken(wal, tok)
	if depsErr != nil {
		return nil, depsErr
	}

	return app.calculate(wal, tok, withs, deps)
}

func (app *repository) calculate(
	wal wallet.Wallet,
	tok token.Token,
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
	bal := createBalance(wal, tok, total)
	return bal, nil
}
