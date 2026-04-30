if (plr.mapID() !== 670010800) {
    npc.sendOk("There's nothing for me to do here.");
} else if (!plr.isLeader()) {
    npc.sendOk("Hit the boxes while time remains. If you want to leave early, ask your party leader.");
} else if (npc.sendYesNo("Would you like to leave Amos' Vault now?")) {
    var members = plr.partyMembersOnMap();
    for (var i = 0; i < members.length; i++) {
        if (members[i].mapID() === 670010800) {
            members[i].removeAll(4031592);
            members[i].removeAll(4031593);
            members[i].removeAll(4031594);
            members[i].removeAll(4031595);
            members[i].removeAll(4031596);
            members[i].removeAll(4031597);
        }
    }
    plr.warpEventMembersToPortal(670010000, "st00");
    plr.finishEvent();
} else {
    npc.sendOk("Use the time you have left wisely.");
}
