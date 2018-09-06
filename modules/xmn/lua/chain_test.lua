-- Unit testing starts
require('luaunit')
local xmn = require("xmn")

TestChain = {} --class
    function TestChain:testCreate_Success()
        -- func handlers:
        function retrieveVideosByID(from, path, params, sig)
            return {}
        end

        function saveVideo(from, path, params, data, sig)
            return {}
        end

        function deleteVideoByID(from, path, params, sig)
            return {}
        end

        -- create the application
        first = xmn.app().new({
            version = "16.02.01",
            beginBlockIndex = 0,
            endBlockIndex = 200,
            router = xmn.router().new({
                key = "this-is-the-router-key:16.02.01",
                routes = {
                    xmn.route().new("retrieve", "/videos/<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>", retrieveVideosByID),
                    xmn.route().new("delete", "/videos/<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>", deleteVideoByID),
                    xmn.route().new("save", "/videos", saveVideo),
                }
            })
        })

        second = xmn.app().new({
            version = "17.03.09",
            beginBlockIndex = 200,
            endBlockIndex = -1,
            router = xmn.router().new({
                key = "this-is-the-router-key:16.02.01",
                routes = {
                    xmn.route().new("retrieve", "/videos/<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>", retrieveVideosByID),
                    xmn.route().new("delete", "/videos/<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>", deleteVideoByID),
                    xmn.route().new("save", "/videos", saveVideo),
                }
            })
        })

        xmn.chain().load({
            namespace = "xmn",
            name = "messages",
            apps = {
                first,
                second
            }
        })

        -- verify:
        -- assertEquals(type(firstRte), "xroute")
    end

-- class TestChain

LuaUnit:run()
