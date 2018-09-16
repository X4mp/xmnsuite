local chain = require("chain")

-- helpers
require("helpers/helper")

-- repositories
require("repositories/repository")

-- models
require("models/deposit")
require("models/token")
require("models/transfer")
require("models/wallet")

-- controllers
require("controllers/delete_token_by_id")
require("controllers/delete_wallet_by_pubkey")
require("controllers/retrieve_deposits_by_token_id")
require("controllers/retrieve_deposits_by_wallet_pubkey_and_token_id")
require("controllers/retrieve_deposits_by_wallet_pubkey")
require("controllers/retrieve_token_by_id")
require("controllers/retrieve_wallet_by_pubkey")
require("controllers/save_token_transfer")
require("controllers/save_token")
require("controllers/save_wallet")

-- chain
chain.chain().load({
    namespace = "xmn",
    name = "token",
    apps = {
        chain.app().new({
            version = "18.09.06",
            beginBlockIndex = 0,
            endBlockIndex = -1,
            router = chain.router().new({
                key = "router-roles",
                routes = {
                    chain.route().new("save", "/wallets", saveWallet),
                    chain.route().new("delete", "/wallets/<pubkey|[0-9a-f]{74}>", deleteWalletByPubKey),
                    chain.route().new("retrieve", "/wallets/<pubkey|[0-9a-f]{74}>", retrieveWalletByPubKey),
                    chain.route().new("retrieve", "/wallets/<pubkey|[0-9a-f]{74}>/deposits", retrieveDepositsByWalletPubKey),
                    chain.route().new("retrieve", "/wallets/<pubkey|[0-9a-f]{74}>/tokens/<tokenID|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>/deposits", retrieveDepositsByWalletPubKeyAndTokenID),
                    chain.route().new("save", "/tokens", saveToken),
                    chain.route().new("delete", "/tokens/<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>", deleteTokenByID),
                    chain.route().new("retrieve", "/tokens/<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>", retrieveTokenByID),
                    chain.route().new("retrieve", "/tokens/<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>/deposits", retrieveDepositsByTokenID),
                    chain.route().new("save", "/tokens/<tokenID|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>/from-<fromPubKey|[0-9a-f]{74}>/to-<toPubKey|[0-9a-f]{74}>", saveTokenTransfer),
                }
            })
        })
    }
})
