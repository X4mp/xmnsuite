-- Unit testing starts
require('luaunit')

TestRouter = {} --class
    function TestRouter:testAddRoute_executeRoute_Success()
        -- variables
        queryHandlerFn = function(uriParams, queryParams)

            -- this is the handler func:

            return {
                header: {
                    statusCode: xrouter:statusOk,
                    lines: {
                        name: "X-Requested-With",
                        value: "XMLHttpRequest",
                    }
                },
                body: "this is the body"
            }
        end

        txHandlerFunc = function(uriParams, queryParams, txData)

            -- this is the handler func:

            return {
                header: {
                    statusCode: xrouter:statusOk,
                    lines: {
                        name: "X-Requested-With",
                        value: "XMLHttpRequest",
                    }
                },
                body: "this is the body"
            }
        end

        firstRoute = {
            pattern="/articles/{category}/{id:[0-9]+}",
            handler=queryHandlerFn
        }

        secondRoute = {
            pattern="/another-article/{category}/{id:[0-9]+}",
            handler=txHandlerFunc
        }

        xrouter:queryRoutes(firstRoute)
        xrouter:transactionRoutes(secondRoute)

    end

-- class TestRouter

LuaUnit:run()
