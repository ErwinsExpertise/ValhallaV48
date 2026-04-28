var mapId = plr.mapID();

if (mapId === 250000100) {
    npc.sendSelection("Where do you want to go?\r\n#L0#Board the ship to Orbis#l\r\n#L1#Go to Herb Town#l");

    var selection = npc.selection();
    if (selection === 0) {
        plr.warp(200090310);
    } else if (selection === 1) {
        plr.warp(251000000);
    } else {
        npc.sendOk("Please choose a valid destination.");
    }
} else if (mapId === 200000141) {
    plr.warp(200090300);
} else {
    npc.sendOk("You cannot board the Mu Lung / Orbis transport from here.");
}
