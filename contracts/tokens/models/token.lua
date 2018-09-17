require("datastore")
local uuid = require("uuid")

Token = {} --class
Token.__index = Token

    -- create create a new Token instance
    function Token:create(uid, symbol, name, description, createdOn)
        local tok = {
            uid = uid,
            symbol = symbol,
            name = name,
            description = description,
            created_on = createdOn
        }

        setmetatable(tok, Token)
        return tok
    end

    -- load loads the data into an object
    function Token:load(data)
        local tok = {
            uid = uuid.new(data.uid),
            symbol = data.symbol,
            name = data.name,
            description = data.description,
            created_on = data.created_on
        }

        setmetatable(tok, Token)
        return tok
    end

    -- toData converts the object to data
    function Token:toData()
        return {
            uid = self.uid:string(),
            symbol = self.symbol,
            name = self.name,
            description = self.description,
            created_on = self.created_on
        }
    end

    -- save saves a token instance to the database.  Returns true if successful, false otherwise
    function Token:save()
        local db = tables.load()
        local keyname = generateKeyname("token", "uuid", self.uid:string())
        local retAmountSaved = db:save({key=keyname, table=self:toData()})

        if retAmountSaved ~= 1 then
            return false
        end

        return true
    end

    -- delete deletes a token instance from the database.  Returns true if successful, false otherwise
    function Token:delete()
        local db = tables.load()
        local keyname = generateKeyname("token", "uuid", self.uid:string())
        local retAmountDeleted = db:delete(keyname)
        if retAmountDeleted ~= 1 then
            return false
        end

        return true
    end
-- class Token
