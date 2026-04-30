if (plr.getEventProperty("secretgate1open") === "yes" || map.reactorStateByName("secretgate1") === 1) {
    portal.warp(990000611, "out00");
} else {
    portal.block("The hidden passage is still closed.");
}
