package account

import (
	"errors"

	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/work"
	"github.com/xmnservices/xmnsuite/crypto"
)

// Account represents an account
type Account interface {
	User() user.User
	Work() work.Work
}

// Normalized represents a normalized account
type Normalized interface {
}

// Service represents an account service
type Service interface {
	Save(ins Account, amountOfWorkToVerify int) error
}

// CreateAccountParams represents a CreateAccount params
type CreateAccountParams struct {
	User user.User
	Work work.Work
}

// CreateServiceParams represents a CreateService params
type CreateServiceParams struct {
	UserRepository   user.Repository
	WalletRepository wallet.Repository
	EntityService    entity.Service
}

// CreateSDKServiceParams represents the CreateSDKService params
type CreateSDKServiceParams struct {
	PK          crypto.PrivateKey
	Client      applications.Client
	RoutePrefix string
}

// SDKFunc represents the account SDK func
var SDKFunc = struct {
	Create           func(params CreateAccountParams) Account
	Normalize        func(ins Account) Normalized
	Denormalize      func(ins interface{}) Account
	CreateService    func(params CreateServiceParams) Service
	CreateSDKService func(params CreateSDKServiceParams) Service
}{
	Create: func(params CreateAccountParams) Account {
		out, outErr := createAccount(params.User, params.Work)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	Normalize: func(ins Account) Normalized {
		out, outErr := createNormalizedAccount(ins)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	Denormalize: func(ins interface{}) Account {

		if data, ok := ins.([]byte); ok {
			ptr := new(normalizedAccount)
			jsErr := cdc.UnmarshalJSON(data, ptr)
			if jsErr != nil {
				panic(jsErr)
			}

			out, outErr := fromNormalizedToAccount(ptr)
			if outErr != nil {
				panic(outErr)
			}

			return out
		}

		if normalized, ok := ins.(Normalized); ok {
			out, outErr := fromNormalizedToAccount(normalized)
			if outErr != nil {
				panic(outErr)
			}

			return out
		}

		panic(errors.New("the given instance cannot be denormalized to an Account instance"))
	},
	CreateService: func(params CreateServiceParams) Service {
		userRepresentation := user.SDKFunc.CreateRepresentation()
		out := createService(params.UserRepository, params.WalletRepository, params.EntityService, userRepresentation)
		return out
	},
	CreateSDKService: func(params CreateSDKServiceParams) Service {
		out := createSDKService(params.PK, params.Client, params.RoutePrefix)
		return out
	},
}
