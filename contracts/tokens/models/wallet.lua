require("datastore")

Wallet = {} --class
Wallet.__index = Wallet

    -- create creates a new wallet instance
    function Wallet:create(pubKey, createdOn)
        local wallet = {
            pub_key = pubKey,
            created_on = createdOn
        }

        setmetatable(wallet, Wallet)
        return wallet
    end

    -- load loads the data into an object
    function Wallet:load(data)
        if data == null then
            return null
        end

        local wallet = {
            pub_key = data.pub_key,
            created_on = data.created_on
        }

        setmetatable(wallet, Wallet)
        return wallet
    end

    -- toData converts the object to data
    function Wallet:toData()
        return {
            pub_key = self.pub_key,
            created_on = self.created_on
        }
    end

    -- save saves a wallet instance to the database.  Returns true if successful, false otherwise
    function Wallet:save()
        local db = tables.load()
        local keyname = generateKeyname("wallet", "pubkey", self.pub_key)
        local retAmountSaved = db:save({key=keyname, table=self:toData()})

        if retAmountSaved ~= 1 then
            return false
        end

        return true
    end

    -- delete deletes a wallet instance from the database.  Returns true if successful, false otherwise
    function Wallet:delete()
        local db = tables.load()
        local keyname = generateKeyname("wallet", "pubkey", self.pub_key)
        local retAmountDeleted = db:delete(keyname)
        if retAmountDeleted ~= 1 then
            return false
        end

        return true
    end
-- class Wallet
