-- delete a token by its ID
function deleteTokenByID(from, path, params, sig)
    local rep = Repository:create()
    local token = rep:retrieveTokenByID(params.id)
    if token == null then
        return {
            code = 1,
            log="the token (path: ".. path ..") was not found",
        }
    end

    isDeleted = token:delete()

    -- if deleted successfully
    if isDeleted then
        return {
            code = 0,
            log="success",
        }
    end

    return {
        code = 2,
        log="the token could not be deleted",
    }
end
