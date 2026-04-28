if (npc.sendYesNo("Do you want to leave #bHappyville#k?")) {
    var location = plr.getSavedLocation("HV_MAP");
    if (location < 0) {
        location = 100000000;
    }
    plr.warp(location);
}
