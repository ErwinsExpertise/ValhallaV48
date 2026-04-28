var rooms = [108000100, 108000101, 108000102];

if (!(plr.job() === 300 && plr.getLevel() >= 30)) {
    npc.sendOk("If you're here for the bowman trial, come back when Athena Pierce sends you.");
} else if (!npc.sendYesNo("You've come for the bowman trial. Inside, you'll face monsters that won't reward you with experience or items. Bring back #b30 #t4031013##k to prove your aim is true. Are you ready to enter?")) {
    npc.sendOk("Come back when you're ready to be tested.");
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
