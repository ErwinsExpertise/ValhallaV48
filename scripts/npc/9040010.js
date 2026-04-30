function clearRewardPoints() {
    var start = plr.getEventProperty("entryTimestamp");
    if (start == null) {
        return 200;
    }
    var elapsedMinutes = Math.floor((Date.now() - parseInt(start, 10)) / 60000);
    if (elapsedMinutes < 25) {
        return 850;
    }
    if (elapsedMinutes < 30) {
        return 800 - ((elapsedMinutes - 25) * 20);
    }
    if (elapsedMinutes < 40) {
        return 580 - ((elapsedMinutes - 30) * 10);
    }
    if (elapsedMinutes < 60) {
        return 400 - ((elapsedMinutes - 40) * 5);
    }
    if (elapsedMinutes < 90) {
        return 260 - ((elapsedMinutes - 60) * 2);
    }
    return 200;
}

var leader = plr.getEventProperty("leader");
var bossClear = plr.getEventProperty("bossclear");

if (leader == null) {
    plr.warp(990001100);
} else if (leader !== plr.name()) {
    npc.sendOk("Ask the registered leader to speak with me.");
} else if (bossClear === true || bossClear === "true") {
    npc.sendOk("You have already reclaimed the treasure vault. Hurry and collect your reward.");
} else if (!plr.haveItem(4001024, 1)) {
    npc.sendOk("Defeat Ergoth and bring me #t4001024#. Only then can I open the treasure vault.");
} else {
    plr.removeAll(4001024);
    plr.setEventProperty("bossclear", true);
    plr.setEventProperty("completed", true);
    plr.setEventProperty("bonusStage", true);
    plr.setEventDuration("40s");

    var points = clearRewardPoints();
    var awarded = plr.getEventProperty("isIncGuildPointState08");
    if (awarded !== "1") {
        plr.gainGuildPoints(points);
        plr.setEventProperty("isIncGuildPointState08", "1");
    }

    plr.logEvent("gpq clear rewardPoints=" + points);
    map.showEffect("quest/party/clear");
    map.playSound("Party1/Clear");
    plr.warpEventMembersToPortal(990001000, "sp");
    npc.sendOk("You have defeated Ergoth and restored glory to Sharenian. I will send your guild to the treasure vault now.");
}
