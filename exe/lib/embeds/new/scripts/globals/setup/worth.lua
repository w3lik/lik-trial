-- 资源设定
Game():worth("lumber", "木头", { "gold", 1000000 }) -- 1木 = 1000000金
Game():worth("gold", "黄金", { "silver", 100 }) -- 1金 = 100银
Game():worth("silver", "白银", { "copper", 100 }) -- 1银 = 100铜
Game():worth("copper", "青铜") -- 无下级
Game():worth("other", "其他") -- 无下级