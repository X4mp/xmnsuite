package balance

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
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
	withsPS, withsPSErr := app.withdrawalRepository.RetrieveSetByFromWalletAndToken(wal, tok)
	if withsPSErr != nil {
		return nil, withsPSErr
	}

	// retrieve all the deposits related to our wallet and token:
	depsPS, depsPSErr := app.depositRepository.RetrieveSetByToWalletAndToken(wal, tok)
	if depsPSErr != nil {
		return nil, depsPSErr
	}

	return app.calculate(wal, tok, withsPS, depsPS)
}

func (app *repository) calculate(
	wal wallet.Wallet,
	tok token.Token,
	withsPS entity.PartialSet,
	depsPS entity.PartialSet,
) (Balance, error) {
	// calculate the withdrawals amount:
	withAmount := 0
	withs := withsPS.Instances()
	for _, oneWithdrawalIns := range withs {
		withAmount += oneWithdrawalIns.(withdrawal.Withdrawal).Amount()
	}

	// calculate the deposits amount:
	depAmount := 0
	deps := depsPS.Instances()
	for _, oneDepIns := range deps {
		depAmount += oneDepIns.(deposit.Deposit).Amount()
	}

	// create the balance:
	total := depAmount - withAmount
	bal := createBalance(wal, tok, total)
	return bal, nil
}
