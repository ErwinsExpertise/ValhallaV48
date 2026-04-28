if (!npc.sendYesNo(plr.mapID() === 100000000
    ? "I can take you to the Amoria Village. Are you ready to go?"
    : "I can take you back to Henesys. Are you ready to go?")) {
    npc.sendOk("Ok, feel free to hang around until you're ready to go!");
} else if (plr.mapID() === 100000000) {
    plr.warp(680000000);
} else {
    plr.warp(100000000);
}
