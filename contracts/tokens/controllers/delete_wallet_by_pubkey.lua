-- delete a wallet by its pubKey
function deleteWalletByPubKey(from, path, params, sig)
    local rep = Repository:create()
    local wallet = rep:retrieveWalletByPubKey(params.pubkey)
    if wallet == null then
        return {
            code = 1,
            log="the wallet (path: ".. path ..") was not found",
        }
    end

    isDeleted = wallet:delete()

    -- if deleted successfully
    if isDeleted then
        return {
            code = 0,
            log="success",
        }
    end

    return {
        code = 2,
        log="the wallet could not be deleted",
    }
end
