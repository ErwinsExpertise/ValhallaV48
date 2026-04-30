if (plr.getEventProperty("isIncGuildPointState06") !== "1") {
    plr.gainGuildPoints(5);
    plr.setEventProperty("isIncGuildPointState06", "1");
}
plr.setEventProperty("mazeRoute:" + plr.id(), "2");
portal.warp(990000700, "st00");
