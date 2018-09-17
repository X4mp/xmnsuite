local uuid = require("uuid")
local json = require("json")

-- retrieve deposits by wallet pubKey
function retrieveDepositsByWalletPubKey(from, path, params, sig)
    -- create repository:
    rep = Repository:create()
    deps = rep:retrieveDepositsByWalletPubKey(params.pubkey, params.index, params.amount)

    return {
        code = 0,
        log="success",
        key=path,
        value=json.encode(toData(deps))
    }
end
