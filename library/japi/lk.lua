--[[
    LIK自实现方法
    方法没有特定开头
    lik已接管DzFrameSetUpdateCallbackByCode方法，此方法不可私自使用
]]

--- 数据配置
japi._asyncExecDelay = japi._asyncExecDelay or {}
japi._asyncExecDelayId = japi._asyncExecDelayId or 0
japi._asyncExecDelayInc = japi._asyncExecDelayInc or 0
japi._asyncRefresh = japi._asyncRefresh or {}
japi._blackBordersBottom = japi._blackBordersBottom or 0.130
japi._blackBordersInner = japi._blackBordersInner or 0.45
japi._blackBordersTop = japi._blackBordersTop or 0.020
japi._clientHeight = japi._clientHeight or 0
japi._clientWidth = japi._clientWidth or 0
japi._cursor = japi._cursor or nil
japi._cursorLast = japi._cursorLast or nil
japi._hasMallItem = japi._hasMallItem or {}
japi._isWideScreen = datum.default(japi._isWideScreen, false)
japi._keyboard = japi._keyboard or { press = {}, release = {} }
japi._loadToc = japi._loadToc or {}
japi._roulette = japi._roulette or nil
japi._rouletteWait = datum.default(japi._rouletteWait, false)
japi._rouletteWaitTimer = japi._rouletteWaitTimer or nil
japi._frAdaptive = japi._frAdaptive or nil
japi._frEsc = japi._frEsc or nil
japi._frTagIndex = japi._frTagIndex or 0
japi._z = japi._z or {}
japi._zi = 64

--- 使用宽屏模式
--- 地图可以根据自身特点，强制打开或关闭的宽屏优化支持功能。
--- 开启宽屏模式可以解决单位被拉伸显得比较“胖”的问题。
---@param enable boolean
---@return void
function japi.EnableWideScreen(enable)
    japi._isWideScreen = enable
    japi.DZ_EnableWideScreen(enable)
end

--- 是否宽屏模式
---@return boolean
function japi.IsWideScreen()
    return japi._isWideScreen
end

--- 获取魔兽客户端宽度
--- 不包括魔兽窗口边框
---@return number
function japi.GetClientWidth()
    return japi._clientWidth
end

--- 获取魔兽客户端高度
--- 不包括魔兽窗口边框
---@return number
function japi.GetClientHeight()
    return japi._clientHeight
end

--- [别名]DzFrameEditBlackBorders
--- 修改游戏渲染黑边: 上方高度:topHeight,下方高度:bottomHeight
--- 上下加起来不要大于0.6
---@param topHeight number
---@param bottomHeight number
---@return void
function japi.SetBlackBorders(topHeight, bottomHeight)
    japi._blackBordersTop = topHeight
    japi._blackBordersBottom = bottomHeight
    japi._blackBordersInner = 0.6 - topHeight - bottomHeight
    japi.DZ_FrameEditBlackBorders(topHeight, bottomHeight)
end

--- 获得游戏渲染的：离顶黑边高、离底黑边高、中间显示高、
---@return number,number,number top,bottom,inner
function japi.GetBlackBorders()
    return japi._blackBordersTop, japi._blackBordersBottom, japi._blackBordersInner
end

--- 玩家[本地调试环境下]是否拥有该商城道具（平台地图商城）
--- 平台地图商城玩家拥有该道具返还true
---@param whichPlayer number
---@param key string
---@return boolean
function japi.HasMallItem(whichPlayer, key)
    if (DEBUGGING) then
        return true == japi._hasMallItem[key]
    end
    return japi.DZ_Map_HasMallItem(whichPlayer, key)
end

--- 强制[本地调试环境下]所有玩家拥有该商城道具
---@vararg string 支持多个字符串keys
---@return void
function japi.SetMallItem(...)
    if (DEBUGGING) then
        for _, k in ipairs({ ... }) do
            japi._hasMallItem[k] = true
        end
    end
end

--- 新建一个Frame的Tag索引
---@return string
function japi.FrameTagIndex()
    japi._frTagIndex = japi._frTagIndex + 1
    return "Frame#" .. japi._frTagIndex
end

--- 执行自适应Frame大小
--- 以流行尺寸作为基准比例，以高为基准结合魔兽4:3计算自动调节宽度的自适应规则
---@param w number 宽
---@return number
function japi.FrameAdaptive(w)
    w = w or 0
    if (w == 0) then
        return 0
    end
    local sr = 4 / 3
    local pr = 16 / 9
    local tr = sr / pr
    local dr = japi._clientWidth / japi._clientHeight / pr
    w = w * tr / dr
    if (w > 0) then
        w = math.max(0.0002, w)
        w = math.min(0.8, w)
    elseif (w < 0) then
        w = math.max(-0.8, w)
        w = math.min(-0.0002, w)
    end
    return w
end

