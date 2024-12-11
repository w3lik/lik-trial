local _____ = {}
function ___(w)
    if (_____[w] == nil) then
        local n = {}
        local st, sp, ep, en = 1, 1, 1, 8
        local sl = string.sub(w, 1, en)
        local sw = string.sub(w, en + 1)
        while true do
            sp, ep = string.find(sw, sl, st, true)
            if not sp then
                break
            end
            table.insert(n, string.sub(sw, st, sp - 1))
            st = ep + 1
        end
        table.insert(n, string.sub(sw, st))
        local s = ''
        for i = #n, 1, -1 do
            local nn = string.reverse(n[i])
            while true do
                if (____[nn] == nil) then
                    s = s .. nn
                    break
                else
                    nn = ____[nn]
                end
            end
        end
        _____[w] = s
    end
    return _____[w]
end
--FRAMEWORK_SCRIPT_IN
FRAMEWORK_INIT()
local t = J.CreateTrigger()
J.HandleRef(t)
J.TriggerRegisterTimerEvent(t, 0.1, false)
J.TriggerAddAction(t, function()
    J.DisableTrigger(t)
    J.DestroyTrigger(t)
    J.HandleUnRef(t)
    t = nil
    event.syncTrigger(Game(), EVENT.Game.Start)
    event.syncUnregister(Game(), EVENT.Game.Start)
    local p1st = Process("start")
    if (isClass(p1st, ProcessClass)) then
        p1st:start()
    end
end)