var route = plr.getEventProperty("mazeRoute:" + plr.id());
if (route === "0") {
    portal.warp(990000611, "st00");
} else if (route === "1") {
    portal.warp(990000620, "st00");
} else if (route === "2") {
    portal.warp(990000631, "st00");
} else {
    portal.warp(990000641, "st00");
}
