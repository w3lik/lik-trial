--- 客户端尺寸
japi._clientHeight = japi.DZ_GetClientHeight()
japi._clientWidth = japi.DZ_GetClientWidth()
--- frame
japi._frAdaptive = Array()
japi._frEsc = Array()
--- 鼠标事件
J.Japi["DzTriggerRegisterMouseEventByCode"](nil, GAME_KEY_MOUSE_LEFT, GAME_KEY_ACTION_PRESS, false, function()
    if (true ~= japi.DZ_IsWindowActive()) then
        return
    end
    if (event.asyncHas("mouse", EVENT.Mouse.LeftClick)) then
        local triggerPlayer = Player(1 + J.GetPlayerId(japi.DZ_GetTriggerKeyPlayer()))
        async.call(triggerPlayer, function()
            local events = event.get(event._async, "mouse", EVENT.Mouse.LeftClick)
            events:backEach(function(_, v)
                local r = J.Promise(v, nil, nil, { triggerPlayer = triggerPlayer, rx = japi.MouseRX(), ry = japi.MouseRY() })
                if (r ~= nil) then
                    if (type(r) == "boolean" and r == false) then
                        return false
                    end
                end
            end)
        end)
    end
end)
J.Japi["DzTriggerRegisterMouseEventByCode"](nil, GAME_KEY_MOUSE_LEFT, GAME_KEY_ACTION_RELEASE, false, function()
    if (true ~= japi.DZ_IsWindowActive()) then
        return
    end
    if (event.asyncHas("mouse", EVENT.Mouse.LeftRelease)) then
        local triggerPlayer = Player(1 + J.GetPlayerId(japi.DZ_GetTriggerKeyPlayer()))
        async.call(triggerPlayer, function()
            local events = event.get(event._async, "mouse", EVENT.Mouse.LeftRelease)
            events:backEach(function(_, v)
                local r = J.Promise(v, nil, nil, { triggerPlayer = triggerPlayer })
                if (r ~= nil) then
                    if (type(r) == "boolean" and r == false) then
                        return false
                    end
                end
            end)
        end)
    end
end)
J.Japi["DzTriggerRegisterMouseEventByCode"](nil, GAME_KEY_MOUSE_RIGHT, GAME_KEY_ACTION_PRESS, false, function()
    if (true ~= japi.DZ_IsWindowActive()) then
        return
    end
    if (event.asyncHas("mouse", EVENT.Mouse.RightClick)) then
        local triggerPlayer = Player(1 + J.GetPlayerId(japi.DZ_GetTriggerKeyPlayer()))
        async.call(triggerPlayer, function()
            local events = event.get(event._async, "mouse", EVENT.Mouse.RightClick)
            events:backEach(function(_, v)
                local r = J.Promise(v, nil, nil, { triggerPlayer = triggerPlayer, rx = japi.MouseRX(), ry = japi.MouseRY() })
                if (r ~= nil) then
                    if (type(r) == "boolean" and r == false) then
                        return false
                    end
                end
            end)
        end)
    end
end)
J.Japi["DzTriggerRegisterMouseEventByCode"](nil, GAME_KEY_MOUSE_RIGHT, GAME_KEY_ACTION_RELEASE, false, function()
    if (true ~= japi.DZ_IsWindowActive()) then
        return
    end
    if (event.asyncHas("mouse", EVENT.Mouse.RightRelease)) then
        local triggerPlayer = Player(1 + J.GetPlayerId(japi.DZ_GetTriggerKeyPlayer()))
        async.call(triggerPlayer, function()
            local events = event.get(event._async, "mouse", EVENT.Mouse.RightRelease)
            events:backEach(function(_, v)
                local r = J.Promise(v, nil, nil, { triggerPlayer = triggerPlayer })
                if (r ~= nil) then
                    if (type(r) == "boolean" and r == false) then
                        return false
                    end
                end
            end)
        end)
    end
end)
J.Japi["DzTriggerRegisterMouseWheelEventByCode"](nil, false, function()
    if (true ~= japi.DZ_IsWindowActive()) then
        return
    end
    if (event.asyncHas("mouse", EVENT.Mouse.Wheel)) then
        local triggerPlayer = Player(1 + J.GetPlayerId(japi.DZ_GetTriggerKeyPlayer()))
        async.call(triggerPlayer, function()
            local events = event.get(event._async, "mouse", EVENT.Mouse.Wheel)
            events:backEach(function(_, v)
                local r = J.Promise(v, nil, nil, { triggerPlayer = triggerPlayer, delta = japi.DZ_GetWheelDelta() })
                if (r ~= nil) then
                    if (type(r) == "boolean" and r == false) then
                        return false
                    end
                end
            end)
        end)
    end
end)
J.Japi["DzTriggerRegisterWindowResizeEventByCode"](nil, false, function()
    if (event.asyncHas("window", EVENT.Window.Resize)) then
        local triggerPlayer = Player(1 + J.GetPlayerId(japi.DZ_GetTriggerKeyPlayer()))
        async.call(triggerPlayer, function()
            local events = event.get(event._async, "window", EVENT.Window.Resize)
            events:backEach(function(_, v)
                J.Promise(v, nil, nil, { triggerPlayer = triggerPlayer })
            end)
        end)
    end
end)
--- 键盘事件
for _, kb in pairx(KEYBOARD) do
    J.Japi["DzTriggerRegisterKeyEventByCode"](nil, kb, GAME_KEY_ACTION_PRESS, false, function()
        if (true ~= japi.DZ_IsWindowActive()) then
            return
        end
        if (true == japi.DZ_IsChatBoxOpen()) then
            return
        end
        local triggerKey = japi.DZ_GetTriggerKey()
        if (event.asyncHas("keyboard" .. triggerKey, EVENT.Keyboard.Press)) then
            local triggerPlayer = Player(1 + J.GetPlayerId(japi.DZ_GetTriggerKeyPlayer()))
            async.call(triggerPlayer, function()
                local events = event.get(event._async, "keyboard" .. triggerKey, EVENT.Keyboard.Press)
                events:backEach(function(_, v)
                    J.Promise(v, nil, nil, { triggerPlayer = triggerPlayer, triggerKey = triggerKey })
                end)
            end)
        end
    end)
    J.Japi["DzTriggerRegisterKeyEventByCode"](nil, kb, GAME_KEY_ACTION_RELEASE, false, function()
        if (true ~= japi.DZ_IsWindowActive()) then
            return
        end
        if (true == japi.DZ_IsChatBoxOpen()) then
            return
        end
        local triggerKey = japi.DZ_GetTriggerKey()
        if (event.asyncHas("keyboard" .. triggerKey, EVENT.Keyboard.Release)) then
            local triggerPlayer = Player(1 + J.GetPlayerId(japi.DZ_GetTriggerKeyPlayer()))
            async.call(triggerPlayer, function()
                local events = event.get(event._async, "keyboard" .. triggerKey, EVENT.Keyboard.Release)
                events:backEach(function(_, v)
                    J.Promise(v, nil, nil, { triggerPlayer = triggerPlayer, triggerKey = triggerKey })
                end)
            end)
        end
    end)
