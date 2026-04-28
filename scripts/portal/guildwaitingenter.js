if (plr.getEventProperty("canEnter") === false || plr.getEventProperty("canEnter") === "false") {
    portal.warp(990000100, "");
} else {
    portal.block("The portal is not open yet.");
}
