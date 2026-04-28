if (plr.isLeader() || plr.isGM()) {
    if (plr.getEventProperty("5stageclear") == null) portal.warp(920010400, "st00");
    else portal.block("You may not go back in this room.");
} else {
    portal.warp(920010400, "st00");
}
