var indoors = plr.mapID() >= 680000100 && plr.mapID() <= 680000500

if (!indoors) {
    npc.sendOk("Please speak with me from inside the wedding venue.")
} else if (!plr.isMarried()) {
    npc.sendOk("Only the newly married couple should be talking to me right now.")
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
        npc.sendOk("The bride and groom must finish their vows first. After that, I can send everyone to the afterparty or escort everyone to the exit.")
    }
}
