var mapId = plr.mapID()

if (mapId === 680000300) {
    var text = "What would you like to do?#b\r\n#L0#Leave for Amoria.#l\r\n#L1#Stay a bit longer.#l"
    npc.sendSelection(text)
    if (npc.selection() === 0) {
        plr.warp(680000000)
    }
} else {
    if (plr.haveItem(4000313, 1) && plr.isMarried()) {
        npc.sendOk("Enjoy your wedding celebration. When you are ready, the afterparty awaits.")
    } else if (npc.sendYesNo("Are you sure you want to leave and return to Amoria?")) {
        plr.warp(680000000)
    } else {
        npc.sendOk("Please enjoy the wedding festivities.")
    }
}
