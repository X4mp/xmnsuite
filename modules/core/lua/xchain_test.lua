-- Unit testing starts
require('luaunit')

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
        first = xapp.new({
            version = "16.02.01",
            beginBlockIndex = 0,
            endBlockIndex = 200,
            router = xrouter.new({
                key = "this-is-the-router-key:16.02.01",
                routes = {
                    xroute.new("retrieve", "/videos/<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>", retrieveVideosByID),
                    xroute.new("delete", "/videos/<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>", deleteVideoByID),
                    xroute.new("save", "/videos", saveVideo),
                }
            })
        })

        second = xapp.new({
            version = "17.03.09",
            beginBlockIndex = 200,
            endBlockIndex = -1,
            router = xrouter.new({
                key = "this-is-the-router-key:16.02.01",
                routes = {
                    xroute.new("retrieve", "/videos/<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>", retrieveVideosByID),
                    xroute.new("delete", "/videos/<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>", deleteVideoByID),
                    xroute.new("save", "/videos", saveVideo),
                }
            })
        })

        xchain.load({
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
