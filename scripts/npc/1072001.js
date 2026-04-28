var rooms = [108000200, 108000201, 108000202];

if (!(plr.job() === 200 && plr.getLevel() >= 30)) {
    npc.sendOk("If you're here for the magician trial, come back when Grendel sends you.");
} else if (!npc.sendYesNo("You've come for the magician trial. Inside, you'll face monsters that won't reward you with experience or items. Bring back #b30 #t4031013##k to prove your skill. Are you ready to enter?")) {
    npc.sendOk("Come back when you're ready to prove your control over magic.");
} else {
    var room = -1;
    for (var i = 0; i < rooms.length; i++) {
        if (map.playerCount(rooms[i], 0) === 0) {
            room = rooms[i];
            break;
        }
    }
    if (room < 0) {
        npc.sendOk("All test chambers are currently occupied. Please wait until one becomes available.");
    } else {
        plr.removeAll(4031013);
        plr.removeAll(4031012);
        plr.warp(room);
    }
}
