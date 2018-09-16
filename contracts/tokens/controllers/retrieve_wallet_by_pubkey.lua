local json = require("json")

-- retrieve wallet by its pubKey
function retrieveWalletByPubKey(from, path, params, sig)
    local rep = Repository:create()
    local wallet = rep:retrieveWalletByPubKey(params.pubkey)
    if wallet == null then
        return {
            code = 1,
            log="the wallet (path: ".. path ..") could not be found",
        }
    end

    return {
        code = 0,
        log="success",
        key=path,
        value=json.encode(wallet)
    }
end
