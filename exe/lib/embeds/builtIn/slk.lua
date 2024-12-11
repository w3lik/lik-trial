---@protected
function _SLK_UNIT(_v)
    _v._class = "unit"
    if (_v._parent == nil) then
        _v._parent = "nfrp"
    end
    _v.goldcost = 0
    _v.lumbercost = 0
    _v.fmade = 0
    _v.fused = 0
    _v.regenMana = 0
    _v.regenHP = 0
    _v.regenType = "none"
    _v.tilesets = "*"
    _v.sides1 = 1
    _v.dice1 = 1
    _v.def = 0
    _v.HP = 1e4
    _v.manaN = 1e4
    _v.defType = "large"
    _v.Art = _v.Art or "Framework\\ui\\default.tga"
    _v.unitSound = _v.unitSound or ""
    _v.unitShadow = _v.unitShadow or "ShadowFlyer"
    _v.scale = _v.scale or 1.00
    _v.dmgpt1 = _v.dmgpt1 or 0.3
    _v.backSw1 = _v.backSw1 or 0.3
    _v.castpt = _v.castpt or 0.1
    _v.castbsw = _v.castbsw or 0.1
    _v.targs1 = _v.targs1 or "vulnerable,ground,structure,organic,air" --攻击目标
    _v.movetp = _v.movetp or "foot"
    _v.moveHeight = _v.moveHeight or 0
    _v.spd = _v.spd or 1
    _v.turnRate = _v.turnRate or 3.0
    _v.rangeN1 = _v.rangeN1 or 100
    _v.acquire = _v.acquire or math.max(600, (_v.rangeN1 + 100))
    _v.dmgplus1 = _v.dmgplus1 or 1
    _v.cool1 = _v.cool1 or 2.0
    _v.race = _v.race or "other"
    _v.sight = _v.sight or 1500
    _v.nsight = _v.nsight or 1000
    _v.weapsOn = _v.weapsOn or 1
    _v.collision = _v.collision or 32
    if (_v.abilList == nil) then
        _v.abilList = ""
    else
        _v.abilList = _v.abilList
    end
    if (_v.weapsOn == 1) then
        _v.weapTp1 = "instant"
        _v.weapType1 = "" -- 攻击声音
        _v.Missileart_1 = "" -- 箭矢模型
        _v.Missilespeed_1 = 99999 -- 箭矢速度
        _v.Missilearc_1 = 0
    else
        _v.weapTp1 = ""
        _v.weapType1 = ""
        _v.Missileart_1 = ""
        _v.Missilespeed_1 = 0
        _v.Missilearc_1 = 0
    end
    return _v
end

---@protected
function _SLK_ABILITY(_v)
    _v._class = "ability"
    if (_v._parent == nil) then
        _v._parent = "ANcl"
    end
    return _v
end

---@protected
function _SLK_DESTRUCTABLE(_v)
    _v._class = "destructable"
    if (_v._parent == nil) then
        _v._parent = "BTsc"
    end
    _v.HP = _v.HP or 1e4
    _v.fogVis = _v.fogVis or 0
    _v.showInMM = _v.showInMM or 0
    _v.selectable = _v.selectable or 0
    _v.selSize = _v.selSize or 0
    _v.selcircsize = _v.selcircsize or 128
    return _v
end
