var mapId = plr.mapID()

if (mapId === 680000200) {
    if (plr.weddingStage(true) === 0) {
        npc.sendOk("The guests are gathering in Saint Maple Lounge right now. Please wait awhile, the ceremony will start soon enough.")
    } else if (plr.enterWeddingAsGuest(true)) {
        npc.sendOk("Please take your seat at the Cathedral altar.")
    } else {
        npc.sendOk("The ceremony is not ready yet. Please wait a little longer.")
    }
} else if (mapId === 680000210) {
    var text = "How can I help you?#b\r\n#L0#When does the wedding begin?#l\r\n#L1#I want to leave.#l"
    npc.sendSelection(text)
    var selection = npc.selection()
    if (selection === 0) {
        npc.sendOk("We will wait until the bride and groom are ready. Please wait a few minutes.")
    } else if (selection === 1) {
        plr.warp(680000000)
    }
} else {
    npc.sendOk("I can only assist guests inside the Cathedral wedding maps.")
}
