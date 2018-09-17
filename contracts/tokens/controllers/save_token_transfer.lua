local json = require("json")
local uuid = require("uuid")

-- transfer a token from a wallet to another
function saveTokenTransfer(from, path, params, data, sig)
    local rep = Repository:create()
    local token = rep:retrieveTokenByUUID(uuid.new(params.tokenID))
    if token == null then
        return {
            code = 1,
            log="the token (id: ".. params.tokenID ..") was not found",
        }
    end

    local from = rep:retrieveWalletByPubKey(params.fromPubKey)
    if from == null then
        return {
            code = 1,
            log="the from wallet (pubKey: ".. params.fromPubKey ..") was not found",
        }
    end

    local to = rep:retrieveWalletByPubKey(params.toPubKey)
    if to == null then
        return {
            code = 1,
            log="the to wallet (pubKey: ".. params.toPubKey ..") was not found",
        }
    end

    -- make sure the from token have enough tokens to execute the transfer:

    -- create the transfer:
    local transf = Transfer:create(params.id, token, from, to, params.created_on)
    isSaved = transf:save()

    if isSaved then
        return {
            code = 0,
            log="success",
            gazUsed=0,
            tags={
                {
                    key=path,
                    value=json.encode(transf)
                }
            }
        }
    end

    return {
        code = 2,
        log="the transfer could not be executed",
    }
end
