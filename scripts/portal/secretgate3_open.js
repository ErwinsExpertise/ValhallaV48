if (plr.getEventProperty("secretgate3open") === "yes" || map.reactorStateByName("secretgate3") === 1) {
    portal.warp(990000641, "out00");
} else {
    portal.block("The hidden passage is still closed.");
}
