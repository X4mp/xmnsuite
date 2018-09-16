require("datastore")

Transfer = {} --class
Transfer.__index = Transfer

    -- create creates a new Transfer instance
    function Transfer:create(id, token, from, to, createdOn)
        local trans = {
            id = id,
            token_id = token.id,
            from = from.pub_key,
            to = to.pub_key,
            created_on = createdOn
        }

        setmetatable(trans, Transfer)
        return trans
    end

-- class Transfer