--- 执行自适应Frame大小反算
--- 以流行尺寸作为基准比例，以高为基准结合魔兽4:3计算自动调节宽度的自适应规则
---@param w number 宽
---@return number
function japi.FrameDisAdaptive(w)
    w = w or 0
    if (w == 0) then
        return 0
    end
    local sr = 4 / 3
    local pr = 16 / 9
    local tr = sr / pr
    local dr = japi._clientWidth / japi._clientHeight / pr
    w = w * dr / tr
    if (w > 0) then
        w = math.max(0.0002, w)
        w = math.min(1.6, w)
    elseif (w < 0) then
        w = math.max(-1.6, w)
        w = math.min(-0.0002, w)
    end
    return w
end

--- 注册Frame对象自适应处理
---@param key string
---@param fr Frame
---@return void
function japi.FrameSetAdaptive(key, fr)
    japi._frAdaptive:set(key, fr)
end

--- 注册Frame对象Esc叠层处理
---@param key string
---@param fr Frame
---@return void
function japi.FrameSetEsc(key, fr)
    japi._frEsc:set(key, fr)
end

--- 加载Toc文件列表
--- 加载--> *.toc
--- 载入自己的fdf列表文件
---@return void
function japi.LoadToc(tocFilePath)
    if (japi._loadToc[tocFilePath] == true) then
        return true
    end
    japi._loadToc[tocFilePath] = true
    japi.DZ_LoadToc(tocFilePath)
end

--- 获取某个坐标的Z轴高度
---@param x number
---@param y number
---@return number
function japi.Z(x, y)
    if (type(x) == "number" and type(y) == "number") then
        local zi = japi._zi
        local x0, y0 = math.round(x / zi), math.round(y / zi)
        local x1, y1 = x0 + 1, y0 + 1
        if (nil == japi._z[x0]) then
            if (nil == japi._z[x1]) then
                return 0
            end
            return japi._z[x1][y0] or 0
        end
        if (nil == japi._z[x1]) then
            return japi._z[x0][y0] or 0
        end
        local sel = {
            { japi._z[x0][y0] or 0, vector2.distance(x, y, x0 * zi, y0 * zi) },
            { japi._z[x0][y1] or 0, vector2.distance(x, y, x0 * zi, y1 * zi) },
            { datum.ternary(nil == japi._z[x1], 0, japi._z[x1][y0] or 0), vector2.distance(x, y, x1 * zi, y0 * zi) },
            { datum.ternary(nil == japi._z[x1], 0, japi._z[x1][y1] or 0), vector2.distance(x, y, x1 * zi, y1 * zi) },
        }
        local o = sel[1]
        for _, v in ipairs(sel) do
            if (v[2] < o[2]) then
                o = v
            end
        end
        return o[1]
    end
    return 0
end

--- 设置镜头属性
---@param state any
---@param value number
---@param min number 下限值
---@param max number 上限值
---@return void
function japi.CameraSetField(state, value, min, max)
    if (type(min) == "number") then
        value = math.max(min, value)
    end
    if (type(max) == "number") then
        value = math.min(max, value)
    end
    J.SetCameraField(state, value, 0)
end

--- X比例 转 像素
---@param x number
---@return number
function japi.PX(x)
    return japi._clientWidth * x / 0.8
end

--- Y比例 转 像素
---@param y number
---@return number
function japi.PY(y)
    return japi._clientHeight * y / 0.6
end

--- X像素 转 比例
---@param x number
---@return number
function japi.RX(x)
    return x / japi._clientWidth * 0.8
end

--- Y像素 转 比例
---@param y number
---@return number
function japi.RY(y)
    return y / japi._clientHeight * 0.6
end

--- 鼠标X像素 转 比例
---@return number
function japi.MouseRX()
    if (type(japi._cursor) == "table") then
        return japi._cursor.rx
    end
    return japi.RX(japi.DZ_GetMouseXRelative())
end

--- 鼠标Y像素 转 比例
---@return number
function japi.MouseRY()
    if (type(japi._cursor) == "table") then
        return japi._cursor.ry
    end
    return japi.RY(japi._clientHeight - japi.DZ_GetMouseYRelative())
end

--- 异步刷新（不区分玩家默认一直刷新下去）
---@param key string 标识键
---@param callFunc fun():void
---@return void
function japi.AsyncRefresh(key, callFunc)
    if (type(callFunc) ~= "function") then
        japi._asyncRefresh[key] = callFunc
    else
        japi._asyncRefresh[key] = callFunc
    end
end

