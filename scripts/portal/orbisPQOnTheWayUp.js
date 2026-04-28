if (plr.isLeader() || plr.isGM()) {
    if (plr.getEventProperty("8stageclear") == null) portal.warp(920010700, "st00");
    else portal.block("You may not go back in this room.");
} else {
    portal.warp(920010700, "st00");
}
