var openedAt = parseInt(plr.getEventProperty("s4_eTime") || "0", 10);
var remaining = plr.eventRemainingTime();
if (openedAt === 0) {
    plr.setEventProperty("s4_eTime", remaining);
    map.message(plr.name() + " has entered Lord Pirate's Servant II. You may use that portal for the next 50 seconds.");
    plr.warp(925100302);
} else if ((openedAt - remaining) < 50) {
    plr.warp(925100302);
} else {
    portal.block("That portal has already closed.");
}