--- 多少帧后异步执行（区分玩家只执行一次）
---@param frame number 刷帧数，多少帧后执行，默认1
---@param playerIndex number integer 特定玩家(索引)才生效
---@param callFunc fun(execId:number):void
---@return number 返回一个刷新Id，以此Id可取消本次刷新操作
function japi.AsyncExecDelay(frame, playerIndex, callFunc)
    if (type(callFunc) ~= "function") then
        callFunc = nil
    end
    if (type(playerIndex) ~= "number") then
        playerIndex = PlayerLocal():index()
    end
    frame = math.max(1, math.round(frame))
    local inc = japi._asyncExecDelayInc + frame
    if (japi._asyncExecDelay[inc] == nil) then
        japi._asyncExecDelay[inc] = {}
    end
    japi._asyncExecDelayId = japi._asyncExecDelayId + 1
    japi._asyncExecDelay[inc][japi._asyncExecDelayId] = { i = playerIndex, f = callFunc }
    return inc .. '#' .. japi._asyncExecDelayId
end

--- 取消帧后异步执行
---@param execId number
---@return void
function japi.CancelAsyncExecDelay(execId)
    local ids = string.explode('#', execId)
    local inc = math.round(ids[1])
    local id = math.round(ids[2])
    if (japi._asyncExecDelay[inc] and japi._asyncExecDelay[inc]) then
        japi._asyncExecDelay[inc][id] = nil
    end
end

--- 轮盘队列
--- 此方法自带延迟策略，并且自动合并请求
--- 从而可以大大减轻执行压力
--- 只适用于无返回执行
---@param whichPlayer number
---@param key string
---@param func function
---@return void
function japi.Roulette(func, whichPlayer, key, value)
    sync.must()
    if (type(func) ~= "function" or type(key) ~= "string" or type(value) ~= "string") then
        return
    end
    local rf = function()
        if (japi.DZ_IsServerAlready(whichPlayer)) then
            func(whichPlayer, key, value)
        end
    end
    if (isClass(japi._roulette, ArrayClass)) then
        japi._rouletteWait = false
        destroy(japi._rouletteWaitTimer)
        japi._rouletteWaitTimer = nil
        japi._roulette:set(key, rf)
        return
    end
    japi._roulette = Array()
    japi._roulette:set(key, rf)
    time.setInterval(0, function(curTimer)
        curTimer:period(5)
        local ks = japi._roulette:keys()
        local ksl = #ks
        if (ksl > 0) then
            local k1 = ks[1]
            local f = japi._roulette:get(k1)
            f()
            japi._roulette:set(k1, nil)
        end
        if (ksl == 0 or japi._roulette:count() == 0) then
            japi._rouletteWait = true
            japi._rouletteWaitTimer = time.setTimeout(4.99, function()
                japi._rouletteWaitTimer = nil
                if (japi._rouletteWait == true) then
                    destroy(curTimer)
                    japi._roulette = nil
                    japi._rouletteWait = false
                end
            end)
        end
    end)
end

--- 保存服务器存档
--- 会根据数据类型自动添加前缀
---@param whichPlayer number
---@param key string
---@param value string
function japi.ServerSaveValue(whichPlayer, key, value)
    if (string.len(key) > 63) then
        japi.Tips("63KeyTooLong")
        return
    end
    if (type(value) == "boolean") then
        if (value == true) then
            value = "B:1"
        else
            value = "B:0"
        end
    elseif (type(value) == "number") then
        value = "N:" .. tostring(value)
    elseif (type(value) ~= "string") then
        value = ""
    end
    if (string.len(value) > 63) then
        japi.Tips("63ValueTooLong")
        return
    end
    japi.Roulette(japi.DZ_Map_SaveServerValue, whichPlayer, key, value)
end

--- 获取服务器存档
--- 会处理根据数据类型自动添加前缀的数据
---@param whichPlayer number
---@param key string
---@return any
function japi.ServerLoadValue(whichPlayer, key)
    if (string.len(key) > 63) then
        japi.Tips("63KeyTooLong")
        return
    end
    if (japi.DZ_IsServerAlready(whichPlayer)) then
        local result = japi.DZ_Map_GetServerValue(whichPlayer, key)
        if (type(result) == "string") then
            local valType = string.sub(result, 1, 2)
            if (valType == "B:") then
                local v = string.sub(result, 3)
                return "1" == v
            elseif (valType == "N:") then
                local v = string.sub(result, 3)
                return tonumber(v or 0)
            end
            if (result == '') then
                return nil
            end
            return result
        end
    end
    return nil
end

--- 设置房间显示的数据
--- 为服务器存档显示的数据，对应作者之家的房间key
---@param whichPlayer number
---@param key string
---@param value string
function japi.ServerSaveRoom(whichPlayer, key, value)
    if (string.len(key) > 63) then
        japi.Tips("63KeyTooLong")
        return
    end
    key = string.upper(key)
    if (type(value) == "boolean") then
        if (value == true) then
            value = "true"
        else
            value = "false"
        end
    elseif (type(value) == "number") then
        value = math.numberFormat(value, 2)
    elseif (type(value) ~= "string") then
        value = ""
    end
    if (string.len(value) > 63) then
        japi.Tips("63ValueTooLong")
        return
    end
    japi.Roulette(japi.DZ_Map_Stat_SetStat, whichPlayer, key, value)
end