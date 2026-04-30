if (plr.getEventProperty("metalgateopen") === "yes" || map.reactorStateByName("metalgate") === 1) {
    portal.warp(990000431, "out00");
} else if (plr.haveItem(4001037, 1)) {
    plr.removeItemsByIDSilent(4001037, 1);
    plr.setEventProperty("metalgateopen", "yes");
    map.setReactorStateByName("metalgate", 1);
    portal.warp(990000431, "out00");
} else {
    portal.block("The metal gate will not open without a proper key.");
}
