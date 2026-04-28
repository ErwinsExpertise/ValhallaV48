if (plr.partnerID() <= 0) {
    npc.sendOk("Welcome to the Cathedral. Couples wanting to marry here should first arrange a reservation with #r#p9201005##k.")
} else if (!plr.hasWeddingReservation(true)) {
    npc.sendOk("There is no Cathedral wedding reservation registered for your couple on this channel.")
} else if (plr.weddingStarted(true)) {
    npc.sendOk("Your Cathedral wedding session has already begun. Please proceed inside and make your way to the altar.")
} else if (npc.sendYesNo("If both you and #b" + plr.partnerName() + "#k are here and ready, I can begin the Cathedral wedding session now. Would you like to start?")) {
    if (plr.startWedding(true)) {
        npc.sendOk("Very well. Let the Cathedral wedding session begin. Please gather your guests in the lounge.")
    } else {
        npc.sendOk("I cannot begin the wedding session right now. Make sure both partners are here, still engaged, and ready.")
    }
} else {
    npc.sendOk("Return when both of you are ready to begin the wedding session.")
}
