package withdrawal

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet"
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

// RetrieveSetByFromWalletAndToken retrieves a withdrawal partial set related to a wallet and token:
func (app *repository) RetrieveSetByFromWalletAndToken(wal wallet.Wallet, tok token.Token) (entity.PartialSet, error) {
	keynames := []string{
		retrieveWithdrawalsByFromWalletIDKeyname(wal.ID()),
		retrieveWithdrawalsByTokenIDKeyname(tok.ID()),
	}

	return app.entityRepository.RetrieveSetByIntersectKeynames(app.depositMetaData, keynames, 0, -1)
}
