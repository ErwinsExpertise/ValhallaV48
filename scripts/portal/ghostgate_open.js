if (plr.getEventProperty("ghostgateopen") === "yes" || map.reactorStateByName("ghostgate") === 1) {
    if (plr.getEventProperty("isIncGuildPointState09") !== "1") {
        plr.gainGuildPoints(10);
        plr.setEventProperty("isIncGuildPointState09", "1");
    }
    portal.warp(990000800, "st00");
} else {
    portal.block("The entrance to the throne remains shut.");
}
