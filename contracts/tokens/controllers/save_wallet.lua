local json = require("json")

-- save a new wallet
function saveWallet(from, path, params, data, sig)
    local newWalletData = json.decode(data)
    local newWallet = Wallet:create(newWalletData.pub_key, newWalletData.created_on)
    isSaved = newWallet:save()

    -- if saved suucessfully
    if isSaved then
        return {
            code = 0,
            log="success",
            gazUsed=0,
            tags={
                {
                    key=path .. "/" .. newWallet.pub_key,
                    value=json.encode(newWallet)
                }
            }
        }
    end

    return {
        code = 2,
        log="the wallet could not be saved",
    }
end
