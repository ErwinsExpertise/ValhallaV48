if (plr.partnerID() <= 0) {
    npc.sendOk("Come back with your fiancee when you are ready to begin the Chapel wedding session.")
} else if (!plr.hasWeddingReservation(false)) {
    npc.sendOk("There is no Chapel wedding reservation registered for your couple on this channel.")
} else if (plr.weddingStarted(false)) {
    npc.sendOk("Your Chapel wedding session has already begun. Please proceed inside and make your way to the altar.")
} else if (npc.sendYesNo("If both you and #b" + plr.partnerName() + "#k are here and ready, I can begin the Chapel wedding session now. Would you like to start?")) {
    if (plr.startWedding(false)) {
        npc.sendOk("Wonderful. Let the Chapel wedding session begin. Please gather your guests in the lounge.")
    } else {
        npc.sendOk("I cannot begin the ceremony right now. Make sure both partners are here, still engaged, and ready.")
    }
} else {
    npc.sendOk("Return when both of you are ready to begin the wedding session.")
}
