if (plr.getEventProperty("watergateopen") === "yes" || map.reactorStateByName("watergate") === 1) {
    if (plr.getEventProperty("isIncGuildPointState03") !== "1") {
        plr.gainGuildPoints(25);
        plr.setEventProperty("isIncGuildPointState03", "1");
    }
    portal.warp(990000600, "st00");
} else {
    portal.block("This way forward is not open yet.");
}
