var openedAt = parseInt(plr.getEventProperty("s3_eTime") || "0", 10);
var remaining = plr.eventRemainingTime();
if (openedAt === 0) {
    plr.setEventProperty("s3_eTime", remaining);
    map.message(plr.name() + " has entered Lord Pirate's Servant I. You may use that portal for the next 50 seconds.");
    plr.warp(925100202);
} else if ((openedAt - remaining) < 50) {
    plr.warp(925100202);
} else {
    portal.block("That portal has already closed.");
}
