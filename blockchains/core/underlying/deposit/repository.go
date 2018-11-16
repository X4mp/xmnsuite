package deposit

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token"
	"github.com/xmnservices/xmnsuite/blockchains/core/wallet"
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

// RetrieveSetByToWalletAndToken retrieves a deposit partial set related to a wallet and token:
func (app *repository) RetrieveSetByToWalletAndToken(wal wallet.Wallet, tok token.Token) (entity.PartialSet, error) {
	keynames := []string{
		retrieveDepositsByToWalletIDKeyname(wal.ID()),
		retrieveDepositsByTokenIDKeyname(tok.ID()),
	}

	return app.entityRepository.RetrieveSetByIntersectKeynames(app.depositMetaData, keynames, 0, -1)
}
