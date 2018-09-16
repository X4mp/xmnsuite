local json = require("json")

-- retrieve wallet by its ID
function retrieveTokenByID(from, path, params, sig)
    local rep = Repository:create()
    local token = rep:retrieveTokenByID(params.id)
    if token == null then
        return {
            code = 1,
            log="the token (path: ".. path ..") was not found",
        }
    end

    return {
        code = 0,
        log="success",
        key=path,
        value=json.encode(token)
    }
end
