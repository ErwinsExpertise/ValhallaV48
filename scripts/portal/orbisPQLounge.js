if (plr.isLeader() || plr.isGM()) {
    if (plr.getEventProperty("7stageclear") == null) portal.warp(920010600, "st00");
    else portal.block("You may not go back in this room.");
} else {
    portal.warp(920010600, "st00");
}
