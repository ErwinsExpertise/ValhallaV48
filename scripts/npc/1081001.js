var returnMap = plr.getSavedLocation("FLORINA");
if (returnMap < 0) {
    returnMap = 104000000;
}

npc.sendNext("So you want to leave #b#m110000000##k? If you want, I can take you back to #b#m" + returnMap + "##k.");

if (npc.sendYesNo("Are you sure you want to return to #b#m" + returnMap + "##k? Alright, we'll have to get going fast. Do you want to head back to #m" + returnMap + "# now?")) {
    plr.warp(returnMap)
} else {
    npc.sendOk("You must have some business to take care of here. It's not a bad idea to take some rest at #m" + returnMap + "# Look at me; I love it here so much that I wound up living here. Hahaha anyway, talk to me when you feel like going back.")
}
