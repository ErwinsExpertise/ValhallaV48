var mapId = plr.mapID();

if (mapId === 670010200) {
    if (!plr.isLeader()) {
        npc.sendOk("Ask your party leader to talk to me.");
    } else if (plr.getEventProperty("apqStage1Clear")) {
        npc.sendOk("You've shattered the mirror. Head back up and speak with Amos.");
    } else if (plr.haveItem(4031595, 1)) {
        plr.removeItemsByID(4031595, 1);
        plr.setEventProperty("apqStage1Clear", true);
        plr.partyGiveExp(2000);
        map.showEffect("quest/party/clear");
        map.playSound("Party1/Clear");
        npc.sendOk("You've shattered the Magik Mirror. Go see Amos and he'll lead your party onward.");
    } else if (!plr.getEventProperty("apqStage1Boss") && plr.countMonster() === 0 && !plr.haveItem(4031596, 1)) {
        plr.setEventProperty("apqStage1Boss", true);
        plr.spawnMonster(9400518, 1000, 1340);
        plr.sendMessage("A special Magik Fierry has appeared somewhere in the map.");
        npc.sendOk("You've cleared the monsters the mirror summoned. One final Magik Fierry has appeared. Defeat it, claim the #b#t4031596##k, and drop it onto the mirror.");
    } else {
        npc.sendOk("Break the mirror with a #b#t4031596##k and bring me back one of its shattered pieces.");
    }
} else if (mapId === 670011000) {
    plr.warp(670010000);
} else {
    npc.sendOk("There's nothing for me to do here.");
}
