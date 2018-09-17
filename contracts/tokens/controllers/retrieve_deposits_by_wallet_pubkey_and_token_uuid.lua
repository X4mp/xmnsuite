local uuid = require("uuid")
local json = require("json")

-- retrieve deposits by wallet pubKey and tokenUUID
function retrieveDepositsByWalletPubKeyAndTokenUUID(from, path, params, sig)
    -- create repository:
    rep = Repository:create()
    deps = rep:retrieveDepositsByWalletPubKeyAndTokenUUID(params.pubkey, uuid.new(params.tokenid), params.index, params.amount)

    return {
        code = 0,
        log="success",
        key=path,
        value=json.encode(toData(deps))
    }
end
