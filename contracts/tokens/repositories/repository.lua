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

    -- retrieveWalletByPubKey retrieves a wallet by its public key.  If the wallet doesn't exists, returns null
    function Repository:retrieveWalletByPubKey(pubKey)
        local db = tables.load()
        local keyname = generateKeyname("wallet", "pubkey", pubKey)
        retWal = db:retrieve(keyname)
        if retWal == null then
            return null
        end

        return Wallet:load(retWal)
    end

    -- retrieveTokenByUUID retrieves a token by its UUID.  If the token doesn't exists, returns null
    function Repository:retrieveTokenByUUID(uid)
        local db = tables.load()
        local keyname = generateKeyname("token", "uuid", uid:string())
        retTok = db:retrieve(keyname)
        if retTok == null then
            return null
        end

        return Token:load(retTok)
    end
-- class Repository
