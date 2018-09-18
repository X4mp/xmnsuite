local chain = require("chain")

-- helpers
require("helpers/helper")

-- models
require("models/deposit")
require("models/token")
require("models/wallet")

-- repositories
require("repositories/repository")

-- controllers
require("controllers/delete_token_by_uuid")
require("controllers/delete_wallet_by_pubkey")
require("controllers/retrieve_deposits_by_token_uuid")
require("controllers/retrieve_deposits_by_wallet_pubkey_and_token_uuid")
require("controllers/retrieve_deposits_by_wallet_pubkey")
require("controllers/retrieve_token_by_uuid")
require("controllers/retrieve_wallet_by_pubkey")
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
                    chain.route().new("delete", "/wallets/<pubkey|[0-9a-f]{64}>", deleteWalletByPubKey),
                    chain.route().new("retrieve", "/wallets/<pubkey|[0-9a-f]{64}>", retrieveWalletByPubKey),
                    chain.route().new("save", "/tokens", saveToken),
                    chain.route().new("delete", "/tokens/<uid|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>", deleteTokenByUUID),
                    chain.route().new("retrieve", "/tokens/<uid|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>", retrieveTokenByUUID),
                    chain.route().new("retrieve", "/deposits/<pubkey|[0-9a-f]{64}>", retrieveDepositsByWalletPubKey),
                    chain.route().new("retrieve", "/deposits/<tokenid|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>", retrieveDepositsByTokenUUID),
                    chain.route().new("retrieve", "/deposits/<pubkey|[0-9a-f]{64}>/<tokenid|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>", retrieveDepositsByWalletPubKeyAndTokenUUID),
                }
            })
        })
    }
})
