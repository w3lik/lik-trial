local isCustom = function(trunk)
    local l = #trunk - 3
    if (l < 2) then
        return false
    end
    local s = string.sub(trunk, l, l)
    return s == '.'
end

---@private
---@param _v noteSlkUnit
function slk_unit(_v)
    _v = _SLK_UNIT(_v)
    _v._id = SLK_ID(_v)
    SLK_GO_SET(_v)
    return _v
end

---@private
---@param _v{checkDep,Requires,Requiresamount,Effectsound,Effectsoundlooped,EditorSuffix,Name,Untip,Unubertip,Tip,Ubertip,Researchtip,Researchubertip,Unorder,Orderon,Order,Orderoff,Unhotkey,Hotkey,Researchhotkey,UnButtonpos_1,UnButtonpos_2,Buttonpos_1,Buttonpos_2,Researchbuttonpos1,Researchbuttonpos2,Unart,Researchart,Art,SpecialArt,Specialattach,Missileart_1,Missilespeed_1,Missilearc_1,MissileHoming_1,LightningEffect,EffectArt,TargetArt,Targetattachcount,Targetattach,Targetattach1,Targetattach2,Targetattach3,Targetattach4,Targetattach5,Areaeffectart,Animnames,CasterArt,Casterattachcount,Casterattach,Casterattach1,hero,item,race,levels,reqLevel,priority,BuffID,EfctID,Tip,Ubertip,targs,DataA,DataB,DataC,DataD,DataE,DataF,Cast,Cool,Dur,HeroDur,Cost,Rng,Area,_id_force,_class,_type,_parent,_desc,_attr,_ring,_remarks,_lv,_onSkillStudy,_onSkillEffect,_onRing}
function slk_ability(_v)
    _v = _SLK_ABILITY(_v)
    _v._id = SLK_ID(_v)
    SLK_GO_SET(_v)
    return _v
end

---@private
---@param _v{file,}
function slk_destructable(_v)
    _v = _SLK_DESTRUCTABLE(_v)
    _v._id = SLK_ID(_v)
    SLK_GO_SET(_v)
    return _v
end

function assets_path2Alias(typ, path)
    if (typ) then
        if (path) then
            local k = typ .. path
            return GO_CHECK.path[k]
        end
    end
end

function assets_alias2Path(typ, alias)
    if (typ) then
        if (alias) then
            local k = typ .. alias
            return GO_CHECK.alias[k]
        end
    end
end

function assets_check(typ, path, alias)
    if (typ) then
        if (isCustom(path) and (nil == alias or alias == "")) then
            local k = path .. ':' .. alias
            error("Unreasonable path: <" .. typ .. ">" .. k)
        end
        if (path) then
            local k = typ .. path
            if (GO_CHECK.path[k] ~= nil) then
                print("Duplicate path: <" .. typ .. ">" .. k)
            end
            GO_CHECK.path[k] = alias
        end
        if (alias) then
            local k = typ .. alias
            if (GO_CHECK.alias[k] ~= nil) then
                print("Duplicate alias: <" .. typ .. ">" .. k)
            end
            GO_CHECK.alias[k] = path
        end
    end
end

function assets_selection(typeName)
    if (type(typeName) == "string") then
        if (isCustom(typeName)) then
            print("selection must not have a suffix")
        else
            GO_RESULT.selection = string.trim(typeName)
        end
    end
end

function assets_font(ttfName)
    if (type(ttfName) == "string") then
        if (isCustom(ttfName)) then
            print("font must not have a suffix")
        else
            GO_RESULT.font = string.trim(ttfName)
        end
    end
end

function assets_loading(tga)
    if (type(tga) == "string") then
        if (isCustom(tga)) then
            print("loading must not have a suffix")
        else
            GO_RESULT.loading = string.trim(tga)
        end
    end
end

function assets_preview(tga)
    if (type(tga) == "string") then
        if (isCustom(tga)) then
            print("preview must not have a suffix")
        else
            GO_RESULT.preview = string.trim(tga)
        end
    end
end

function assets_icon(path, alias)
    if (path) then
        path = string.trim(path)
        alias = alias or path
        assets_check("assets_icon", path, alias)
        if (false == isCustom(path)) then
            table.insert(GO_RESULT.icons, { path, alias })
        end
    end
end

function assets_model(path, alias)
    if (path) then
        path = string.trim(path)
        alias = alias or path
        assets_check("assets_model", path, alias)
        if (false == isCustom(path)) then
            table.insert(GO_RESULT.model, { path, alias })
        end
    end
end

function assets_sound(path, alias, soundType)
    if (path) then
        path = string.trim(path)
        alias = alias or path
        assets_check("assets_sound", path, alias)
        if (false == isCustom(path)) then
            table.insert(GO_RESULT.sound, { path, alias, soundType })
        end
    end
end

function assets_plugins(kit)
    if (kit) then
        assets_check("assets_plugins", kit)
        table.insert(GO_RESULT.plugins, kit)
    end
end

function assets_ui(kit)
    if (kit) then
        assets_check("assets_ui", kit)
        table.insert(GO_RESULT.ui, kit)
    end
end

function assets_speech_extra(_v)
    return _v
end

function assets_speech(name, extra)
    if (name ~= nil) then
        name = string.trim(name)
        assets_check("assets_speech", name, nil)
        local v = {
            weapsOn = 1,
            file = ".mdl",
            unitSound = name,
            collision = 32,
            Name = name .. "|D",
            abilList = "AInv",
            canFlee = 1,
        }
        slk_unit(v)
        if (type(extra) == "table") then
            for k, e in pairx(extra) do
                if (type(e) == "table") then
                    local ve = {}
                    for k2, v2 in pairs(v) do
                        ve[k2] = v2
                    end
                    for k2, v2 in pairs(e) do
                        ve[k2] = v2
                    end
                    ve.unitSound = name
                    ve.Name = name .. "|EX|" .. k
                    ve.canFlee = ve.canFlee or 1
                    local abilList = { "AInv" }
                    if (type(ve.abilList) == "string") then
                        for _, a in ipairs(string.explode(',', ve.abilList)) do
                            abilList[#abilList + 1] = a
                        end
                    end
                    ve.abilList = table.concat(abilList, ',')
                    if (ve.model) then
                        ve.file = assets_alias2Path("assets_model", ve.model)
                        if (ve.file == nil) then
                            error("model【" .. ve.model .. "】 not exist")
                        end
                        if (false == isCustom(ve.file)) then
                            local a = string.gsub(ve.model, "/", "\\", -1)
                            ve.file = "war3mapModel\\" .. a .. ".mdl"
                        end
                        ve.model = nil
                    end
                    ve.file = ve.file or ".mdl"
                    slk_unit(ve)
                end
            end
        end
    end
end