if (plr.getEventProperty("isIncGuildPointState05") !== "1") {
    plr.gainGuildPoints(5);
    plr.setEventProperty("isIncGuildPointState05", "1");
}
plr.setEventProperty("mazeRoute:" + plr.id(), "1");
portal.warp(990000700, "st00");
