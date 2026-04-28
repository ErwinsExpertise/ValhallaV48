var mapId = plr.mapID()

if (mapId === 680000100) {
    if (plr.weddingStage(false) === 0) {
        npc.sendOk("The guests are gathering in White Wedding Lounge right now. Please wait awhile, the ceremony will start soon enough.")
    } else if (plr.enterWeddingAsGuest(false)) {
        npc.sendOk("Please take your seat at the Chapel altar.")
    } else {
        npc.sendOk("The ceremony is not ready yet. Please wait a little longer.")
    }
} else if (mapId === 680000110) {
    if (!plr.isMarried()) {
        npc.sendOk("The superstars must receive Pelvis Bebop's word before heading to the afterparty.")
    } else {
        var text = "What can I help you with?#b\r\n#L0#Go to the afterparty.#l\r\n#L1#What should I be doing?#l"
        npc.sendSelection(text)
        var selection = npc.selection()
        if (selection === 0) {
            if (plr.startWeddingAfterParty()) {
                npc.sendOk("Enjoy! Cherish your wedding memories forever!")
            } else {
                npc.sendOk("I can't begin the afterparty just yet.")
            }
        } else if (selection === 1) {
            npc.sendOk("Once Pelvis Bebop finalizes the vows, I can send everyone to the afterparty or escort everyone to the exit.")
        }
    }
} else {
    npc.sendOk("I can only assist guests inside the Chapel wedding maps.")
}
