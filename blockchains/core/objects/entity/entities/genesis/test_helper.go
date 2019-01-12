package genesis

import (
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/information"
	"github.com/xmnservices/xmnsuite/crypto"
)

// CreateGenesisWithPubKeyForTests creates a Genesis for tests
func CreateGenesisWithPubKeyForTests(pubKey crypto.PublicKey) Genesis {
	id := uuid.NewV4()
	dep := deposit.CreateDepositWithPubKeyForTests(pubKey)
	concensusNeeded := int(dep.Amount()/2) - 1
	usr := user.CreateUserWithWalletAndPublicKeyAndSharesForTests(dep.To(), pubKey, dep.To().ConcensusNeeded())
	inf := information.CreateInformationWithConcensusNeededForTests(concensusNeeded)
	tok := token.CreateTokenForTests()
	out, outErr := createGenesis(&id, inf, dep, usr, tok)
	if outErr != nil {
		panic(outErr)
	}

	return out
}

// CompareGenesisForTests compares Genesis instances for tests
func CompareGenesisForTests(t *testing.T, first Genesis, second Genesis) {
	information.CompareInformationForTests(t, first.Info(), second.Info())
	deposit.CompareDepositForTests(t, first.Deposit(), second.Deposit())
	user.CompareUserForTests(t, first.User(), second.User())
}
