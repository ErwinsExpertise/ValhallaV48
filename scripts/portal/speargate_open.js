if (plr.getEventProperty("speargateopen") === "yes" || map.reactorStateByName("speargate") === 4) {
    if (plr.getEventProperty("isIncGuildPointState02") !== "1") {
        plr.gainGuildPoints(20);
        plr.setEventProperty("isIncGuildPointState02", "1");
    }
    plr.setEventProperty("speargateopen", "yes");
    portal.warp(990000401, "st00");
} else {
    portal.block("This way forward is not open yet.");
}
