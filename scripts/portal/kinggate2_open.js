if (plr.getEventProperty("kinggateopen") === "yes" || map.reactorStateByName("kinggate") === 1) {
    portal.warp(990000900, "st01");
} else {
    portal.block("The king's gate is sealed.");
}
