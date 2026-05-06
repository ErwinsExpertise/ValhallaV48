var GLOVE_ITEM = 1472063;

if (plr.mapID() === 209000000) {
    if (!plr.isWearingItem(GLOVE_ITEM)) {
        npc.sendOk("The Extra Frosty Snow Zone is dangerously cold. Please equip #b#t" + GLOVE_ITEM + "##k before I send you there. If you don't have a pair yet, try opening the gift boxes around Happyville.");
    } else if (plr.itemCount(GLOVE_ITEM) > 0) {
        npc.sendOk("One pair of magical mittens is all you need. Please store or drop the extra pair before heading in.");
    } else if (npc.sendYesNo("Ready to head to the #bExtra Frosty Snow Zone#k?")) {
        plr.warp(209080000);
    } else {
        npc.sendOk("Come back when you are ready. I'll still be here.");
    }
} else if (plr.mapID() === 209080000) {
    if (npc.sendYesNo("Would you like to return to #bHappyville#k?")) {
        plr.warp(209000000);
    } else {
        npc.sendOk("Stay warm out there.");
    }
} else {
    npc.sendOk("I'm supposed to be helping out in Happyville right now.");
}
