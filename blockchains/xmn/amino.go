package xmn

import (
	amino "github.com/tendermint/go-amino"
	applications "github.com/xmnservices/xmnsuite/routers"
)

const (

	// XMNSuiteApplicationsXMNWallet represents the xmnsuite xmn Wallet resource
	XMNSuiteApplicationsXMNWallet = "xmnsuite/xmn/Wallet"

	// XMNSuiteApplicationsXMNWalletPartialSet represents the xmnsuite xmn WalletPartialSet resource
	XMNSuiteApplicationsXMNWalletPartialSet = "xmnsuite/xmn/WalletPartialSet"

	// XMNSuiteApplicationsXMNUser represents the xmnsuite xmn User resource
	XMNSuiteApplicationsXMNUser = "xmnsuite/xmn/User"

	// XMNSuiteApplicationsXMNUserRequest represents the xmnsuite xmn UserRequest resource
	XMNSuiteApplicationsXMNUserRequest = "xmnsuite/xmn/UserRequest"

	// XMNSuiteApplicationsXMNUserRequestVote represents the xmnsuite xmn UserRequestVote resource
	XMNSuiteApplicationsXMNUserRequestVote = "xmnsuite/xmn/UserRequestVote"

	// XMNSuiteApplicationsXMNInitialDeposit represents the xmnsuite xmn InitialDeposit resource
	XMNSuiteApplicationsXMNInitialDeposit = "xmnsuite/xmn/InitialDeposit"

	// XMNSuiteApplicationsXMNToken represents the xmnsuite xmn Token resource
	XMNSuiteApplicationsXMNToken = "xmnsuite/xmn/Token"

	// XMNSuiteApplicationsXMNGenesis represents the xmnsuite xmn Genesis resource
	XMNSuiteApplicationsXMNGenesis = "xmnsuite/xmn/Genesis"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Dependencies
	applications.Register(codec)

	// Wallet
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Wallet)(nil), nil)
		codec.RegisterConcrete(&wallet{}, XMNSuiteApplicationsXMNWallet, nil)
	}()

	// WalletPartialSet
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*WalletPartialSet)(nil), nil)
		codec.RegisterConcrete(&walletPartialSet{}, XMNSuiteApplicationsXMNWalletPartialSet, nil)
	}()

	// User
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*User)(nil), nil)
		codec.RegisterConcrete(&user{}, XMNSuiteApplicationsXMNUser, nil)
	}()

	// UserRequest
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*UserRequest)(nil), nil)
		codec.RegisterConcrete(&userRequest{}, XMNSuiteApplicationsXMNUserRequest, nil)
	}()

	// UserRequestVote
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*UserRequestVote)(nil), nil)
		codec.RegisterConcrete(&userRequestVote{}, XMNSuiteApplicationsXMNUserRequestVote, nil)
	}()

	// InitialDeposit
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*InitialDeposit)(nil), nil)
		codec.RegisterConcrete(&initialDeposit{}, XMNSuiteApplicationsXMNInitialDeposit, nil)
	}()

	// Token
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Token)(nil), nil)
		codec.RegisterConcrete(&token{}, XMNSuiteApplicationsXMNToken, nil)
	}()

	// Genesis
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Genesis)(nil), nil)
		codec.RegisterConcrete(&genesis{}, XMNSuiteApplicationsXMNGenesis, nil)
	}()
}
