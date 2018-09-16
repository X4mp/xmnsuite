local json = require("json")

-- save a new token
function saveToken(from, path, params, data, sig)
    -- retrieve the token data:
    local newTokenData = json.decode(data)

    -- retrieve the wallet:
    local rep = Repository:create()
    local wallet = rep:retrieveWalletByPubKey(newTokenData.wallet)
    if wallet == null then
        return {
            code = 1,
            log="the wallet (pubKey: ".. newTokenData.wallet ..") could not be found",
        }
    end

    -- create the new token instance, then save it:
    local newToken = Token:create(newTokenData.id, newTokenData.symbol, newTokenData.name, newTokenData.description, newTokenData.created_on)
    isTokenSaved = newToken:save()

    -- if token not saved successfully:
    if isTokenSaved == false then
        return {
            code = 2,
            log="the token could not be saved",
        }
    end

    -- save the deposit to the wallet:
    local deposit = Deposit:create("b081e434-cfb2-4f1f-9f0b-015a2077af31", wallet.id, newToken.id, newTokenData.amount, os.time())
    isDepositSaved = deposit:save()

    -- if deposit not saved successfully:
    if isDepositSaved == false then
        return {
            code = 2,
            log="the initial token deposit could not be saved",
        }
    end

    return {
        code = 0,
        log="success",
        gazUsed=0,
        tags={
            {
                key=path .. "/" .. newToken.id,
                value=json.encode(newToken)
            }
        }
    }
end
