if (plr.getEventProperty("isIncGuildPointState07") !== "1") {
    plr.gainGuildPoints(5);
    plr.setEventProperty("isIncGuildPointState07", "1");
}
plr.setEventProperty("mazeRoute:" + plr.id(), "3");
portal.warp(990000700, "st00");
