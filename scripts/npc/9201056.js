var fee = 15000;
var mapId = plr.mapID();

if (mapId === 682000000) {
    if (!npc.sendYesNo("Would you like to return back to #bcivilization#k? The fee is " + fee + " mesos.")) {
        npc.sendOk("All right, see you next time.");
    } else if (plr.getMesos() < fee) {
        npc.sendOk("Hey, what are you trying to pull on? You don't have enough meso to pay the fee.");
    } else {
        plr.gainMesos(-fee);
        plr.warp(600000000);
    }
} else if (mapId === 600000000) {
    if (!npc.sendYesNo("Would you like to go to the #bHaunted Mansion#k? The fee is " + fee + " mesos.")) {
        npc.sendOk("All right, see you next time.");
    } else if (plr.getMesos() < fee) {
        npc.sendOk("Hey, what are you trying to pull on? You don't have enough meso to pay the fee.");
    } else {
        plr.gainMesos(-fee);
        plr.warp(682000000);
    }
} else {
    npc.sendOk("I can only take you between New Leaf City and the Haunted Mansion.");
}
