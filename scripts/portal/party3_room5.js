if (plr.getEventProperty("stage5_clear") === "1") {
    portal.block("There is nothing left to do in this room.");
} else if (!plr.isLeader()) {
    var count = map.playerCountInMap(920010600) + map.playerCountInMap(920010601) + map.playerCountInMap(920010602) + map.playerCountInMap(920010603) + map.playerCountInMap(920010604);
    if (count > 0) {
        portal.warp(920010600, "st00");
    } else {
        portal.block("You may only enter the room your party leader is already in.");
    }
} else {
    plr.sendMessage("Your party leader entered the Lounge.");
    portal.warp(920010600, "st00");
}
