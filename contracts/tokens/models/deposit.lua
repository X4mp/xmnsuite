require("datastore")

Deposit = {} --class
Deposit.__index = Deposit

    -- create create a new Deposit instance
    function Deposit:create(id, walletID, tokenID, amount, createdOn)
        local deposit = {
            id = id,
            wallet_id = walletID,
            token_id = tokenID,
            amount = amount,
            created_on = createdOn
        }

        setmetatable(deposit, Deposit)
        return deposit
    end

    -- save saves a deposit instance to the database.  Returns true if successful, false otherwise
    function Deposit:save()
        local db = tables.load()
        local keynameDepositByID = generateKeyname("deposit", "id", self.id)
        local retAmountSaved = db:save({key=keynameDepositByID, table={
            id = self.id,
            wallet_id = self.wallet_id,
            token_id = self.token_id,
            amount = self.amount,
            created_on = self.created_on
        }})

        if retAmountSaved ~= 1 then
            return false
        end

        -- save the desposit IDs in these keys as well, to make it possible to retrieve them by tokenID and walletPubKey:
        keys = {
            generateKeyname("deposit", "wallet_id", self.wallet_id),
            generateKeyname("deposit", "token_id", self.token_id)
        }

        local lst = lists.load()


        return true
    end

-- class Deposit
