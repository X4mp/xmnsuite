package balance

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/token"
	"github.com/xmnservices/xmnsuite/blockchains/core/wallet"
)

type repository struct {
	entityRepository entity.Repository
}

func createRepository(entityRepository entity.Repository) Repository {
	out := repository{
		entityRepository: entityRepository,
	}
	return &out
}

// RetrieveByWalletAndToken retrieves a Balance instance by wallet and token
func (app *repository) RetrieveByWalletAndToken(wal wallet.Wallet, tok token.Token) (Balance, error) {
	// create the representations:

	// retrieve all the withdrawals related to our wallet and token:

	// retrieve all the deposits related to our wallet and token:

	// create the balance:
	return nil, nil
}

// RetrieveSetByWallet retrieves a Balance PartialSet instance by wallet
func (app *repository) RetrieveSetByWallet(wal wallet.Wallet) (entity.PartialSet, error) {
	return nil, nil
}

// RetrieveSetByToken retrieves a Balance PartialSet instance by token
func (app *repository) RetrieveSetByToken(tok token.Token) (entity.PartialSet, error) {
	return nil, nil
}
