require("datastore")
local uuid = require("uuid")

Deposit = {} --class
Deposit.__index = Deposit

    -- create create a new Deposit instance
    function Deposit:create(uid, walletPubKey, tokenUUID, amount, createdOn)
        local deposit = {
            uid = uid,
            wallet_pub_key = walletPubKey,
            token_uuid = tokenUUID,
            amount = amount,
            created_on = createdOn
        }

        setmetatable(deposit, Deposit)
        return deposit
    end

    -- load loads the data into an object
    function Deposit:load(data)
        if data == null then
            return null
        end

        local dep = {
            uid = uuid.new(data.uid),
            wallet_pub_key = data.wallet_pub_key,
            token_uuid = uuid.new(data.token_uuid),
            amount = data.amount,
            created_on = data.created_on
        }

        setmetatable(dep, Deposit)
        return dep
    end

    -- toData converts the object to data
    function Deposit:toData()
        return {
            uid = self.uid:string(),
            wallet_pub_key = self.wallet_pub_key,
            token_uuid = self.token_uuid:string(),
            amount = self.amount,
            created_on = self.created_on
        }
    end

    -- save saves a deposit instance to the database.  Returns true if successful, false otherwise
    function Deposit:save()
        -- load the databases:
        local sts = sets.load()
        local db = tables.load()

        -- add the deposit by wallet_pub_key:
        firstKeyname = generateKeyname("deposit", "wallet_pub_key", self.wallet_pub_key)
        firstAmountSaved = sts:add(firstKeyname, self.uid:string())
        if firstAmountSaved ~= 1 then
            return false
        end

        -- add the deposit by token uuid:
        secondKeyname = generateKeyname("deposit", "token_uuid", self.token_uuid:string())
        secondAmountSaved = sts:add(secondKeyname, self.uid:string())
        if secondAmountSaved ~= 1 then
            return false
        end

        -- save the deposit:
        local keynameDepositByID = generateKeyname("deposit", "uuid", self.uid:string())
        if db:save({key=keynameDepositByID, table=self:toData()}) ~= 1 then
            return false
        end

        return true
    end

-- class Deposit