end

Game():onEvent(EVENT.Game.Start, "_lk_japi", function()
    -- adaptive
    local smooth = {}
    view.onResize("_li_resize", function()
        -- adaptive
        local pIdx = PlayerLocal():index()
        if (isClass(smooth[pIdx], TimerAsyncClass)) then
            smooth[pIdx]:setRemain(smooth[pIdx]:period())
            return
        end
        smooth[pIdx] = async.setTimeout(9, function()
            smooth[pIdx] = nil
            -- client
            japi._clientHeight = japi.DZ_GetClientHeight()
            japi._clientWidth = japi.DZ_GetClientWidth()
            ---@param fr Frame
            japi._frAdaptive:backEach(function(_, fr)
                if (instanceof(fr, FrameClass)) then
                    local s = fr:prop("unAdaptiveSize")
                    local r = fr:prop("unAdaptiveRelation")
                    if (type(s) == "table") then
                        fr:size(table.unpack(s))
                    end
                    if (type(r) == "table") then
                        fr:relation(table.unpack(r))
                    end
                end
            end)
        end)
    end)
    -- Esc
    keyboard.onRelease(KEYBOARD["Esc"], "_lk_esc", function()
        if (japi.DZ_IsChatBoxOpen()) then
            return
        end
        local keys = japi._frEsc:keys()
        local l = #keys
        if (l > 0) then
            local key = keys[l]
            local fr = japi._frEsc:get(key)
            japi._frEsc:set(key, nil)
            fr:show(false)
        end
    end)
    -- common refresh
    local rf = function()
        if (false == japi.DZ_IsWindowActive()) then
            return
        end
        -- cursor
        local tmpCursor = {
            triggerPlayer = PlayerLocal(),
            rx = japi.RX(japi.DZ_GetMouseXRelative()),
            ry = japi.RY(japi.DZ_GetClientHeight() - japi.DZ_GetMouseYRelative())
        }
        if (event.asyncHas("mouse", EVENT.Mouse.Move)) then
            local d
            if (type(japi._cursorLast) == "table") then
                local csl = japi._cursorLast
                d = vector2.distance(csl.rx, csl.ry, tmpCursor.rx, tmpCursor.ry)
            end
            if (d == nil or d > 0.003) then
                japi._cursorLast = tmpCursor
                async.call(PlayerLocal(), function()
                    local events = event.get(event._async, "mouse", EVENT.Mouse.Move)
                    events:backEach(function(_, v)
                        J.Promise(v, nil, nil, tmpCursor)
                    end)
                end)
            end
        end
        if (japi._cursor == nil or tmpCursor.rx ~= japi._cursor.rx or tmpCursor.ry ~= japi._cursor.ry) then
            japi._cursor = tmpCursor
            if (event.asyncHas("mouse", EVENT.Mouse.Instant)) then
                async.call(PlayerLocal(), function()
                    local events = event.get(event._async, "mouse", EVENT.Mouse.Instant)
                    events:backEach(function(_, v)
                        J.Promise(v, nil, nil, tmpCursor)
                    end)
                end)
            end
        end
    end
    --- 初始化 异步帧钟
    J.Japi["DzFrameSetUpdateCallbackByCode"](function()
        rf()
        local pid = PlayerLocal():index()
        for _, re in pairs(japi._asyncRefresh) do
            async._id = pid
            J.Promise(re)
            async._id = 0
        end
        japi._asyncExecDelayInc = japi._asyncExecDelayInc + 1
        local ts = japi._asyncExecDelay[japi._asyncExecDelayInc]
        if (ts ~= nil) then
            for fid, fcl in pairs(ts) do
                if (pid == fcl.i) then
                    async._id = fcl.i
                    J.Promise(fcl.f, nil, nil, fid)
                    async._id = 0
                end
            end
            japi._asyncExecDelay[japi._asyncExecDelayInc] = nil
        end
    end)
    --- camera lock token
    local lock = {}
    for i = 0, BJ_MAX_PLAYERS, 1 do
        local k = tostring(i)
        lock[k] = J.CreateUnit(PlayerPassive._handle, FRAMEWORK_ID["unit_token"], 0, 0, 270)
        J.HandleRef(lock[k])
    end
    --- sync同步处理
    local tgr = J.CreateTrigger()
    J.HandleRef(tgr)
    japi.DZ_TriggerRegisterSyncData(tgr, "_lk_sync", false)
    J.TriggerAddAction(tgr, function()
        local pid = 1 + J.GetPlayerId(japi.DZ_GetTriggerSyncPlayer())
        for _, v in ipairs(string.explode("||", japi.DZ_GetTriggerSyncData())) do
            local trData = string.explode("|", v)
            local k = trData[1]
            table.remove(trData, 1)
            sync.exec(k, trData, Player(pid))
        end
    end)
    --- Z轴坐标初始化
    local xMin = RegionWorld:xMin() // japi._zi
    local yMin = RegionWorld:yMin() // japi._zi
    local xMax = RegionWorld:xMax() // japi._zi
    local yMax = RegionWorld:yMax() // japi._zi
    local loc = J.Location(0, 0)
    J.HandleRef(loc)
    local x, y = xMin, 0
    while (x < xMax) do
        y = yMin
        while (y < yMax) do
            J.MoveLocation(loc, x * japi._zi, y * japi._zi)
            local z = math.ceil(J.Common["GetLocationZ"](loc))
            if (z ~= 0) then
                if (nil == japi._z[x]) then
                    japi._z[x] = {}
                end
                japi._z[x][y] = z
            end
            y = y + 1
        end
        x = x + 1
    end
    J.RemoveLocation(loc)
    J.HandleUnRef(loc)
end)