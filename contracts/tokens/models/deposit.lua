require("datastore")

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

    -- save saves a deposit instance to the database.  Returns true if successful, false otherwise
    function Deposit:save()
        local db = tables.load()
        local keynameDepositByID = generateKeyname("deposit", "uuid", self.uid:string())
        local retAmountSaved = db:save({key=keynameDepositByID, table={
            uid = self.uid,
            wallet_pub_key = self.wallet_pub_key,
            token_uuid = self.token_uuid,
            amount = self.amount,
            created_on = self.created_on
        }})

        if retAmountSaved ~= 1 then
            return false
        end

        -- save the desposit IDs in these keys as well, to make it possible to retrieve them by tokenUUID and walletPubKey:
        keys = {
            generateKeyname("deposit", "wallet_pub_key", self.wallet_pub_key),
            generateKeyname("deposit", "token_uuid", self.token_uuid:string())
        }

        local lst = lists.load()


        return true
    end

-- class Deposit
