require("datastore")

Token = {} --class
Token.__index = Token

    -- create create a new Token instance
    function Token:create(id, symbol, name, description, createdOn)
        local tok = {
            id = id,
            symbol = symbol,
            name = name,
            description = description,
            created_on = createdOn
        }

        setmetatable(tok, Token)
        return tok
    end

    -- save saves a token instance to the database.  Returns true if successful, false otherwise
    function Token:save()
        local db = tables.load()
        local keyname = generateKeyname("token", "id", self.id)
        local retAmountSaved = db:save({key=keyname, table={
            id = self.id,
            symbol = self.symbol,
            name = self.name,
            description = self.description,
            created_on = self.created_on
        }})

        if retAmountSaved ~= 1 then
            return false
        end

        return true
    end

    -- delete deletes a token instance from the database.  Returns true if successful, false otherwise
    function Token:delete()
        local db = tables.load()
        local keyname = generateKeyname("token", "id", self.id)
        local retAmountDeleted = db:delete(keyname)
        if retAmountDeleted ~= 1 then
            return false
        end

        return true
    end
-- class Token
