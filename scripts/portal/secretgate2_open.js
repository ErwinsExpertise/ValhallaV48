if (plr.getEventProperty("secretgate2open") === "yes" || map.reactorStateByName("secretgate2") === 1) {
    portal.warp(990000631, "out00");
} else {
    portal.block("The hidden passage is still closed.");
}
