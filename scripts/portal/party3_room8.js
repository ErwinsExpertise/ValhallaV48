if (!plr.isLeader()) {
    if (map.playerCountInMap(920011000) > 0) {
        portal.warp(920011000, "st00");
    } else {
        portal.block("You may only enter the room your party leader is already in.");
    }
} else {
    plr.sendMessage("Your party leader entered the Room of Darkness.");
    portal.warp(920011000, "st00");
}
