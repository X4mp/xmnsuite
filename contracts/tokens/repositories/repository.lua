require("datastore")
local uuid = require("uuid")

Repository = {} --class
Repository.__index = Repository

    -- create creates a new repository
    function Repository:create()
        local rep = {}
        setmetatable(rep, Repository)
        return rep
    end

    -- retrieveByKeyname retrieves data by its keyname.  If the data doesn't exists, returns null
    function Repository:retrieveByKeyname(keyname)
        local db = tables.load()
        retValue = db:retrieve(keyname)
        if retValue == null then
            return null
        end

        return retValue
    end

    -- retrieveWalletByPubKey retrieves a wallet by its public key.  If the wallet doesn't exists, returns null
    function Repository:retrieveWalletByPubKey(pubKey)
        -- create the keyname:
        local keyname = generateKeyname("wallet", "pubkey", pubKey)

        -- retrieve then return:
        value = self:retrieveByKeyname(keyname)
        return Wallet:load(value)
    end

    -- retrieveTokenByUUID retrieves a token by its UUID.  If the token doesn't exists, returns null
    function Repository:retrieveTokenByUUID(uid)
        -- create the keyname:
        local keyname = generateKeyname("token", "uuid", uid:string())

        -- retrieve then return:
        value = self:retrieveByKeyname(keyname)
        return Token:load(value)
    end

    -- retrieveDepositByUUID retrieves a deposit by its UUID.  If the deposit doesn't exists, returns null
    function Repository:retrieveDepositByUUID(uid)
        -- create keyname:
        local keyname = generateKeyname("deposit", "uuid", uid:string())

        -- retrieve then return:
        value = self:retrieveByKeyname(keyname)
        return Deposit:load(value)
    end

    -- retrieveDepositsByUUIDs retrieve deposits by their UUID
    function Repository:retrieveDepositsByUUIDs(uuids, index, amount)
        -- for each uuid, retrieve the deposit:
        local deposits = {}
        for index, uuidAsString in pairs(uuids) do
            oneDeposit = self:retrieveDepositByUUID(uuid.new(uuidAsString))
            if oneDeposit ~= null then
                table.insert(deposits, oneDeposit)
            end
        end

        return deposits
    end

    -- retrieveDepositsByKeyname retrieve a set of deposits by its keyname
    function Repository:retrieveDepositsByKeyname(keyname, index, amount)
        -- load the databases:
        local sts = sets.load()

        -- index:
        ind = 0
        if index ~= null then
            ind = index
        end

        -- amount:
        am = -1
        if amount ~= null then
            am = amount
        end

        -- retrieve the list of uuids:
        local uuids = sts:retrieve(keyname, ind, am)
        if uuids == null then
            return {}
        end

        -- return:
        return self:retrieveDepositsByUUIDs(uuids)
    end

    -- retrieveDepositsByTokenUUID retrieve deposits by their the tokenUUID
    function Repository:retrieveDepositsByTokenUUID(tokenUUID, index, amount)
        -- create the keyname:
        local keyname = generateKeyname("deposit", "token_uuid", tokenUUID:string())

        -- return:
        return self:retrieveDepositsByKeyname(keyname, index, amount)
    end

    -- retrieveDepositsByWalletPubKey retrieve deposits by the wallet pub key:
    function Repository:retrieveDepositsByWalletPubKey(pubKey, index, amount)
        -- create the keyname:
        local keyname = generateKeyname("deposit", "wallet_pub_key", pubKey)

        -- return:
        return self:retrieveDepositsByKeyname(keyname, index, amount)
    end

    -- retrieveDepositsByWalletPubKeyAndTokenUUID retrieve deposits by the wallet pub key and tokenID:
    function Repository:retrieveDepositsByWalletPubKeyAndTokenUUID(pubKey, tokenUUID, index, amount)
        -- load the databases:
        local sts = sets.load()

        -- create the walletPubKey keyname:
        local walletKeyname = generateKeyname("deposit", "wallet_pub_key", pubKey)

        -- create the tokenID keyname:
        local tokenKeyname = generateKeyname("deposit", "token_uuid", tokenUUID:string())

        -- interstore:
        local outputKeyname = "inter:" .. walletKeyname .. "|" .. tokenKeyname
        sts:interstore(outputKeyname, walletKeyname, tokenKeyname)

        -- return:
        return self:retrieveDepositsByKeyname(outputKeyname, index, amount)
    end

-- class Repository
