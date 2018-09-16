require("datastore")

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
        setmetatable(retWal, Wallet)
        return retWal
    end

    -- retrieveTokenByID retrieves a token by its ID.  If the token doesn't exists, returns null
    function Repository:retrieveTokenByID(id)
        local db = tables.load()
        local keyname = generateKeyname("token", "id", id)
        retTok = db:retrieve(keyname)
        setmetatable(retTok, Token)
        return retTok
    end
-- class Repository
