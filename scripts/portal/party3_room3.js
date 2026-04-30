if (plr.getEventProperty("stage3_clear") === "1") {
    portal.block("There is nothing left to do in this room.");
} else if (!plr.isLeader()) {
    if (map.playerCountInMap(920010400) > 0) {
        portal.warp(920010400, "st00");
    } else {
        portal.block("You may only enter the room your party leader is already in.");
    }
} else {
    plr.sendMessage("Your party leader entered the Lobby.");
    portal.warp(920010400, "st00");
}
