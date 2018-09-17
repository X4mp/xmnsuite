local uuid = require("uuid")
local json = require("json")

-- retrieve depots by the tokenUUID
function retrieveDepositsByTokenUUID(from, path, params, sig)
    -- create repository:
    rep = Repository:create()
    deps = rep:retrieveDepositsByTokenUUID(uuid.new(params.tokenid), params.index, params.amount)

    return {
        code = 0,
        log="success",
        key=path,
        value=json.encode(toData(deps))
    }
end
