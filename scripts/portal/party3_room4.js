if (plr.getEventProperty("stage4_clear") === "1") {
    portal.block("There is nothing left to do in this room.");
} else if (!plr.isLeader()) {
    if (map.playerCountInMap(920010500) > 0) {
        portal.warp(920010500, "st00");
    } else {
        portal.block("You may only enter the room your party leader is already in.");
    }
} else {
    plr.sendMessage("Your party leader entered the Sealed Room.");
    portal.warp(920010500, "st00");
}
