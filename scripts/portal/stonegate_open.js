if (plr.getEventProperty("stonegateopen") === "yes" || map.reactorStateByName("stonegate") === 1) {
    portal.warp(990000430, "out00");
} else if (plr.haveItem(4001026, 1)) {
    plr.removeItemsByIDSilent(4001026, 1);
    plr.setEventProperty("stonegateopen", "yes");
    map.setReactorStateByName("stonegate", 1);
    portal.warp(990000430, "out00");
} else {
    portal.block("The entrance is sealed. You need the key from the Hall of Courage.");
}
