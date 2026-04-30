var reqItem = 4000001
var reqQty = 50
var reward = 4031367

if (plr.haveItem(reward, 1)) {
    npc.sendOk("You already earned my #bProof of Love#k. Treasure it carefully.")
} else if (!plr.haveItem(reqItem, reqQty)) {
    npc.sendOk("If you want my #bProof of Love#k, bring me #b" + reqQty + " #t" + reqItem + "##k first.")
} else if (!plr.canHold(reward, 1)) {
    npc.sendOk("Please free an ETC slot before I hand over your #bProof of Love#k.")
} else if (!npc.sendYesNo("You brought everything I asked for. Would you like to trade it for my #bProof of Love#k?")) {
    npc.sendOk("Come back when you are ready.")
} else {
    plr.gainItem(reqItem, -reqQty)
    plr.gainItem(reward, 1)
    npc.sendOk("Wonderful. Here is my #bProof of Love#k.")
}
