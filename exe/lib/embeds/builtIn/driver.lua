--- 分隔字符串
---@param str string
---@param size number 每隔[size]字符切一次
---@return string
function string.split(str, size)
    local sp = {}
    local len = string.len(str)
    if (len <= 0) then
        return sp
    end
    size = size or 1
    local i = 1
    while (i <= len) do
        table.insert(sp, string.sub(str, i, i + size - 1))
        i = i + size
    end
    return sp
end

--- 把字符串以分隔符打散为数组
---@param delimeter string
---@param str string
---@return string[]
function string.explode(delimeter, str)
    if (delimeter == '') then
        return string.split(str, 1)
    end
    local res = {}
    local s, a, b = 1, 1, 1
    while true do
        a, b = string.find(str, delimeter, s, true)
        if not a then
            break
        end
        table.insert(res, string.sub(str, s, a - 1))
        s = b + 1
    end
    table.insert(res, string.sub(str, s))
    return res
end

--- 移除字符串两侧的空白字符或其他预定义字符
---@param str string
---@return string
function string.trim(str)
    local res = string.gsub(str, "^%s*(.-)%s*$", "%1")
    return res
end

--- 在数组内
---@param arr table
---@param val any
---@return boolean
function table.includes(arr, val)
    local isIn = false
    if (val == nil or #arr <= 0) then
        return isIn
    end
    for _, v in ipairs(arr) do
        if (v == val) then
            isIn = true
            break
        end
    end
    return isIn
end

--- 合并table
---@vararg table
---@return table
function table.merge(...)
    local tempTable = {}
    local tables = { ... }
    if (tables == nil) then
        return {}
    end
    for _, tn in ipairs(tables) do
        if (type(tn) == "table") then
            for k, v in pairs(tn) do
                tempTable[k] = v
            end
        end
    end
    return tempTable
end

--- 四舍五入取整
---@param decimal number
---@return number integer
function math.round(decimal)
    return math.floor(decimal + 0.5)
end