var returnMap = 100000000;

if (plr.mapID() === 108010301) {
    returnMap = 102000000;
} else if (plr.mapID() === 108010201) {
    returnMap = 101000000;
} else if (plr.mapID() === 108010101) {
    returnMap = 100000000;
} else if (plr.mapID() === 108010401) {
    returnMap = 103000000;
}

if (npc.sendYesNo("You can use the Sparkling Crystal to return to the real world. Do you want to go back now?")) {
    plr.warp(returnMap);
}
