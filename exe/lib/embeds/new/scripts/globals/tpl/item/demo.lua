---@param getData noteOnItemGetData
TPL_ITEM.DEMO = ItemTpl()
    :modelAlias("TreasureChest")
    :name("勇气之剑")
    :description("就是一把剑")
    :ability(TPL_ABILITY.DEMO)
    :icon("ReplaceableTextures\\CommandButtons\\BTNArcaniteMelee.blp")
    :worth({ gold = 10 })
    :onEvent(EVENT.Item.Get,
    function(getData)
        echo(getData.triggerUnit:name() .. " 获得物品 " .. getData.triggerItem:name())
    end)