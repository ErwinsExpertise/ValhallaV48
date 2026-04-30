if (plr.getEventProperty("isIncGuildPointState04") !== "1") {
    plr.gainGuildPoints(5);
    plr.setEventProperty("isIncGuildPointState04", "1");
}
plr.setEventProperty("mazeRoute:" + plr.id(), "0");
portal.warp(990000700, "st00");
