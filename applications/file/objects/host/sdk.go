package host

import (
	"net"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet"
)

// Host represents a host
type Host interface {
	ID() *uuid.UUID
	Wallet() wallet.Wallet
	BandwithPricePerChunk() int
	StoragePricePerBlock() int
	IP() net.IP
	Port() int
}
