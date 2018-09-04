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
        first = xapp.new("16.02.01", 0, 200, xrouter.new(
            xroute.new("retrieve", "/videos/<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>", retrieveVideosByID),
            xroute.new("delete", "/videos/<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>", deleteVideoByID),
            xroute.new("save", "/videos", saveVideo)
        ))

        second = xapp.new("17.03.09", 200, -1, xrouter.new(
            xroute.new("retrieve", "/videos/<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>", retrieveVideosByID),
            xroute.new("delete", "/videos/<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>", deleteVideoByID),
            xroute.new("save", "/videos", saveVideo)
        ))

        xchain.load(first, second)

        -- verify:
        -- assertEquals(type(firstRte), "xroute")
    end

-- class TestChain

LuaUnit:run()
