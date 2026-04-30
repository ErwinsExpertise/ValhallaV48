if (plr.getEventProperty("stage1_clear") === "1") {
    portal.block("There is nothing left to do in this room.");
} else if (!plr.isLeader()) {
    if (map.playerCountInMap(920010200) > 0) {
        portal.warp(920010200, "st00");
    } else {
        portal.block("You may only enter the room your party leader is already in.");
    }
} else {
    plr.sendMessage("Your party leader entered the Walkway.");
    portal.warp(920010200, "st00");
}
