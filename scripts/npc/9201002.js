if (plr.partnerID() <= 0) {
    npc.sendOk("I oversee Cathedral weddings. Come back with your fiancee when you are ready for the vows.")
} else if (!plr.hasWeddingReservation(true)) {
    npc.sendOk("There is no Cathedral wedding reservation registered for your couple on this channel.")
} else if (!plr.weddingStarted(true)) {
    npc.sendOk("Your Cathedral wedding session has not begun yet. Please speak with #b#p9201005##k first.")
} else if (npc.sendYesNo("If both you and #b" + plr.partnerName() + "#k are ready at the altar, I can complete the Cathedral vows now. Shall I proceed?")) {
    if (plr.completeWedding(true)) {
        npc.sendOk("Very well. By the power vested here, you are now married.")
    } else {
        npc.sendOk("I cannot complete the vows right now. Make sure both partners are here at the altar with their engagement rings.")
    }
} else {
    npc.sendOk("Return when both of you are ready to complete the vows.")
}
