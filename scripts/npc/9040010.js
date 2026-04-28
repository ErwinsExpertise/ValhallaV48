var leader = plr.getEventProperty("leader");

if (leader == null) {
    plr.warp(990001100);
} else if (leader !== plr.name()) {
    npc.sendOk("Ask your guild leader to talk to me.");
} else if (!plr.haveItem(4001024, 1)) {
    npc.sendOk("This is your final challenge. Defeat the evil lurking within the Rubian and return it to me. That is all.");
} else {
    plr.removeAll(4001024);
    var prev = plr.setEventProperty("bossclear", "true");
    if (prev == null && plr.inGuild()) {
        var start = plr.getEventProperty("entryTimestamp");
        if (start != null) {
            var diff = Date.now() - parseInt(start, 10);
            var points = 1000 - Math.floor(diff / (100 * 60));
            if (points < 0) {
                points = 0;
            }
            plr.gainGuildPoints(points);
        }
    }
    plr.finishEvent();
}
