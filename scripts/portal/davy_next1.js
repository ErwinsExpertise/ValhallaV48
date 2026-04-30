if (String(plr.getEventProperty("mobGen") || "0") !== "end") {
    portal.block("The seal on this portal has not been broken yet.");
} else {
    plr.setEventProperty("clear_2", true);
    plr.warp(925100200);
}
