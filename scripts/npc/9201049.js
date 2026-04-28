var coupleBox = 4031424
var guestBox = 4031423

if (!plr.haveItem(4000313, 1)) {
    npc.sendOk("Hey there, did you enjoy the wedding? I will head you back to #bAmoria#k now.")
    plr.warp(680000000)
} else {
    var reward = plr.isMarried() ? coupleBox : guestBox
    if (!plr.canHold(reward, 1)) {
        npc.sendOk("Please make room in your ETC inventory to receive your Onyx Chest.")
    } else {
        plr.gainItem(4000313, -1)
        plr.gainItem(reward, 1)
        npc.sendOk("You just received an Onyx Chest. Search for #b#p9201014##k at the top of Amoria if you want to open it.")
        plr.warp(680000000)
    }
}
