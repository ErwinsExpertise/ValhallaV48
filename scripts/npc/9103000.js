if (!plr.isLeader()) {
    npc.sendOk("Please tell the #bParty Leader#k to talk to me.");
} else {
    var coupons = plr.itemCount(4001106);
    if (coupons < 30) {
        npc.sendOk("You must have 30 or more coupons to complete the quest.");
    } else {
        npc.sendNext("Wow, you picked up " + coupons + " coupons! Congratulations! You will be taken to another map to get your #bEXP#k and Rolly will give you your rewards!");
        plr.partyGiveExp(5 * coupons);
        plr.removeAll(4001106);
        plr.warpEventMembers(809050016);
        plr.finishEvent();
    }
}
