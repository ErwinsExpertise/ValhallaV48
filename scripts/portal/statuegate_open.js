if (plr.getEventProperty("statuegateopen") === "yes" || map.reactorStateByName("statuegate") === 1) {
    if (plr.haveItem(1032033, 2)) {
        portal.block("The Protector Rock rejects your greed. You may only carry one.");
    } else {
        if (plr.getEventProperty("isIncGuildPointState01") !== "1") {
            plr.gainGuildPoints(15);
            plr.setEventProperty("isIncGuildPointState01", "1");
        }
        portal.warp(990000301, "st00");
    }
} else {
    portal.block("The gate is closed.");
}
