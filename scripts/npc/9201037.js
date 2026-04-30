if (plr.isMarried()) {
    npc.sendOk("You already completed your vows. I hope your marriage stays blessed.")
} else if (plr.partnerID() <= 0) {
    npc.sendOk("Only engaged couples may receive the Parent's Blessing for a Cathedral wedding.")
} else if (plr.haveItem(4031373, 1) || plr.haveItem(4031374, 1)) {
    npc.sendOk("Your couple already received the Parent's Blessing needed for the Cathedral.")
} else if (!plr.canHold(4031373, 1)) {
    npc.sendOk("Please free an ETC slot before I hand over the Parent's Blessing.")
} else if (!npc.sendYesNo("If your couple is truly ready for a Cathedral wedding, I can issue the #b#t4031373##k now. Shall I proceed?")) {
    npc.sendOk("Return when your couple is ready for the blessing.")
} else {
    plr.gainItem(4031373, 1)
    npc.sendOk("Take this #b#t4031373##k to #b#p9201002##k and have it converted into the Officiator's Permission.")
}
