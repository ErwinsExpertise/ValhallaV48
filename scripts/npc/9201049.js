var chapelStage = plr.weddingStage(false)
var cathedralStage = plr.weddingStage(true)
var activeWedding = chapelStage >= 0 || cathedralStage >= 0

if (!activeWedding) {
    npc.sendOk("Hey there, did you enjoy the wedding? I will head you back to #bAmoria#k now.")
    plr.warp(680000000)
} else {
    var result = plr.claimWeddingExitReward()
    if (result === 3) {
        npc.sendOk("Please make room in your ETC inventory to receive your #b#t4031424##k first.")
    } else if (result === 2) {
        npc.sendOk("You just received an #b#t4031424##k. Search for #b#p9201014##k at the top of Amoria if you want to open it.")
        plr.warp(680000000)
    } else {
        npc.sendOk("I hope you enjoyed the wedding. I will send you back to #bAmoria#k now.")
        plr.warp(680000000)
    }
}
